package routes

import (
	"github.com/eryajf/go-ldap-admin/controller"
	"github.com/eryajf/go-ldap-admin/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitOperationLogRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	operation_log := r.Group("/log")
	// Enable JWT authentication middleware
	operation_log.Use(authMiddleware.MiddlewareFunc())
	// Enable Casbin authorization middleware
	operation_log.Use(middleware.CasbinMiddleware())
	{
		operation_log.GET("/operation/list", controller.OperationLog.List)
		operation_log.POST("/operation/delete", controller.OperationLog.Delete)
		operation_log.DELETE("/operation/clean", controller.OperationLog.Clean)
	}
	return r
}
