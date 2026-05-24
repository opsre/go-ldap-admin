package logic

import (
	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/model/response"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/public/version"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
)

type BaseLogic struct{}

// SendCode 发送验证码
func (l BaseLogic) SendCode(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseSendCodeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// 判断邮箱是否正确
	user := new(model.User)
	err := isql.User.Find(tools.H{"mail": r.Mail}, user)
	if err != nil {
		return nil, tools.NewMySqlI18nError("base.email_query_failed", i18n.Args{"error": err.Error()})
	}
	if user.Status != 1 || user.SyncState != 1 {
		return nil, tools.NewMySqlI18nError("base.reset_password_user_unavailable", nil)
	}
	err = tools.SendCodeI18n([]string{r.Mail}, i18n.LocaleFromContext(c))
	if err != nil {
		return nil, tools.NewLdapI18nError("base.send_email_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// ChangePwd 重置密码
func (l BaseLogic) ChangePwd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseChangePwdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// 判断邮箱是否正确
	if !isql.User.Exist(tools.H{"mail": r.Mail}) {
		return nil, tools.NewValidatorI18nError("base.email_not_found", nil)
	}
	// 判断验证码是否过期
	cacheCode, ok := tools.VerificationCodeCache.Get(r.Mail)
	if !ok {
		return nil, tools.NewValidatorI18nError("base.code_expired", nil)
	}
	// 判断验证码是否正确
	if cacheCode != r.Code {
		return nil, tools.NewValidatorI18nError("base.code_invalid", nil)
	}

	user := new(model.User)
	err := isql.User.Find(tools.H{"mail": r.Mail}, user)
	if err != nil {
		return nil, tools.NewMySqlI18nError("base.email_query_failed", i18n.Args{"error": err.Error()})
	}

	newpass, err := ildap.User.NewPwd(user.Username)
	if err != nil {
		return nil, tools.NewLdapI18nError("base.ldap_new_password_failed", i18n.Args{"error": err.Error()})
	}

	err = tools.SendMailI18n([]string{user.Mail}, newpass, i18n.LocaleFromContext(c))
	if err != nil {
		return nil, tools.NewLdapI18nError("base.send_email_failed", i18n.Args{"error": err.Error()})
	}

	// 更新数据库密码
	err = isql.User.ChangePwd(user.Username, tools.NewGenPasswd(newpass))
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.mysql_password_update_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// Dashboard 仪表盘
func (l BaseLogic) Dashboard(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.BaseDashboardReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	userCount, err := isql.User.Count()
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.count_failed", i18n.Args{"error": err.Error()})
	}
	groupCount, err := isql.Group.Count()
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.count_failed", nil)
	}
	roleCount, err := isql.Role.Count()
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.count_failed", nil)
	}
	menuCount, err := isql.Menu.Count()
	if err != nil {
		return nil, tools.NewMySqlI18nError("menu.count_failed", nil)
	}
	apiCount, err := isql.Api.Count()
	if err != nil {
		return nil, tools.NewMySqlI18nError("api.count_failed", nil)
	}
	logCount, err := isql.OperationLog.Count()
	if err != nil {
		return nil, tools.NewMySqlI18nError("operation_log.count_failed", nil)
	}

	rst := make([]*response.DashboardList, 0)

	rst = append(rst,
		&response.DashboardList{
			DataType:  "user",
			DataName:  i18n.TC(c, "dashboard.user", nil),
			DataCount: userCount,
			Icon:      "people",
			Path:      "#/personnel/user",
		},
		&response.DashboardList{
			DataType:  "group",
			DataName:  i18n.TC(c, "dashboard.group", nil),
			DataCount: groupCount,
			Icon:      "peoples",
			Path:      "#/personnel/group",
		},
		&response.DashboardList{
			DataType:  "role",
			DataName:  i18n.TC(c, "dashboard.role", nil),
			DataCount: roleCount,
			Icon:      "eye-open",
			Path:      "#/system/role",
		},
		&response.DashboardList{
			DataType:  "menu",
			DataName:  i18n.TC(c, "dashboard.menu", nil),
			DataCount: menuCount,
			Icon:      "tree-table",
			Path:      "#/system/menu",
		},
		&response.DashboardList{
			DataType:  "api",
			DataName:  i18n.TC(c, "dashboard.api", nil),
			DataCount: apiCount,
			Icon:      "tree",
			Path:      "#/system/api",
		},
		&response.DashboardList{
			DataType:  "log",
			DataName:  i18n.TC(c, "dashboard.log", nil),
			DataCount: logCount,
			Icon:      "documentation",
			Path:      "#/log/operation-log",
		},
	)

	return rst, nil
}

// EncryptPasswd
func (l BaseLogic) EncryptPasswd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.EncryptPasswdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	passwd, err := tools.RSAEncrypt([]byte(r.Passwd), config.Conf.System.RSAPublicBytes)
	if err != nil {
		return nil, tools.NewValidatorI18nError("user.password_encrypt_failed", nil)
	}
	return string(passwd), nil
}

// DecryptPasswd
func (l BaseLogic) DecryptPasswd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.DecryptPasswdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	passwd, err := tools.RSADecrypt([]byte(r.Passwd), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, tools.NewValidatorI18nError("user.password_decrypt_failed", nil)
	}
	return string(passwd), nil
}

// GetConfig 获取系统配置
func (l BaseLogic) GetConfig(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.BaseConfigReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 安全获取配置值，防止配置段缺失导致空指针
	rsp := &response.BaseConfigRsp{}
	if config.Conf.Ldap != nil {
		rsp.LdapEnableSync = config.Conf.Ldap.EnableSync
	}
	if config.Conf.DingTalk != nil {
		rsp.DingTalkEnableSync = config.Conf.DingTalk.EnableSync
	}
	if config.Conf.FeiShu != nil {
		rsp.FeiShuEnableSync = config.Conf.FeiShu.EnableSync
	}
	if config.Conf.WeCom != nil {
		rsp.WeComEnableSync = config.Conf.WeCom.EnableSync
	}

	return rsp, nil
}

// GetVersion 获取版本信息
func (l BaseLogic) GetVersion(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.BaseVersionReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	return version.GetVersion(), nil
}
