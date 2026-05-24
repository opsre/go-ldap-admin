package tools

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/gin-gonic/gin"
)

const (
	SystemErr    = 500
	MySqlErr     = 501
	LdapErr      = 505
	OperationErr = 506
	ValidatorErr = 412
)

type RspError struct {
	code   int
	err    error
	msgKey string
	args   i18n.Args
}

type legacyErrorKey struct {
	raw    string
	key    string
	prefix bool
}

var legacyErrorKeys = []legacyErrorKey{
	{raw: "获取当前登陆用户信息失败", key: "legacy.common.current_user_failed"},
	{raw: "获取当前登陆用户失败", key: "legacy.common.current_user_failed"},
	{raw: "获取资源列表失败: ", key: "legacy.common.resource_list_failed", prefix: true},
	{raw: "该记录不存在", key: "legacy.common.record_not_found"},
	{raw: "有用户不存在", key: "legacy.user.some_not_found"},
	{raw: "该用户不存在", key: "legacy.user.not_found"},
	{raw: "原密码错误", key: "legacy.user.old_password_incorrect"},
	{raw: "原密码解析失败", key: "legacy.user.old_password_parse_failed"},
	{raw: "新密码解析失败", key: "legacy.user.new_password_parse_failed"},
	{raw: "用户已经是在职状态", key: "legacy.user.already_active"},
	{raw: "用户已经是离职状态", key: "legacy.user.already_inactive"},
	{raw: "只有管理员才能更改用户状态", key: "legacy.user.only_admin_change_status"},
	{raw: "获取用户列表失败: ", key: "legacy.user.list_failed", prefix: true},
	{raw: "获取用户总数失败", key: "legacy.user.count_failed"},
	{raw: "添加用户失败", key: "legacy.user.create_failed", prefix: true},
	{raw: "更新用户失败", key: "legacy.user.update_failed", prefix: true},
	{raw: "获取用户信息失败: ", key: "legacy.user.info_failed", prefix: true},
	{raw: "在LDAP删除用户失败", key: "legacy.user.ldap_delete_failed", prefix: true},
	{raw: "在LDAP添加用户失败", key: "legacy.user.ldap_add_failed", prefix: true},
	{raw: "在LDAP更新密码失败", key: "legacy.user.ldap_password_update_failed", prefix: true},
	{raw: "在MySQL删除用户失败: ", key: "legacy.user.mysql_delete_failed", prefix: true},
	{raw: "在MySQL更新密码失败: ", key: "legacy.user.mysql_password_update_failed", prefix: true},
	{raw: "在MySQL更新用户状态失败: ", key: "legacy.user.mysql_status_update_failed", prefix: true},
	{raw: "创建接口失败: ", key: "legacy.api.create_failed", prefix: true},
	{raw: "获取接口列表失败: ", key: "legacy.api.list_failed", prefix: true},
	{raw: "获取接口总数失败", key: "legacy.api.count_failed"},
	{raw: "接口不存在", key: "legacy.api.not_found"},
	{raw: "更新接口失败: ", key: "legacy.api.update_failed", prefix: true},
	{raw: "删除接口失败: ", key: "legacy.api.delete_failed", prefix: true},
	{raw: "该角色名已存在", key: "legacy.role.name_exists"},
	{raw: "该角色名不已存在", key: "legacy.role.not_found"},
	{raw: "获取当前用户最高角色等级失败: ", key: "legacy.role.current_max_sort_failed", prefix: true},
	{raw: "当前用户没有权限更新角色", key: "legacy.role.no_update_permission"},
	{raw: "不能创建比自己等级高或相同等级的角色", key: "legacy.role.cannot_create_higher_or_equal"},
	{raw: "创建角色失败: ", key: "legacy.role.create_failed", prefix: true},
	{raw: "获取菜单列表失败: ", key: "legacy.role.list_failed", prefix: true},
	{raw: "获取角色信息失败: ", key: "legacy.role.info_failed", prefix: true},
	{raw: "不能更新比自己角色等级高或相等的角色", key: "legacy.role.cannot_update_higher_or_equal"},
	{raw: "不能把角色等级更新得比当前用户的等级高或相同", key: "legacy.role.cannot_set_higher_or_equal"},
	{raw: "更新角色失败: ", key: "legacy.role.update_failed", prefix: true},
	{raw: "未能获取到角色信息", key: "legacy.role.no_role_info"},
	{raw: "不能删除比自己角色等级高或相等的角色", key: "legacy.role.cannot_delete_higher_or_equal"},
	{raw: "删除角色失败: ", key: "legacy.role.delete_failed", prefix: true},
	{raw: "获取角色的权限菜单失败: ", key: "legacy.role.menu_list_failed", prefix: true},
	{raw: "获取父级组信息失败", key: "legacy.group.parent_info_failed"},
	{raw: "该分组对应DN已存在", key: "legacy.group.dn_exists"},
	{raw: "向LDAP创建分组失败", key: "legacy.group.ldap_create_failed", prefix: true},
	{raw: "向MySQL创建分组失败", key: "legacy.group.mysql_create_failed"},
	{raw: "添加用户到分组失败: ", key: "legacy.group.add_user_failed", prefix: true},
	{raw: "获取分组列表失败: ", key: "legacy.group.list_failed", prefix: true},
	{raw: "获取分组总数失败", key: "legacy.group.count_failed"},
	{raw: "分组不存在", key: "legacy.group.not_found"},
	{raw: "向LDAP更新分组失败：", key: "legacy.group.ldap_update_failed", prefix: true},
	{raw: "向MySQL更新分组失败", key: "legacy.group.mysql_update_failed"},
	{raw: "有分组不存在", key: "legacy.group.some_not_found"},
	{raw: "存在子分组，请先删除子分组，再执行该分组的删除操作！", key: "legacy.group.has_children"},
	{raw: "向LDAP删除分组失败：", key: "legacy.group.ldap_delete_failed", prefix: true},
	{raw: "获取分组失败: ", key: "legacy.group.get_failed", prefix: true},
	{raw: "ou类型的分组不能添加用户", key: "legacy.group.ou_cannot_add_user"},
	{raw: "ou类型的分组内没有用户", key: "legacy.group.ou_has_no_user"},
	{raw: "向LDAP添加用户到分组失败", key: "legacy.group.ldap_add_user_failed", prefix: true},
	{raw: "处理用户的部门数据失败:", key: "legacy.group.user_department_failed", prefix: true},
	{raw: "在MySQL更新用户失败：", key: "legacy.group.mysql_user_update_failed", prefix: true},
	{raw: "将用户从ldap移除失败", key: "legacy.group.ldap_remove_user_failed", prefix: true},
	{raw: "将用户从MySQL移除失败: ", key: "legacy.group.mysql_remove_user_failed", prefix: true},
	{raw: "通过邮箱查询用户失败", key: "legacy.base.email_query_failed", prefix: true},
	{raw: "LDAP生成新密码失败", key: "legacy.base.ldap_new_password_failed", prefix: true},
	{raw: "获取角色总数失败", key: "legacy.role.count_failed"},
	{raw: "获取菜单总数失败", key: "legacy.menu.count_failed"},
	{raw: "获取日志总数失败", key: "legacy.log.count_failed"},
}

