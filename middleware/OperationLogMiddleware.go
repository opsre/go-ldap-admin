package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
)

// Operation log channel
var OperationLogChan = make(chan *model.OperationLog, 30)

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		startTime := time.Now()

		// Handle the request
		c.Next()

		// End time
		endTime := time.Now()

		// Execution duration
		timeCost := endTime.Sub(startTime).Milliseconds()

		// Get the current logged-in user
		var username string
		ctxUser, _ := c.Get("user")
		user, ok := ctxUser.(model.User)
		if !ok {
			username = "未登录"
		} else {
			username = user.Username
		}

		// Get the request path
		path := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)
		// Request method
		method := c.Request.Method

		// Get the API description
		api := new(model.Api)
		_ = isql.Api.Find(tools.H{"path": path, "method": method}, api)

		operationLog := model.OperationLog{
			Username:   username,
			Ip:         c.ClientIP(),
			IpLocation: "",
			Method:     method,
			Path:       path,
			Remark:     api.Remark,
			Status:     c.Writer.Status(),
			StartTime:  fmt.Sprintf("%v", startTime),
			TimeCost:   timeCost,
			UserAgent:  c.Request.UserAgent(),
		}

		// Ideally, send logs to RabbitMQ or Kafka
		// For now, send logs to a channel and start three goroutines to process them
		OperationLogChan <- &operationLog
	}
}
