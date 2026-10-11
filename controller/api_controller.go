package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type ApiController struct{}

// List records
// @Summary Get the API list
// Description: Get the API list
// @Tags API management
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

// GetTree gets the API tree
// @Summary Get the API tree
// Description: Get the API tree
// @Tags API management
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

// Add a record
// @Summary Create an API
// @Tags API management
// Description: Create an API
// @Accept application/json
// @Produce application/json
// @Param data body request.ApiAddReq true "New API"
// @Success 200 {object} response.ResponseBody
// @Router /api/add [post]
// @Security ApiKeyAuth
func (m *ApiController) Add(c *gin.Context) {
	req := new(request.ApiAddReq)
	Run(c, req, func() (any, any) {
		return logic.Api.Add(c, req)
	})
}

// Update a record
// @Summary Update the API
// @Tags API management
// Description: Update the API
// @Accept application/json
// @Produce application/json
// @Param data body request.ApiUpdateReq true "Update API"
// @Success 200 {object} response.ResponseBody
// @Router /api/update [post]
// @Security ApiKeyAuth
func (m *ApiController) Update(c *gin.Context) {
	req := new(request.ApiUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Api.Update(c, req)
	})
}

// Delete records
// @Summary Delete the API
// @Tags API management
// Description: Delete the API
// @Accept application/json
// @Produce application/json
// @Param data body request.ApiDeleteReq true "Delete API"
// @Success 200 {object} response.ResponseBody
// @Router /api/delete [post]
// @Security ApiKeyAuth
func (m *ApiController) Delete(c *gin.Context) {
	req := new(request.ApiDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Api.Delete(c, req)
	})
}
