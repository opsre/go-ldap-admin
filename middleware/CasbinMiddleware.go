package middleware

import (
	"strings"
	"sync"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
)

var checkLock sync.Mutex

// Casbin middleware based on the RBAC access control model
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := isql.User.GetCurrentLoginUser(c)
		if err != nil {
			tools.Response(c, 401, 401, nil, i18n.TC(c, "auth.not_logged_in", nil))
			c.Abort()
			return
		}
		if user.Status != 1 {
			tools.Response(c, 401, 401, nil, i18n.TC(c, "auth.user_disabled", nil))
			c.Abort()
			return
		}
		// Get all user roles
		roles := user.Roles
		// Get the keywords of all enabled user roles
		var subs []string
		for _, role := range roles {
			if role.Status == 1 {
				subs = append(subs, role.Keyword)
			}
		}
		// Get the request URL path
		obj := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)
		// Get the request method
		act := c.Request.Method
		isPass := check(subs, obj, act)
		if !isPass {
			tools.Response(c, 401, 401, nil, i18n.TC(c, "auth.no_permission", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

func check(subs []string, obj string, act string) bool {
	// Allow only one request to run validation at a time, or validation may fail
	checkLock.Lock()
	defer checkLock.Unlock()
	isPass := false
	for _, sub := range subs {
		pass, _ := common.CasbinEnforcer.Enforce(sub, obj, act)
		if pass {
			isPass = true
			break
		}
	}
	return isPass
}
