package controller

import (
	"github.com/eryajf/go-ldap-admin/logic"
	"github.com/eryajf/go-ldap-admin/model/request"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// SendCode 给用户邮箱发送验证码
// @Summary 发送验证码
// @Description 向指定邮箱发送验证码
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseSendCodeReq true "发送验证码请求数据"
// @Success 200 {object} response.ResponseBody
// @Router /base/sendcode [post]
func (m *BaseController) SendCode(c *gin.Context) {
	req := new(request.BaseSendCodeReq)
	Run(c, req, func() (any, any) {
		return logic.Base.SendCode(c, req)
	})
}

// ChangePwd 用户通过邮箱修改密码
// @Summary 用户通过邮箱修改密码
// @Description 使用邮箱验证码修改密码
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.BaseChangePwdReq true "发送验证码请求数据"
// @Success 200 {object} response.ResponseBody
// @Router /base/changePwd [post]
func (m *BaseController) ChangePwd(c *gin.Context) {
	req := new(request.BaseChangePwdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.ChangePwd(c, req)
	})
}

// Dashboard 系统首页展示数据
// @Summary 获取仪表盘数据
// @Description 获取系统仪表盘概览数据
// @Tags 基础管理
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

// EncryptPasswd 密码加密
// @Summary 密码加密
// @Description 将明文密码加密
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param passwd query string true "需要加密的明文密码"
// @Success 200 {object} response.ResponseBody
// @Router /base/encryptpwd [get]
func (m *BaseController) EncryptPasswd(c *gin.Context) {
	req := new(request.EncryptPasswdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.EncryptPasswd(c, req)
	})
}

// DecryptPasswd 密码解密为明文
// @Summary 密码解密
// @Description 将加密后的密码解密为明文
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param passwd query string true "需要解密的加密密码"
// @Success 200 {object} response.ResponseBody
// @Router /base/decryptpwd [get]
func (m *BaseController) DecryptPasswd(c *gin.Context) {
	req := new(request.DecryptPasswdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.DecryptPasswd(c, req)
	})
}
