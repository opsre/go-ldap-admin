package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type FieldRelationController struct{}

// List 记录列表
// @Summary 获字段关系管理列表
// Description: 获字段关系管理列表
// @Tags 字段关系管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/list [get]
// @Security ApiKeyAuth
func (m *FieldRelationController) List(c *gin.Context) {
	req := new(request.FieldRelationListReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.List(c, req)
	})
}

// Add 新建记录
// @Summary 新建字段关系管理记录
// Description: 新建字段关系管理记录
// @Tags 字段关系管理
// @Accept application/json
// @Produce application/json
// @Param data body request.FieldRelationAddReq true "新建字段关系管理记录"
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/add [post]
// @Security ApiKeyAuth
func (m *FieldRelationController) Add(c *gin.Context) {
	req := new(request.FieldRelationAddReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.Add(c, req)
	})
}

// Update 更新记录
// @Summary 更新字段关系管理记录
// Description: 更新字段关系管理记录
// @Tags 字段关系管理
// @Accept application/json
// @Produce application/json
// @Param data body request.FieldRelationUpdateReq true "更新字段关系管理记录"
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/update [post]
// @Security ApiKeyAuth
func (m *FieldRelationController) Update(c *gin.Context) {
	req := new(request.FieldRelationUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.Update(c, req)
	})
}

// Delete 删除记录
// @Summary 删除字段关系管理记录
// Description: 删除字段关系管理记录
// @Tags 字段关系管理
// @Accept application/json
// @Produce application/json
// @Param data body request.FieldRelationDeleteReq true "删除字段关系管理记录"
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/delete [post]
// @Security ApiKeyAuth
func (m *FieldRelationController) Delete(c *gin.Context) {
	req := new(request.FieldRelationDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.Delete(c, req)
	})
}
