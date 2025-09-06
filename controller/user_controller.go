package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

// Add 添加用户记录
// @Summary 添加用户记录
// @Description 添加用户记录
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserAddReq true "添加用户记录的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /user/add [post]
// @Security ApiKeyAuth
func (m *UserController) Add(c *gin.Context) {
	req := new(request.UserAddReq)
	Run(c, req, func() (any, any) {
		return logic.User.Add(c, req)
	})
}

// Update 更新用户记录
// @Summary 更新用户记录
// @Description 添加用户记录
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserUpdateReq true "更改用户记录的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /user/update [post]
// @Security ApiKeyAuth
func (m *UserController) Update(c *gin.Context) {
	req := new(request.UserUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.User.Update(c, req)
	})
}

// List 记录列表
// @Summary 获取所有用户记录列表
// @Description 获取所有用户记录列表
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /user/list [get]
// @Security ApiKeyAuth
func (m *UserController) List(c *gin.Context) {
	req := new(request.UserListReq)
	Run(c, req, func() (any, any) {
		return logic.User.List(c, req)
	})
}

// Delete 删除用户记录
// @Summary 删除用户记录
// @Description 删除用户记录
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserDeleteReq true "删除用户记录的结构体ID"
// @Success 200 {object} response.ResponseBody
// @Router /user/delete [post]
// @Security ApiKeyAuth
func (m UserController) Delete(c *gin.Context) {
	req := new(request.UserDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.User.Delete(c, req)
	})
}

// ChangePwd 更新密码
// @Summary 更新密码
// @Description 更新密码
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserChangePwdReq true "更改用户密码的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /user/changePwd [post]
// @Security ApiKeyAuth
func (m UserController) ChangePwd(c *gin.Context) {
	req := new(request.UserChangePwdReq)
	Run(c, req, func() (any, any) {
		return logic.User.ChangePwd(c, req)
	})
}

// ChangeUserStatus 更改用户状态
// @Summary 更改用户状态
// @Description 更改用户状态
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserChangeUserStatusReq true "更改用户状态的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /user/changeUserStatus [post]
// @Security ApiKeyAuth
func (m UserController) ChangeUserStatus(c *gin.Context) {
	req := new(request.UserChangeUserStatusReq)
	Run(c, req, func() (any, any) {
		return logic.User.ChangeUserStatus(c, req)
	})
}

// GetUserInfo 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Description 获取当前登录用户信息
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /user/info [get]
// @Security ApiKeyAuth
func (uc UserController) GetUserInfo(c *gin.Context) {
	req := new(request.UserGetUserInfoReq)
	Run(c, req, func() (any, any) {
		return logic.User.GetUserInfo(c, req)
	})
}

// SyncDingTalkUsers 同步钉钉用户信息
// @Summary 同步钉钉用户信息
// @Description 同步钉钉用户信息
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.SyncDingUserReq true "同步钉钉用户信息"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncDingTalkUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncDingTalkUsers(c *gin.Context) {
	req := new(request.SyncDingUserReq)
	Run(c, req, func() (any, any) {
		return logic.DingTalk.SyncDingTalkUsers(c, req)
	})
}

// SyncWeComUsers 同步企业微信用户信息
// @Summary 同步企业微信用户信息
// @Description 同步企业微信用户信息
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.SyncWeComUserReq true "同步企业微信用户信息"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncWeComUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncWeComUsers(c *gin.Context) {
	req := new(request.SyncWeComUserReq)
	Run(c, req, func() (any, any) {
		return logic.WeCom.SyncWeComUsers(c, req)
	})
}

// SyncFeiShuUsers 同步飞书用户信息
// @Summary 同步飞书用户信息
// @Description 同步飞书用户信息
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.SyncFeiShuUserReq true "同步飞书用户信息"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncFeiShuUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncFeiShuUsers(c *gin.Context) {
	req := new(request.SyncFeiShuUserReq)
	Run(c, req, func() (any, any) {
		return logic.FeiShu.SyncFeiShuUsers(c, req)
	})
}

// SyncOpenLdapUsers 同步ldap用户信息
// @Summary 同步ldap用户信息
// @Description 同步ldap用户信息
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.SyncOpenLdapUserReq true "同步ldap用户信息"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncOpenLdapUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncOpenLdapUsers(c *gin.Context) {
	req := new(request.SyncOpenLdapUserReq)
	Run(c, req, func() (any, any) {
		return logic.OpenLdap.SyncOpenLdapUsers(c, req)
	})
}

// SyncSqlUsers 同步sql用户信息到ldap
// @Summary 同步sql用户信息到ldap
// @Description 同步sql用户信息到ldap
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.SyncSqlUserReq true "更改用户状态的结构体"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncSqlUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncSqlUsers(c *gin.Context) {
	req := new(request.SyncSqlUserReq)
	Run(c, req, func() (any, any) {
		return logic.Sql.SyncSqlUsers(c, req)
	})
}
