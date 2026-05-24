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

// Casbin中间件, 基于RBAC的权限访问控制模型
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
		// 获得用户的全部角色
		roles := user.Roles
		// 获得用户全部未被禁用的角色的Keyword
		var subs []string
		for _, role := range roles {
			if role.Status == 1 {
				subs = append(subs, role.Keyword)
			}
		}
		// 获得请求路径URL
		obj := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)
		// 获取请求方式
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
	// 同一时间只允许一个请求执行校验, 否则可能会校验失败
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
