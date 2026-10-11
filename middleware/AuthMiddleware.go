package middleware

import (
	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"time"

	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/model/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Initialize the JWT middleware
func InitAuth() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.Conf.Jwt.Realm,                                 // JWT realm
		Key:             []byte(config.Conf.Jwt.Key),                           // Server key
		Timeout:         time.Hour * time.Duration(config.Conf.Jwt.Timeout),    // Token expiration time
		MaxRefresh:      time.Hour * time.Duration(config.Conf.Jwt.MaxRefresh), // Maximum token refresh time (RefreshToken expiration = Timeout + MaxRefresh)
		PayloadFunc:     payloadFunc,                                           // Handle the payload
		IdentityHandler: identityHandler,                                       // Parse claims
		Authenticator:   login,                                                 // Validate the token and handle login
		Authorizator:    authorizator,                                          // Handle successful user authentication
		Unauthorized:    unauthorized,                                          // Handle failed user authentication
		LoginResponse:   loginResponse,                                         // Response after successful login
		LogoutResponse:  logoutResponse,                                        // Response after logout
		RefreshResponse: refreshResponse,                                       // Response after token refresh
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",    // Automatically look for the request token in these locations
		TokenHeadName:   "Bearer",                                              // Header name
		TimeFunc:        time.Now,
	})
	return authMiddleware, err
}

// Handle the payload
func payloadFunc(data any) jwt.MapClaims {
	if v, ok := data.(tools.H); ok {
		var user model.User
		// Convert user JSON to a struct
		tools.JsonI2Struct(v["user"], &user)
		return jwt.MapClaims{
			jwt.IdentityKey: user.ID,
			"user":          v["user"],
		}
	}
	return jwt.MapClaims{}
}

// Parse claims
func identityHandler(c *gin.Context) any {
	claims := jwt.ExtractClaims(c)
	// The return type map[string]any must match the data type used by payloadFunc and authorizator, or authorization can fail without an obvious cause
	return tools.H{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

// Validate the token and handle login
func login(c *gin.Context) (any, error) {
	var req request.RegisterAndLoginReq
	// Bind the request JSON
	if err := c.ShouldBind(&req); err != nil {
		return "", err
	}

	// Decrypt the password with RSA
	decodeData, err := tools.RSADecrypt([]byte(req.Password), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Username: req.Username,
		Password: string(decodeData),
	}

	// Validate the password
	user, err := isql.User.Login(u)
	if err != nil {
		return nil, err
	}
	// Write the user as JSON for payloadFunc and authorizator
	return tools.H{
		"user": tools.Struct2Json(user),
	}, nil
}

// Handle successful user authentication
func authorizator(data any, c *gin.Context) bool {
	if v, ok := data.(tools.H); ok {
		userStr := v["user"].(string)
		var user model.User
		// Convert user JSON to a struct
		tools.Json2Struct(userStr, &user)
		// Save the user in the context for easy access during API calls
		c.Set("user", user)
		return true
	}
	return false
}

// Handle failed user authentication
func unauthorized(c *gin.Context, code int, message string) {
	common.Log.Debugf("JWT authentication failed, code: %d, message: %s", code, message)
	message = localizeAuthFailure(c, message)
	response.Response(c, code, code, nil, i18n.TC(c, "auth.jwt_failed", i18n.Args{
		"code":    code,
		"message": message,
	}))
}

func localizeAuthFailure(c *gin.Context, message string) string {
	switch message {
	case isql.ErrUserNotLoggedIn.Error():
		return i18n.TC(c, "auth.not_logged_in", nil)
	case isql.ErrUserNotFound.Error():
		return i18n.TC(c, "auth.user_not_found", nil)
	case isql.ErrUserDisabled.Error():
		return i18n.TC(c, "auth.user_disabled", nil)
	case isql.ErrPasswordIncorrect.Error():
		return i18n.TC(c, "auth.password_incorrect", nil)
	default:
		return message
	}
}

// Response after successful login
func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.Response(c, code, code,
		gin.H{
			"token":   token,
			"expires": expires.Format("2006-01-02 15:04:05"),
		},
		i18n.TC(c, "auth.login_success", nil))
}

// Response after logout
func logoutResponse(c *gin.Context, code int) {
	response.Success(c, nil, i18n.TC(c, "auth.logout_success", nil))
}

// Response after token refresh
func refreshResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.Response(c, code, code,
		gin.H{
			"token":   token,
			"expires": expires,
		},
		i18n.TC(c, "auth.refresh_success", nil))
}
