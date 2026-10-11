package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

// Add a user record
// @Summary Add a user record
// @Description Add a user record
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserAddReq true "User record request structure"
// @Success 200 {object} response.ResponseBody
// @Router /user/add [post]
// @Security ApiKeyAuth
func (m *UserController) Add(c *gin.Context) {
	req := new(request.UserAddReq)
	Run(c, req, func() (any, any) {
		return logic.User.Add(c, req)
	})
}

// Update a user record
// @Summary Update a user record
// @Description Add a user record
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserUpdateReq true "User update request structure"
// @Success 200 {object} response.ResponseBody
// @Router /user/update [post]
// @Security ApiKeyAuth
func (m *UserController) Update(c *gin.Context) {
	req := new(request.UserUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.User.Update(c, req)
	})
}

// List records
// @Summary Get the list of all user records
// @Description Get the list of all user records
// @Tags User management
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

// Delete user records
// @Summary Delete user records
// @Description Delete user records
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserDeleteReq true "IDs in the user deletion request structure"
// @Success 200 {object} response.ResponseBody
// @Router /user/delete [post]
// @Security ApiKeyAuth
func (m UserController) Delete(c *gin.Context) {
	req := new(request.UserDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.User.Delete(c, req)
	})
}

// ChangePwd updates the password
// @Summary Update the password
// @Description Update the password
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserChangePwdReq true "User password change request structure"
// @Success 200 {object} response.ResponseBody
// @Router /user/changePwd [post]
// @Security ApiKeyAuth
func (m UserController) ChangePwd(c *gin.Context) {
	req := new(request.UserChangePwdReq)
	Run(c, req, func() (any, any) {
		return logic.User.ChangePwd(c, req)
	})
}

// ResetPassword resets a user password
// @Summary Reset the user password
// @Description Reset the user password to a random value and send an email notification
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserResetPasswordReq true "User password reset request structure"
// @Success 200 {object} response.ResponseBody
// @Router /user/resetPassword [post]
// @Security ApiKeyAuth
func (m UserController) ResetPassword(c *gin.Context) {
	req := new(request.UserResetPasswordReq)
	Run(c, req, func() (any, any) {
		return logic.User.ResetPassword(c, req)
	})
}

// ChangeUserStatus changes the user status
// @Summary Change the user status
// @Description Change the user status
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.UserChangeUserStatusReq true "User status change request structure"
// @Success 200 {object} response.ResponseBody
// @Router /user/changeUserStatus [post]
// @Security ApiKeyAuth
func (m UserController) ChangeUserStatus(c *gin.Context) {
	req := new(request.UserChangeUserStatusReq)
	Run(c, req, func() (any, any) {
		return logic.User.ChangeUserStatus(c, req)
	})
}

// GetUserInfo gets the current logged-in user information
// @Summary Get the current logged-in user information
// @Description Get the current logged-in user information
// @Tags User management
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

// SyncDingTalkUsers syncs DingTalk user information
// @Summary Sync DingTalk user information
// @Description Sync DingTalk user information
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.SyncDingUserReq true "Sync DingTalk user information"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncDingTalkUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncDingTalkUsers(c *gin.Context) {
	req := new(request.SyncDingUserReq)
	Run(c, req, func() (any, any) {
		return logic.DingTalk.SyncDingTalkUsers(c, req)
	})
}

// SyncWeComUsers syncs WeCom user information
// @Summary Sync WeCom user information
// @Description Sync WeCom user information
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.SyncWeComUserReq true "Sync WeCom user information"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncWeComUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncWeComUsers(c *gin.Context) {
	req := new(request.SyncWeComUserReq)
	Run(c, req, func() (any, any) {
		return logic.WeCom.SyncWeComUsers(c, req)
	})
}

// SyncFeiShuUsers syncs Feishu user information
// @Summary Sync Feishu user information
// @Description Sync Feishu user information
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.SyncFeiShuUserReq true "Sync Feishu user information"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncFeiShuUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncFeiShuUsers(c *gin.Context) {
	req := new(request.SyncFeiShuUserReq)
	Run(c, req, func() (any, any) {
		return logic.FeiShu.SyncFeiShuUsers(c, req)
	})
}

// SyncOpenLdapUsers syncs LDAP user information
// @Summary Sync LDAP user information
// @Description Sync LDAP user information
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.SyncOpenLdapUserReq true "Sync LDAP user information"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncOpenLdapUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncOpenLdapUsers(c *gin.Context) {
	req := new(request.SyncOpenLdapUserReq)
	Run(c, req, func() (any, any) {
		return logic.OpenLdap.SyncOpenLdapUsers(c, req)
	})
}

// SyncSqlUsers syncs SQL user information to LDAP
// @Summary Sync SQL user information to LDAP
// @Description Sync SQL user information to LDAP
// @Tags User management
// @Accept application/json
// @Produce application/json
// @Param  data body request.SyncSqlUserReq true "User status change request structure"
// @Success 200 {object} response.ResponseBody
// @Router /user/syncSqlUsers [post]
// @Security ApiKeyAuth
func (uc UserController) SyncSqlUsers(c *gin.Context) {
	req := new(request.SyncSqlUserReq)
	Run(c, req, func() (any, any) {
		return logic.Sql.SyncSqlUsers(c, req)
	})
}
