package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/eryajf/go-ldap-admin/controller"
	"github.com/gin-gonic/gin"
)

// LoginHandler
// @Summary Login endpoint (manually add: Bearer + token (password encryption endpoint))
// @Description User login
// @Tags Basic management
// @Accept application/json
// @Produce application/json
// @Param  data body request.RegisterAndLoginReq true "User login information: account and password"
// @Success 200 {object} response.ResponseBody
// @Router /base/login [post]
func LoginHandler() {}

// LogoutHandler
// @Summary Log out
// @Description User logout
// @Tags Basic management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/logout [post]
func LogoutHandler() {
}

// RefreshHandler
// @Summary Refresh token
// @Description Get a new token using the old token
// @Tags Basic management
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer old token"
// @Success 200 {object} response.ResponseBody
// @Router /base/refreshToken [post]
func RefreshHandler() {

}

// Register base routes
func InitBaseRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	base := r.Group("/base")
	{
		base.GET("ping", controller.Demo)
		base.GET("encryptpwd", controller.Base.EncryptPasswd) // Generate an encrypted password
		base.GET("decryptpwd", controller.Base.DecryptPasswd) // Decrypt the password to plain text
		base.GET("config", controller.Base.GetConfig)         // Get the system configuration
		base.GET("version", controller.Base.GetVersion)       // Get version information
		// Login, logout, and token refresh do not require authorization
		base.POST("/login", authMiddleware.LoginHandler)
		base.POST("/logout", authMiddleware.LogoutHandler)
		base.POST("/refreshToken", authMiddleware.RefreshHandler)
		base.POST("/sendcode", controller.Base.SendCode)   // Send a verification code to the user email
		base.POST("/changePwd", controller.Base.ChangePwd) // Change the user password
		base.GET("/dashboard", controller.Base.Dashboard)  // Display system home page data
	}
	return r
}
