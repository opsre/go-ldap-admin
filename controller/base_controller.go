package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// SendCode sends a verification code to the user email
// @Summary Send a verification code
// @Description Send a verification code to the specified email
// @Tags Basic management
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseSendCodeReq true "Verification code request data"
// @Success 200 {object} response.ResponseBody
// @Router /base/sendcode [post]
func (m *BaseController) SendCode(c *gin.Context) {
	req := new(request.BaseSendCodeReq)
	Run(c, req, func() (any, any) {
		return logic.Base.SendCode(c, req)
	})
}

// ChangePwd lets users change their password by email
// @Summary Change password by email
// @Description Change the password using an email verification code
// @Tags Basic management
// @Accept application/json
// @Produce application/json
// @Param  data body request.BaseChangePwdReq true "Verification code request data"
// @Success 200 {object} response.ResponseBody
// @Router /base/changePwd [post]
func (m *BaseController) ChangePwd(c *gin.Context) {
	req := new(request.BaseChangePwdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.ChangePwd(c, req)
	})
}

// Dashboard displays the system home page data
// @Summary Get dashboard data
// @Description Get an overview of the system dashboard
// @Tags Basic management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/dashboard [get]
func (m *BaseController) Dashboard(c *gin.Context) {
	req := new(request.BaseDashboardReq)
	Run(c, req, func() (any, any) {
		return logic.Base.Dashboard(c, req)
	})
}

// EncryptPasswd encrypts a password
// @Summary Encrypt the password
// @Description Encrypt the plain-text password
// @Tags Basic management
// @Accept application/json
// @Produce application/json
// @Param passwd query string true "Plain-text password to encrypt"
// @Success 200 {object} response.ResponseBody
// @Router /base/encryptpwd [get]
func (m *BaseController) EncryptPasswd(c *gin.Context) {
	req := new(request.EncryptPasswdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.EncryptPasswd(c, req)
	})
}

// DecryptPasswd decrypts a password to plain text
// @Summary Decrypt the password
// @Description Decrypt the encrypted password to plain text
// @Tags Basic management
// @Accept application/json
// @Produce application/json
// @Param passwd query string true "Encrypted password to decrypt"
// @Success 200 {object} response.ResponseBody
// @Router /base/decryptpwd [get]
func (m *BaseController) DecryptPasswd(c *gin.Context) {
	req := new(request.DecryptPasswdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.DecryptPasswd(c, req)
	})
}

// GetConfig gets the system configuration
// @Summary Get the system configuration
// @Description Get system configuration so the frontend can determine whether to show the sync button
// @Tags Basic management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/config [get]
func (m *BaseController) GetConfig(c *gin.Context) {
	req := new(request.BaseConfigReq)
	Run(c, req, func() (any, any) {
		return logic.Base.GetConfig(c, req)
	})
}

// GetVersion gets version information
// @Summary Get version information
// @Description Get the system version, Git commit hash, and build time
// @Tags Basic management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/version [get]
func (m *BaseController) GetVersion(c *gin.Context) {
	req := new(request.BaseVersionReq)
	Run(c, req, func() (any, any) {
		return logic.Base.GetVersion(c, req)
	})
}
