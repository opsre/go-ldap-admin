package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type OperationLogController struct{}

// List records
// @Summary Get the operation log list
// Description: Get the operation log list
// @Tags Operation log management
// @Accept application/json
// @Produce application/json
// @Param username query string false "Username"
// @Param ip query string false "IP address"
// @Param path query string false "Path"
// @Param method query string false "Method"
// @Param status query int false "Status code"
// @Param pageNum query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} response.ResponseBody
// @Router /log/operation/list [get]
// @Security ApiKeyAuth
func (m *OperationLogController) List(c *gin.Context) {
	req := new(request.OperationLogListReq)
	Run(c, req, func() (any, any) {
		return logic.OperationLog.List(c, req)
	})
}

// Delete records
// @Summary Delete operation log records
// Description: Delete operation log records
// @Tags Operation log management
// @Accept application/json
// @Produce application/json
// @Param data body request.OperationLogDeleteReq true "IDs of the logs to delete"
// @Success 200 {object} response.ResponseBody
// @Router /log/operation/delete [post]
// @Security ApiKeyAuth
func (m *OperationLogController) Delete(c *gin.Context) {
	req := new(request.OperationLogDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.OperationLog.Delete(c, req)
	})
}

// Clean clears the records
// @Summary Clear operation log records
// Description: Clear operation log records
// @Tags Operation log management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /log/operation/clean [delete]
// @Security ApiKeyAuth
func (m *OperationLogController) Clean(c *gin.Context) {
	req := new(request.OperationLogListReq)
	Run(c, req, func() (any, any) {
		return logic.OperationLog.Clean(c, req)
	})
}
