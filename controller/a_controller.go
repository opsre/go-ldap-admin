package controller

import (
	"fmt"
	"net/http"

	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	Api           = &ApiController{}
	Group         = &GroupController{}
	Menu          = &MenuController{}
	Role          = &RoleController{}
	User          = &UserController{}
	OperationLog  = &OperationLogController{}
	Base          = &BaseController{}
	FieldRelation = &FieldRelationController{}
)

func Run(c *gin.Context, req any, fn func() (any, any)) {
	var err error
	// bind struct
	err = c.Bind(req)
	if err != nil {
		tools.Err(c, tools.NewValidatorError(err), nil)
		return
	}
	// 校验
	err = common.Validate.Struct(req)
	if err != nil {
		trans := common.TranslatorForLocale(i18n.LocaleFromContext(c))
		for _, err := range err.(validator.ValidationErrors) {
			tools.Err(c, tools.NewValidatorError(fmt.Errorf("%s", err.Translate(trans))), nil)
			return
		}
	}
	data, err1 := fn()
	if err1 != nil {
		tools.Err(c, tools.ReloadErr(err1), data)
		return
	}
	tools.Success(c, data)
}

// Demo
// @Summary 健康检测
// @Tags 基础管理
// @Produce json
// @Description 健康检测
// @Success 200 {object} response.ResponseBody
// @router /base/ping [get]
func Demo(c *gin.Context) {
	// 健康检测
	CodeDebug()
	c.JSON(http.StatusOK, tools.H{"code": 200, "msg": "ok", "data": "pong"})
}

func CodeDebug() {
}
