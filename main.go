package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/eryajf/go-ldap-admin/logic"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/middleware"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/routes"
	"github.com/eryajf/go-ldap-admin/service/isql"
)

// @title Go Ldap Admin
// @version 1.0
// @description OpenLDAP admin system built with Go and Vue
// @termsOfService https://github.com/eryajf/go-ldap-admin

// @contact.name Project author: Er Ya Jiang Fan; Swagger author: Nangong Chengfeng
// @contact.url https://github.com/eryajf/go-ldap-admin
// @contact.email https://github.com/eryajf/go-ldap-admin

// @host 127.0.0.1:8888
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	// Load the configuration file into the global configuration struct
	config.InitConfig()

	if err := i18n.Init(config.Conf.I18n.DefaultLocale, config.Conf.I18n.SupportedLocales); err != nil {
		panic(fmt.Errorf("i18n initialization failed:%s", err))
	}

	// Initialize logging
	common.InitLogger()

	// Initialize the MySQL database
	common.InitDB()

	// Initialize the LDAP connection
	common.InitLDAP()

	// Initialize the Casbin policy manager
	common.InitCasbinEnforcer()

	// Initialize Validator data validation
	common.InitValidate()

	// Initialize MySQL data
	common.InitData()

	// The operation log middleware sends logs to a channel instead of RabbitMQ or Kafka
	// Start three goroutines to process the channel and write logs to the database
	for i := 0; i < 3; i++ {
		go isql.OperationLog.SaveOperationLogChannel(middleware.OperationLogChan)
	}

	// Register all routes
	r := routes.InitRoutes()

	host := "0.0.0.0"
	port := config.Conf.System.Port

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Log.Fatalf("listen: %s\n", err)
		}
	}()

	// Start scheduled tasks
	logic.InitCron()

	common.Log.Info(fmt.Sprintf("Server is running at http://%s:%d", host, port))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(quit, os.Interrupt)
	<-quit
	common.Log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		common.Log.Fatal("Server forced to shutdown:", err)
	}

	common.Log.Info("Server exiting!")

}
