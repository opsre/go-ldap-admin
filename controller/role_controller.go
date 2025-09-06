package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

// List 角色记录列表
// @Summary 获取角色记录列表
// @Description 获取角色记录列表
// @Tags 角色管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /role/list [get]
// @Security ApiKeyAuth
func (m *RoleController) List(c *gin.Context) {
	req := new(request.RoleListReq)
	Run(c, req, func() (any, any) {
		return logic.Role.List(c, req)
	})
}

// Add 新建
// @Summary 新建角色记录
// @Description 新建角色记录
// @Tags 角色管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.RoleAddReq true "添加角色记录的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /role/add [post]
// @Security ApiKeyAuth
func (m *RoleController) Add(c *gin.Context) {
	req := new(request.RoleAddReq)
	Run(c, req, func() (any, any) {
		return logic.Role.Add(c, req)
	})
}

// Update 更新记录
// @Summary 更新角色记录
// @Description 更新角色记录
// @Tags 角色管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.RoleUpdateReq true "更新角色记录的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /role/update [post]
// @Security ApiKeyAuth
func (m *RoleController) Update(c *gin.Context) {
	req := new(request.RoleUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Role.Update(c, req)
	})
}

// Delete 删除记录
// @Summary 删除角色记录
// @Description 删除角色记录
// @Tags 角色管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.RoleDeleteReq true "删除角色记录的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /role/delete [post]
// @Security ApiKeyAuth
func (m *RoleController) Delete(c *gin.Context) {
	req := new(request.RoleDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Role.Delete(c, req)
	})
}

// GetMenuList 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取菜单列表
// @Tags 角色管理
// @Accept application/json
// @Produce application/json
// @Param roleId query int true "角色ID"
// @Success 200 {object} response.ResponseBody
// @Router /role/getmenulist [get]
// @Security ApiKeyAuth
func (m *RoleController) GetMenuList(c *gin.Context) {
	req := new(request.RoleGetMenuListReq)
	Run(c, req, func() (any, any) {
		return logic.Role.GetMenuList(c, req)
	})
}

// GetApiList 获取接口列表
// @Summary 获取接口列表
// @Description 获取接口列表
// @Tags 角色管理
// @Accept application/json
// @Produce application/json
// @Param roleId query int true "角色ID"
// @Success 200 {object} response.ResponseBody
// @Router /role/getapilist [get]
// @Security ApiKeyAuth
func (m *RoleController) GetApiList(c *gin.Context) {
	req := new(request.RoleGetApiListReq)
	Run(c, req, func() (any, any) {
		return logic.Role.GetApiList(c, req)
	})
}

// UpdateMenus 更新菜单
// @Summary 更新菜单
// @Description 更新菜单
// @Tags 角色管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.RoleUpdateMenusReq true "更新菜单的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /role/updatemenus [post]
// @Security ApiKeyAuth
func (m *RoleController) UpdateMenus(c *gin.Context) {
	req := new(request.RoleUpdateMenusReq)
	Run(c, req, func() (any, any) {
		return logic.Role.UpdateMenus(c, req)
	})
}

// UpdateApis 更新接口
// @Summary 更新接口
// @Description 更新接口
// @Tags 角色管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.RoleUpdateApisReq true "更新接口的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /role/updateapis [post]
// @Security ApiKeyAuth
func (m *RoleController) UpdateApis(c *gin.Context) {
	req := new(request.RoleUpdateApisReq)
	Run(c, req, func() (any, any) {
		return logic.Role.UpdateApis(c, req)
	})
}
