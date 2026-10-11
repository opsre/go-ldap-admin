package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type GroupController struct{}

// List records
// @Summary Get the group list
// @Description Get the group list
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/list [get]
// @Security ApiKeyAuth
func (m *GroupController) List(c *gin.Context) {
	req := new(request.GroupListReq)
	Run(c, req, func() (any, any) {
		return logic.Group.List(c, req)
	})
}

// UserInGroup lists users in the group
// @Summary Get users in the group
// @Description Get users in the group
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Param groupId query int true "Group ID"
// @Param nickname query string false "Nickname"
// @Success 200 {object} response.ResponseBody
// @Router /group/useringroup [get]
// @Security ApiKeyAuth
func (m *GroupController) UserInGroup(c *gin.Context) {
	req := new(request.UserInGroupReq)
	Run(c, req, func() (any, any) {
		return logic.Group.UserInGroup(c, req)
	})
}

// UserNoInGroup lists users outside the group
// @Summary Users outside the group
// @Description Users outside the group
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Param groupId query int true "Group ID"
// @Param nickname query string false "Nickname"
// @Success 200 {object} response.ResponseBody
// @Router /group/usernoingroup [get]
// @Security ApiKeyAuth
func (m *GroupController) UserNoInGroup(c *gin.Context) {
	req := new(request.UserNoInGroupReq)
	Run(c, req, func() (any, any) {
		return logic.Group.UserNoInGroup(c, req)
	})
}

// GetTree gets the API tree
// @Summary Get the group API tree
// @Description Get the group API tree
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/tree [get]
// @Security ApiKeyAuth
func (m *GroupController) GetTree(c *gin.Context) {
	req := new(request.GroupListReq)
	Run(c, req, func() (any, any) {
		return logic.Group.GetTree(c, req)
	})
}

// Add a group record
// @Summary Add a group record
// @Description Add a group record
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Param  data body request.GroupAddReq true "User record request structure"
// @Success 200 {object} response.ResponseBody
// @Router /group/add [post]
// @Security ApiKeyAuth
func (m *GroupController) Add(c *gin.Context) {
	req := new(request.GroupAddReq)
	Run(c, req, func() (any, any) {
		return logic.Group.Add(c, req)
	})
}

// Update a record
// @Summary Update a group record
// @Description Update a group record
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Param  data body request.GroupUpdateReq true "User record update structure"
// @Success 200 {object} response.ResponseBody
// @Router /group/update [post]
// @Security ApiKeyAuth
func (m *GroupController) Update(c *gin.Context) {
	req := new(request.GroupUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Group.Update(c, req)
	})
}

// Delete records
// @Summary Delete a group record
// @Description Delete a group record
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Param  data body request.GroupDeleteReq true "User record deletion structure"
// @Success 200 {object} response.ResponseBody
// @Router /group/delete [post]
// @Security ApiKeyAuth
func (m *GroupController) Delete(c *gin.Context) {
	req := new(request.GroupDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Group.Delete(c, req)
	})
}

// AddUser adds a user
// @Summary Add a user
// @Description Add a user
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Param  data body request.GroupAddUserReq true "User record request structure"
// @Success 200 {object} response.ResponseBody
// @Router /group/adduser [post]
// @Security ApiKeyAuth
func (m *GroupController) AddUser(c *gin.Context) {
	req := new(request.GroupAddUserReq)
	Run(c, req, func() (any, any) {
		return logic.Group.AddUser(c, req)
	})
}

// RemoveUser removes a user
// @Summary Remove a user
// @Description Remove a user
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Param  data body request.GroupRemoveUserReq true "User removal request structure"
// @Success 200 {object} response.ResponseBody
// @Router /group/removeuser [post]
// @Security ApiKeyAuth
func (m *GroupController) RemoveUser(c *gin.Context) {
	req := new(request.GroupRemoveUserReq)
	Run(c, req, func() (any, any) {
		return logic.Group.RemoveUser(c, req)
	})
}

// SyncDingTalkDepts syncs DingTalk department information
// @Summary Sync DingTalk department information
// @Description Sync DingTalk department information
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/syncDingTalkDepts [post]
// @Security ApiKeyAuth
func (m *GroupController) SyncDingTalkDepts(c *gin.Context) {
	req := new(request.SyncDingTalkDeptsReq)
	Run(c, req, func() (any, any) {
		return logic.DingTalk.SyncDingTalkDepts(c, req)
	})
}

// SyncWeComDepts syncs WeCom department information
// @Summary Sync WeCom department information
// @Description Sync WeCom department information
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/syncWeComDepts [post]
// @Security ApiKeyAuth
func (m *GroupController) SyncWeComDepts(c *gin.Context) {
	req := new(request.SyncWeComDeptsReq)
	Run(c, req, func() (any, any) {
		return logic.WeCom.SyncWeComDepts(c, req)
	})
}

// SyncFeiShuDepts syncs Feishu department information
// @Summary Sync Feishu department information
// @Description Sync Feishu department information
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/syncFeiShuDepts [post]
// @Security ApiKeyAuth
func (m *GroupController) SyncFeiShuDepts(c *gin.Context) {
	req := new(request.SyncFeiShuDeptsReq)
	Run(c, req, func() (any, any) {
		return logic.FeiShu.SyncFeiShuDepts(c, req)
	})
}

// SyncOpenLdapDepts syncs existing LDAP department information
// @Summary Sync existing LDAP department information
// @Description Sync existing LDAP department information
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/syncOpenLdapDepts [post]
// @Security ApiKeyAuth
func (m *GroupController) SyncOpenLdapDepts(c *gin.Context) {
	req := new(request.SyncOpenLdapDeptsReq)
	Run(c, req, func() (any, any) {
		return logic.OpenLdap.SyncOpenLdapDepts(c, req)
	})
}

// SyncSqlGroups syncs SQL group information to LDAP
// @Summary Sync SQL group information to LDAP
// @Description Sync SQL group information to LDAP
// @Tags Group management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/syncSqlGroups [post]
// @Security ApiKeyAuth
func (m *GroupController) SyncSqlGroups(c *gin.Context) {
	req := new(request.SyncSqlGrooupsReq)
	Run(c, req, func() (any, any) {
		return logic.Sql.SyncSqlGroups(c, req)
	})
}
