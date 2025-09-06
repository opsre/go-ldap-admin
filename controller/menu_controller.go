package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type MenuController struct{}

// GetTree 菜单树
// @Summary 获取菜单树
// @Tags 菜单管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /menu/tree [get]
// @Security ApiKeyAuth
func (m *MenuController) GetTree(c *gin.Context) {
	req := new(request.MenuGetTreeReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.GetTree(c, req)
	})
}

// GetAccessTree GetUserMenuTreeByUserId 获取用户菜单树
// @Summary 获取用户菜单树
// @Tags 菜单管理
// @Accept application/json
// @Produce application/json
// @Param id query int true "分组ID"
// @Success 200 {object} response.ResponseBody
// @Router /menu/access/tree [get]
// @Security ApiKeyAuth
func (m *MenuController) GetAccessTree(c *gin.Context) {
	req := new(request.MenuGetAccessTreeReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.GetAccessTree(c, req)
	})
}

// Add 新建
// @Summary 新建菜单
// @Tags 菜单管理
// @Accept application/json
// @Produce application/json
// @Param data body request.MenuAddReq true "新建菜单"
// @Success 200 {object} response.ResponseBody
// @Router /menu/add [post]
// @Security ApiKeyAuth
func (m *MenuController) Add(c *gin.Context) {
	req := new(request.MenuAddReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.Add(c, req)
	})
}

// Update 更新记录
// @Summary 更新菜单
// @Tags 菜单管理
// @Accept application/json
// @Produce application/json
// @Param data body request.MenuUpdateReq true "更新菜单"
// @Success 200 {object} response.ResponseBody
// @Router /menu/update [post]
// @Security ApiKeyAuth
func (m *MenuController) Update(c *gin.Context) {
	req := new(request.MenuUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.Update(c, req)
	})
}

// Delete 删除记录
// @Summary 删除菜单
// @Tags 菜单管理
// @Accept application/json
// @Produce application/json
// @Param data body request.MenuDeleteReq true "删除菜单"
// @Success 200 {object} response.ResponseBody
// @Router /menu/delete [post]
// @Security ApiKeyAuth
func (m *MenuController) Delete(c *gin.Context) {
	req := new(request.MenuDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.Delete(c, req)
	})
}
