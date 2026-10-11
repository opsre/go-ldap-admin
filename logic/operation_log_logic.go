package logic

import (
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/model/response"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
)

type OperationLogLogic struct{}

// List data
func (l OperationLogLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.OperationLogListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// Get the data list
	logs, err := isql.OperationLog.List(r)
	if err != nil {
		return nil, tools.NewMySqlI18nError("operation_log.list_failed", i18n.Args{"error": err.Error()})
	}

	rets := make([]model.OperationLog, 0)
	for _, log := range logs {
		rets = append(rets, *log)
	}
	count, err := isql.OperationLog.Count()
	if err != nil {
		return nil, tools.NewMySqlI18nError("operation_log.count_failed", nil)
	}

	return response.LogListRsp{
		Total: count,
		Logs:  rets,
	}, nil

	// Get
	// logs, err := isql.OperationLog.List(&r)
	// if err != nil {
	// 	response.Fail(c, nil, "获取操作日志列表失败: "+err.Error())
	// 	return
	// }
	// return nil, nil
}

// Delete data
func (l OperationLogLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.OperationLogDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.OperationLogIds {
		filter := tools.H{"id": int(id)}
		if !isql.OperationLog.Exist(filter) {
			return nil, tools.NewMySqlI18nError("operation_log.not_found", nil)
		}
	}
	// Delete the API
	err := isql.OperationLog.Delete(r.OperationLogIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("operation_log.delete_failed", i18n.Args{"error": err.Error()})
	}
	return nil, nil
}

func (l OperationLogLogic) Clean(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.OperationLogListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	err := isql.OperationLog.Clean()
	if err != nil {
		return nil, tools.NewMySqlI18nError("operation_log.clean_failed", i18n.Args{"error": err.Error()})
	}
	return nil, nil
}
