package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type OperationLogController struct{}

// List 记录列表
// @Summary 获取操作日志记录列表
// Description: 获取操作日志记录列表
// @Tags 操作日志管理
// @Accept application/json
// @Produce application/json
// @Param username query string false "用户名"
// @Param ip query string false "IP地址"
// @Param path query string false "路径"
// @Param method query string false "方法"
// @Param status query int false "状态码"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.ResponseBody
// @Router /log/operation/list [get]
// @Security ApiKeyAuth
func (m *OperationLogController) List(c *gin.Context) {
	req := new(request.OperationLogListReq)
	Run(c, req, func() (any, any) {
		return logic.OperationLog.List(c, req)
	})
}

// Delete 删除记录
// @Summary 删除操作日志记录
// Description: 删除操作日志记录
// @Tags 操作日志管理
// @Accept application/json
// @Produce application/json
// @Param data body request.OperationLogDeleteReq true "删除日志的ID"
// @Success 200 {object} response.ResponseBody
// @Router /log/operation/delete [post]
// @Security ApiKeyAuth
func (m *OperationLogController) Delete(c *gin.Context) {
	req := new(request.OperationLogDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.OperationLog.Delete(c, req)
	})
}

// Clean 清空记录
// @Summary 清空操作日志记录
// Description: 清空操作日志记录
// @Tags 操作日志管理
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
