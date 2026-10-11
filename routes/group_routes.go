package routes

import (
	"github.com/eryajf/go-ldap-admin/controller"
	"github.com/eryajf/go-ldap-admin/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitGroupRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	group := r.Group("/group")
	// Enable JWT authentication middleware
	group.Use(authMiddleware.MiddlewareFunc())
	// Enable Casbin authorization middleware
	group.Use(middleware.CasbinMiddleware())
	{
		group.GET("/list", controller.Group.List)
		group.GET("/tree", controller.Group.GetTree)
		group.POST("/add", controller.Group.Add)
		group.POST("/update", controller.Group.Update)
		group.POST("/delete", controller.Group.Delete)
		group.POST("/adduser", controller.Group.AddUser)
		group.POST("/removeuser", controller.Group.RemoveUser)

		group.GET("/useringroup", controller.Group.UserInGroup)
		group.GET("/usernoingroup", controller.Group.UserNoInGroup)

		group.POST("/syncDingTalkDepts", controller.Group.SyncDingTalkDepts) // Sync DingTalk departments to the platform
		group.POST("/syncWeComDepts", controller.Group.SyncWeComDepts)       // Sync WeCom departments to the platform
		group.POST("/syncFeiShuDepts", controller.Group.SyncFeiShuDepts)     // Sync Feishu departments to the platform
		group.POST("/syncOpenLdapDepts", controller.Group.SyncOpenLdapDepts) // Sync LDAP groups to the platform InitGroupRoutes
		group.POST("/syncSqlGroups", controller.Group.SyncSqlGroups)         // Sync SQL groups to LDAP
	}

	return r
}