func (re *RspError) Error() string {
	return re.err.Error()
}

func (re *RspError) Code() int {
	return re.code
}

func (re *RspError) MsgKey() string {
	return re.msgKey
}

func (re *RspError) Args() i18n.Args {
	return re.args
}

// NewRspError New
func NewRspError(code int, err error) *RspError {
	return &RspError{
		code: code,
		err:  err,
	}
}

func NewRspI18nError(code int, msgKey string, args i18n.Args) *RspError {
	return &RspError{
		code:   code,
		err:    errors.New(msgKey),
		msgKey: msgKey,
		args:   args,
	}
}

// NewMySqlError mysql错误
func NewMySqlError(err error) *RspError {
	return NewRspError(MySqlErr, err)
}

func NewMySqlI18nError(msgKey string, args i18n.Args) *RspError {
	return NewRspI18nError(MySqlErr, msgKey, args)
}

// NewValidatorError 验证错误
func NewValidatorError(err error) *RspError {
	return NewRspError(ValidatorErr, err)
}

func NewValidatorI18nError(msgKey string, args i18n.Args) *RspError {
	return NewRspI18nError(ValidatorErr, msgKey, args)
}

// NewLdapError ldap错误
func NewLdapError(err error) *RspError {
	return NewRspError(LdapErr, err)
}

func NewLdapI18nError(msgKey string, args i18n.Args) *RspError {
	return NewRspI18nError(LdapErr, msgKey, args)
}

// NewOperationError 操作错误
func NewOperationError(err error) *RspError {
	return NewRspError(OperationErr, err)
}

func NewOperationI18nError(msgKey string, args i18n.Args) *RspError {
	return NewRspI18nError(OperationErr, msgKey, args)
}

// ReloadErr 重新加载错误
func ReloadErr(err any) *RspError {
	rspErr, ok := err.(*RspError)
	if !ok {
		rspError, ok := err.(error)
		if !ok {
			return &RspError{
				code: SystemErr,
				err:  fmt.Errorf("unknow error"),
			}
		}
		return &RspError{
			code: SystemErr,
			err:  rspError,
		}
	}
	return rspErr
}

// Success http 成功
func Success(c *gin.Context, data any) {
	SuccessI18n(c, data, "common.success", nil)
}

func SuccessI18n(c *gin.Context, data any, msgKey string, args i18n.Args) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  i18n.TC(c, msgKey, args),
		"data": data,
	})
}

// Err http 错误
func Err(c *gin.Context, err *RspError, data any) {
	body := gin.H{
		"code": err.Code(),
		"msg":  err.Error(),
		"data": data,
	}
	msgKey, args := err.MsgKey(), err.Args()
	if msgKey == "" {
		msgKey, args = lookupLegacyErrorKey(err.Error())
	}
	if msgKey == "" && i18n.LocaleFromContext(c) != i18n.NormalizeLocale("") && containsCJK(err.Error()) {
		msgKey = fallbackErrorKey(err.Code())
	}
	if msgKey != "" {
		body["msg"] = i18n.TC(c, msgKey, args)
		body["msgKey"] = msgKey
	}
	c.JSON(http.StatusOK, body)
}

func lookupLegacyErrorKey(message string) (string, i18n.Args) {
	for _, item := range legacyErrorKeys {
		if item.prefix {
			if strings.HasPrefix(message, item.raw) {
				return item.key, i18n.Args{"error": strings.TrimPrefix(message, item.raw)}
			}
			continue
		}
		if message == item.raw {
			return item.key, nil
		}
	}
	return "", nil
}

func fallbackErrorKey(code int) string {
	switch code {
	case MySqlErr:
		return "error.database"
	case LdapErr:
		return "error.ldap"
	case OperationErr:
		return "error.operation"
	case ValidatorErr:
		return "error.validator"
	default:
		return "common.unknown_error"
	}
}

func containsCJK(message string) bool {
	for _, r := range message {
		if r >= '\u4e00' && r <= '\u9fff' {
			return true
		}
	}
	return false
}

// 返回前端
func Response(c *gin.Context, httpStatus int, code int, data gin.H, message string) {
	c.JSON(httpStatus, gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	})
}
