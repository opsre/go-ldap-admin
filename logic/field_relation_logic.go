package logic

import (
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/isql"
	"gorm.io/datatypes"

	"github.com/gin-gonic/gin"
)

type FieldRelationLogic struct{}

// Add 添加数据
func (l FieldRelationLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.FieldRelationAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if isql.FieldRelation.Exist(tools.H{"flag": r.Flag}) {
		return nil, tools.NewValidatorI18nError("field_relation.exists", nil)
	}

	attr, err := tools.MapToJson(r.Attributes)
	if err != nil {
		return nil, tools.NewOperationI18nError("field_relation.map_to_json_failed", i18n.Args{"error": err.Error()})
	}

	frObj := model.FieldRelation{
		Flag:       r.Flag,
		Attributes: datatypes.JSON(attr),
	}

	// 创建接口
	err = isql.FieldRelation.Add(&frObj)
	if err != nil {
		return nil, tools.NewMySqlI18nError("field_relation.create_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// List 数据列表
func (l FieldRelationLogic) List(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.FieldRelationListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 获取数据列表
	frs, err := isql.FieldRelation.List()
	if err != nil {
		return nil, tools.NewMySqlI18nError("field_relation.list_failed", i18n.Args{"error": err.Error()})
	}

	return frs, nil
}

// Update 更新数据
func (l FieldRelationLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.FieldRelationUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"flag": r.Flag}

	if !isql.FieldRelation.Exist(filter) {
		return nil, tools.NewValidatorI18nError("field_relation.not_found", nil)
	}

	oldData := new(model.FieldRelation)
	err := isql.FieldRelation.Find(filter, oldData)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}

	attr, err := tools.MapToJson(r.Attributes)
	if err != nil {
		return nil, tools.NewOperationI18nError("field_relation.map_to_json_failed", i18n.Args{"error": err.Error()})
	}

	frObj := model.FieldRelation{
		Model:      oldData.Model,
		Flag:       r.Flag,
		Attributes: datatypes.JSON(attr),
	}

	err = isql.FieldRelation.Update(&frObj)
	if err != nil {
		return nil, tools.NewMySqlI18nError("field_relation.update_failed", i18n.Args{"error": err.Error()})
	}
	return nil, nil
}

// Delete 删除数据
func (l FieldRelationLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.FieldRelationDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.FieldRelationIds {
		filter := tools.H{"id": int(id)}
		if !isql.FieldRelation.Exist(filter) {
			return nil, tools.NewMySqlI18nError("field_relation.not_found", nil)
		}
	}
	// 删除
	err := isql.FieldRelation.Delete(r.FieldRelationIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("field_relation.delete_failed", i18n.Args{"error": err.Error()})
	}
	return nil, nil
}
