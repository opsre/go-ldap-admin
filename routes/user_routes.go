package routes

import (
	"github.com/eryajf/go-ldap-admin/controller"
	"github.com/eryajf/go-ldap-admin/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Register user routes
func InitUserRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	user := r.Group("/user")
	// Enable JWT authentication middleware
	user.Use(authMiddleware.MiddlewareFunc())
	// Enable Casbin authorization middleware
	user.Use(middleware.CasbinMiddleware())
	{
		user.GET("/info", controller.User.GetUserInfo)                   // Not yet complete
		user.GET("/list", controller.User.List)                          // User list
		user.POST("/add", controller.User.Add)                           // Add a user
		user.POST("/update", controller.User.Update)                     // Update the user
		user.POST("/delete", controller.User.Delete)                     // Delete the user
		user.POST("/changePwd", controller.User.ChangePwd)               // Change the user password
		user.POST("/resetPassword", controller.User.ResetPassword)       // Reset the user password
		user.POST("/changeUserStatus", controller.User.ChangeUserStatus) // Change the user status

		user.POST("/syncDingTalkUsers", controller.User.SyncDingTalkUsers) // Sync DingTalk users to the platform
		user.POST("/syncWeComUsers", controller.User.SyncWeComUsers)       // Sync WeCom users to the platform
		user.POST("/syncFeiShuUsers", controller.User.SyncFeiShuUsers)     // Sync Feishu users to the platform
		user.POST("/syncOpenLdapUsers", controller.User.SyncOpenLdapUsers) // Sync LDAP users to the platform
		user.POST("/syncSqlUsers", controller.User.SyncSqlUsers)           // Sync SQL users to LDAP
	}
	return r
}
