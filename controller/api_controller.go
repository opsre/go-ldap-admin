package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type ApiController struct{}

// List 记录列表
// @Summary 获取API接口列表
// Description: 获取API接口列表
// @Tags 接口管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /api/list [get]
// @Security ApiKeyAuth
func (m *ApiController) List(c *gin.Context) {
	req := new(request.ApiListReq)
	Run(c, req, func() (any, any) {
		return logic.Api.List(c, req)
	})
}

// GetTree 接口树
// @Summary 获取API接口树
// Description: 获取API接口树
// @Tags 接口管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /api/tree [get]
// @Security ApiKeyAuth
func (m *ApiController) GetTree(c *gin.Context) {
	req := new(request.ApiGetTreeReq)
	Run(c, req, func() (any, any) {
		return logic.Api.GetTree(c, req)
	})
}

// Add 新建记录
// @Summary 新建API接口
// @Tags 接口管理
// Description: 新建API接口
// @Accept application/json
// @Produce application/json
// @Param data body request.ApiAddReq true "新建API"
// @Success 200 {object} response.ResponseBody
// @Router /api/add [post]
// @Security ApiKeyAuth
func (m *ApiController) Add(c *gin.Context) {
	req := new(request.ApiAddReq)
	Run(c, req, func() (any, any) {
		return logic.Api.Add(c, req)
	})
}

// Update 更新记录
// @Summary 更新API接口
// @Tags 接口管理
// Description: 更新API接口
// @Accept application/json
// @Produce application/json
// @Param data body request.ApiUpdateReq true "更新API"
// @Success 200 {object} response.ResponseBody
// @Router /api/update [post]
// @Security ApiKeyAuth
func (m *ApiController) Update(c *gin.Context) {
	req := new(request.ApiUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Api.Update(c, req)
	})
}

// Delete 删除记录
// @Summary 删除API接口
// @Tags 接口管理
// Description: 删除API接口
// @Accept application/json
// @Produce application/json
// @Param data body request.ApiDeleteReq true "删除API"
// @Success 200 {object} response.ResponseBody
// @Router /api/delete [post]
// @Security ApiKeyAuth
func (m *ApiController) Delete(c *gin.Context) {
	req := new(request.ApiDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Api.Delete(c, req)
	})
}
