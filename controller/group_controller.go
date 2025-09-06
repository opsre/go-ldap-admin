package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type GroupController struct{}

// List 记录列表
// @Summary 获取分组记录列表
// @Description 获取分组记录列表
// @Tags 分组管理
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

// UserInGroup 在分组内的用户
// @Summary 获取分组内用户
// @Description 获取分组内用户
// @Tags 分组管理
// @Accept application/json
// @Produce application/json
// @Param groupId query int true "分组ID"
// @Param nickname query string false "昵称"
// @Success 200 {object} response.ResponseBody
// @Router /group/useringroup [get]
// @Security ApiKeyAuth
func (m *GroupController) UserInGroup(c *gin.Context) {
	req := new(request.UserInGroupReq)
	Run(c, req, func() (any, any) {
		return logic.Group.UserInGroup(c, req)
	})
}

// UserNoInGroup 不在分组的用户
// @Summary 不在分组的用户
// @Description 不在分组的用户
// @Tags 分组管理
// @Accept application/json
// @Produce application/json
// @Param groupId query int true "分组ID"
// @Param nickname query string false "昵称"
// @Success 200 {object} response.ResponseBody
// @Router /group/usernoingroup [get]
// @Security ApiKeyAuth
func (m *GroupController) UserNoInGroup(c *gin.Context) {
	req := new(request.UserNoInGroupReq)
	Run(c, req, func() (any, any) {
		return logic.Group.UserNoInGroup(c, req)
	})
}

// GetTree 接口树
// @Summary 获取分组接口树
// @Description 获取分组接口树
// @Tags 分组管理
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

// Add 新建分组记录
// @Summary 添加分组记录
// @Description 添加分组记录
// @Tags 分组管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.GroupAddReq true "添加用户记录的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /group/add [post]
// @Security ApiKeyAuth
func (m *GroupController) Add(c *gin.Context) {
	req := new(request.GroupAddReq)
	Run(c, req, func() (any, any) {
		return logic.Group.Add(c, req)
	})
}

// Update 更新记录
// @Summary 更新分组记录
// @Description 更新分组记录
// @Tags 分组管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.GroupUpdateReq true "更新用户记录的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /group/update [post]
// @Security ApiKeyAuth
func (m *GroupController) Update(c *gin.Context) {
	req := new(request.GroupUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Group.Update(c, req)
	})
}

// Delete 删除记录
// @Summary 删除分组记录
// @Description 删除分组记录
// @Tags 分组管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.GroupDeleteReq true "删除用户记录的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /group/delete [post]
// @Security ApiKeyAuth
func (m *GroupController) Delete(c *gin.Context) {
	req := new(request.GroupDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Group.Delete(c, req)
	})
}

// AddUser 添加用户
// @Summary 添加用户
// @Description 添加用户
// @Tags 分组管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.GroupAddUserReq true "添加用户记录的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /group/adduser [post]
// @Security ApiKeyAuth
func (m *GroupController) AddUser(c *gin.Context) {
	req := new(request.GroupAddUserReq)
	Run(c, req, func() (any, any) {
		return logic.Group.AddUser(c, req)
	})
}

// RemoveUser 移除用户
// @Summary 移除用户
// @Description 移除用户
// @Tags 分组管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.GroupRemoveUserReq true "移除用户记录的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /group/removeuser [post]
// @Security ApiKeyAuth
func (m *GroupController) RemoveUser(c *gin.Context) {
	req := new(request.GroupRemoveUserReq)
	Run(c, req, func() (any, any) {
		return logic.Group.RemoveUser(c, req)
	})
}

// SyncDingTalkDepts 同步钉钉部门信息
// @Summary 同步钉钉部门信息
// @Description 同步钉钉部门信息
// @Tags 分组管理
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

// SyncWeComDepts 同步企业微信部门信息
// @Summary 同步企业微信部门信息
// @Description 同步企业微信部门信息
// @Tags 分组管理
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

// SyncFeiShuDepts 同步飞书部门信息
// @Summary 同步飞书部门信息
// @Description 同步飞书部门信息
// @Tags 分组管理
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

// SyncOpenLdapDepts 同步原ldap部门信息
// @Summary 同步原ldap部门信息
// @Description 同步原ldap部门信息
// @Tags 分组管理
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

// SyncSqlGroups 同步Sql中的分组信息到ldap
// @Summary 同步Sql中的分组信息到ldap
// @Description 同步Sql中的分组信息到ldap
// @Tags 分组管理
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
