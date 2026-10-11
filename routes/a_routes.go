package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eryajf/go-ldap-admin/config"
	_ "github.com/eryajf/go-ldap-admin/docs"
	"github.com/eryajf/go-ldap-admin/middleware"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Initialize
func InitRoutes() *gin.Engine {
	//Set the mode
	gin.SetMode(config.Conf.System.Mode)

	// Create a router with default middleware:
	// Logging and recovery middleware
	r := gin.Default()
	// Create a router without middleware:
	// r := gin.New()
	// r.Use(gin.Recovery())

	r.Use(middleware.Serve("/", middleware.EmbedFolder(static.Static, "dist")))
	r.NoRoute(func(c *gin.Context) {
		data, err := static.Static.ReadFile("dist/index.html")
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	r.Use(middleware.LocaleMiddleware())

	// Enable rate limiting middleware
	// By default, add one token every 50 milliseconds, up to 200 tokens
	fillInterval := time.Duration(config.Conf.RateLimit.FillInterval)
	capacity := config.Conf.RateLimit.Capacity
	r.Use(middleware.RateLimitMiddleware(time.Millisecond*fillInterval, capacity))

	// Enable global CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Enable operation log middleware
	r.Use(middleware.OperationLogMiddleware())

	// Initialize JWT authentication middleware
	authMiddleware, err := middleware.InitAuth()
	if err != nil {
		common.Log.Panicf("failed to initialize JWT middleware：%v", err)
		panic(fmt.Sprintf("failed to initialize JWT middleware：%v", err))
	}
	// swag
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiGroup := r.Group("/" + config.Conf.System.UrlPathPrefix)
	// swag
	// Register routes
	InitBaseRoutes(apiGroup, authMiddleware)          // Register base routes without JWT authentication or Casbin middleware
	InitUserRoutes(apiGroup, authMiddleware)          // Register user routes with JWT authentication and Casbin authorization middleware
	InitGroupRoutes(apiGroup, authMiddleware)         // Register group routes with JWT authentication and Casbin authorization middleware
	InitRoleRoutes(apiGroup, authMiddleware)          // Register role routes with JWT authentication and Casbin authorization middleware
	InitMenuRoutes(apiGroup, authMiddleware)          // Register menu routes with JWT authentication and Casbin authorization middleware
	InitApiRoutes(apiGroup, authMiddleware)           // Register API routes with JWT authentication and Casbin authorization middleware
	InitOperationLogRoutes(apiGroup, authMiddleware)  // Register operation log routes with JWT authentication and Casbin authorization middleware
	InitFieldRelationRoutes(apiGroup, authMiddleware) // Register operation log routes with JWT authentication and Casbin authorization middleware

	common.Log.Info("routes initialized")
	return r
}
