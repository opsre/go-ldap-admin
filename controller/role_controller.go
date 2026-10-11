package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

// List role records
// @Summary Get the role list
// @Description Get the role list
// @Tags Role management
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

// Add a new item
// @Summary Create a role record
// @Description Create a role record
// @Tags Role management
// @Accept application/json
// @Produce application/json
// @Param  data body request.RoleAddReq true "Role creation request structure"
// @Success 200 {object} response.ResponseBody
// @Router /role/add [post]
// @Security ApiKeyAuth
func (m *RoleController) Add(c *gin.Context) {
	req := new(request.RoleAddReq)
	Run(c, req, func() (any, any) {
		return logic.Role.Add(c, req)
	})
}

// Update a record
// @Summary Update a role record
// @Description Update a role record
// @Tags Role management
// @Accept application/json
// @Produce application/json
// @Param  data body request.RoleUpdateReq true "Role update request structure"
// @Success 200 {object} response.ResponseBody
// @Router /role/update [post]
// @Security ApiKeyAuth
func (m *RoleController) Update(c *gin.Context) {
	req := new(request.RoleUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Role.Update(c, req)
	})
}

// Delete records
// @Summary Delete role records
// @Description Delete role records
// @Tags Role management
// @Accept application/json
// @Produce application/json
// @Param  data body request.RoleDeleteReq true "Role deletion request structure"
// @Success 200 {object} response.ResponseBody
// @Router /role/delete [post]
// @Security ApiKeyAuth
func (m *RoleController) Delete(c *gin.Context) {
	req := new(request.RoleDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Role.Delete(c, req)
	})
}

// GetMenuList gets the menu list
// @Summary Get the menu list
// @Description Get the menu list
// @Tags Role management
// @Accept application/json
// @Produce application/json
// @Param roleId query int true "Role ID"
// @Success 200 {object} response.ResponseBody
// @Router /role/getmenulist [get]
// @Security ApiKeyAuth
func (m *RoleController) GetMenuList(c *gin.Context) {
	req := new(request.RoleGetMenuListReq)
	Run(c, req, func() (any, any) {
		return logic.Role.GetMenuList(c, req)
	})
}

// GetApiList gets the API list
// @Summary Get the API list
// @Description Get the API list
// @Tags Role management
// @Accept application/json
// @Produce application/json
// @Param roleId query int true "Role ID"
// @Success 200 {object} response.ResponseBody
// @Router /role/getapilist [get]
// @Security ApiKeyAuth
func (m *RoleController) GetApiList(c *gin.Context) {
	req := new(request.RoleGetApiListReq)
	Run(c, req, func() (any, any) {
		return logic.Role.GetApiList(c, req)
	})
}

// UpdateMenus updates menus
// @Summary Update the menu
// @Description Update the menu
// @Tags Role management
// @Accept application/json
// @Produce application/json
// @Param  data body request.RoleUpdateMenusReq true "Menu update request structure"
// @Success 200 {object} response.ResponseBody
// @Router /role/updatemenus [post]
// @Security ApiKeyAuth
func (m *RoleController) UpdateMenus(c *gin.Context) {
	req := new(request.RoleUpdateMenusReq)
	Run(c, req, func() (any, any) {
		return logic.Role.UpdateMenus(c, req)
	})
}

// UpdateApis updates APIs
// @Summary Update APIs
// @Description Update APIs
// @Tags Role management
// @Accept application/json
// @Produce application/json
// @Param  data body request.RoleUpdateApisReq true "API update request structure"
// @Success 200 {object} response.ResponseBody
// @Router /role/updateapis [post]
// @Security ApiKeyAuth
func (m *RoleController) UpdateApis(c *gin.Context) {
	req := new(request.RoleUpdateApisReq)
	Run(c, req, func() (any, any) {
		return logic.Role.UpdateApis(c, req)
	})
}
