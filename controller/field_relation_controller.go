package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type FieldRelationController struct{}

// List records
// @Summary Get the field relation management list
// Description: Get the field relation management list
// @Tags Field relation management
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

// Add a record
// @Summary Create a field relation record
// Description: Create a field relation record
// @Tags Field relation management
// @Accept application/json
// @Produce application/json
// @Param data body request.FieldRelationAddReq true "Create a field relation record"
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/add [post]
// @Security ApiKeyAuth
func (m *FieldRelationController) Add(c *gin.Context) {
	req := new(request.FieldRelationAddReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.Add(c, req)
	})
}

// Update a record
// @Summary Update a field relation record
// Description: Update a field relation record
// @Tags Field relation management
// @Accept application/json
// @Produce application/json
// @Param data body request.FieldRelationUpdateReq true "Update a field relation record"
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/update [post]
// @Security ApiKeyAuth
func (m *FieldRelationController) Update(c *gin.Context) {
	req := new(request.FieldRelationUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.Update(c, req)
	})
}

// Delete records
// @Summary Delete a field relation record
// Description: Delete a field relation record
// @Tags Field relation management
// @Accept application/json
// @Produce application/json
// @Param data body request.FieldRelationDeleteReq true "Delete a field relation record"
// @Success 200 {object} response.ResponseBody
// @Router /fieldrelation/delete [post]
// @Security ApiKeyAuth
func (m *FieldRelationController) Delete(c *gin.Context) {
	req := new(request.FieldRelationDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.FieldRelation.Delete(c, req)
	})
}
