package routes

import (
	"github.com/eryajf/go-ldap-admin/controller"
	"github.com/eryajf/go-ldap-admin/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 注册用户路由
func InitUserRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	user := r.Group("/user")
	// 开启jwt认证中间件
	user.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	user.Use(middleware.CasbinMiddleware())
	{
		user.GET("/info", controller.User.GetUserInfo)                   // 暂时未完成
		user.GET("/list", controller.User.List)                          // 用户列表
		user.POST("/add", controller.User.Add)                           // 添加用户
		user.POST("/update", controller.User.Update)                     // 更新用户
		user.POST("/delete", controller.User.Delete)                     // 删除用户
		user.POST("/changePwd", controller.User.ChangePwd)               // 修改用户密码
		user.POST("/resetPassword", controller.User.ResetPassword)       // 重置用户密码
		user.POST("/changeUserStatus", controller.User.ChangeUserStatus) // 修改用户状态

		user.POST("/syncDingTalkUsers", controller.User.SyncDingTalkUsers) // 同步钉钉用户到平台
		user.POST("/syncWeComUsers", controller.User.SyncWeComUsers)       // 同步企业微信用户到平台
		user.POST("/syncFeiShuUsers", controller.User.SyncFeiShuUsers)     // 同步飞书用户到平台
		user.POST("/syncOpenLdapUsers", controller.User.SyncOpenLdapUsers) // 同步Ldap用户到平台
		user.POST("/syncSqlUsers", controller.User.SyncSqlUsers)           // 同步Sql用户到Ldap
	}
	return r
}
