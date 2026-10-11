package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type MenuController struct{}

// GetTree gets the menu tree
// @Summary Get the menu tree
// @Tags Menu management
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

// GetAccessTree / GetUserMenuTreeByUserId gets the user menu tree
// @Summary Get the user menu tree
// @Tags Menu management
// @Accept application/json
// @Produce application/json
// @Param id query int true "Group ID"
// @Success 200 {object} response.ResponseBody
// @Router /menu/access/tree [get]
// @Security ApiKeyAuth
func (m *MenuController) GetAccessTree(c *gin.Context) {
	req := new(request.MenuGetAccessTreeReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.GetAccessTree(c, req)
	})
}

// Add a new item
// @Summary Create a menu
// @Tags Menu management
// @Accept application/json
// @Produce application/json
// @Param data body request.MenuAddReq true "Create a menu"
// @Success 200 {object} response.ResponseBody
// @Router /menu/add [post]
// @Security ApiKeyAuth
func (m *MenuController) Add(c *gin.Context) {
	req := new(request.MenuAddReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.Add(c, req)
	})
}

// Update a record
// @Summary Update the menu
// @Tags Menu management
// @Accept application/json
// @Produce application/json
// @Param data body request.MenuUpdateReq true "Update the menu"
// @Success 200 {object} response.ResponseBody
// @Router /menu/update [post]
// @Security ApiKeyAuth
func (m *MenuController) Update(c *gin.Context) {
	req := new(request.MenuUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.Update(c, req)
	})
}

// Delete records
// @Summary Delete the menu
// @Tags Menu management
// @Accept application/json
// @Produce application/json
// @Param data body request.MenuDeleteReq true "Delete the menu"
// @Success 200 {object} response.ResponseBody
// @Router /menu/delete [post]
// @Security ApiKeyAuth
func (m *MenuController) Delete(c *gin.Context) {
	req := new(request.MenuDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Menu.Delete(c, req)
	})
}
