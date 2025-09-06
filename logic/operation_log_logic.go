package logic

import (
	"fmt"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/model/response"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
)

type OperationLogLogic struct{}

// List 数据列表
func (l OperationLogLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.OperationLogListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	fmt.Println(r)
	// 获取数据列表
	logs, err := isql.OperationLog.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取接口列表失败: %s", err.Error()))
	}

	rets := make([]model.OperationLog, 0)
	for _, log := range logs {
		rets = append(rets, *log)
	}
	count, err := isql.OperationLog.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取接口总数失败"))
	}

	return response.LogListRsp{
		Total: count,
		Logs:  rets,
	}, nil

	// 获取
	// logs, err := isql.OperationLog.List(&r)
	// if err != nil {
	// 	response.Fail(c, nil, "获取操作日志列表失败: "+err.Error())
	// 	return
	// }
	// return nil, nil
}

// Delete 删除数据
func (l OperationLogLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.OperationLogDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.OperationLogIds {
		filter := tools.H{"id": int(id)}
		if !isql.OperationLog.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("该条记录不存在"))
		}
	}
	// 删除接口
	err := isql.OperationLog.Delete(r.OperationLogIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除该改条记录失败: %s", err.Error()))
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
		return err, nil
	}
	return "操作日志清空完成", nil
}
