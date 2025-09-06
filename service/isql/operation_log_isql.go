package isql

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"

	"gorm.io/gorm"
)

type OperationLogService struct{}

// var Logs []model.OperationLog //全局变量多个线程需要加锁，所以每个线程自己维护一个
// 处理OperationLogChan将日志记录到数据库
func (s OperationLogService) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	// 只会在线程开启的时候执行一次
	Logs := make([]model.OperationLog, 0)
	// 5s 自动同步一次
	duration := 5 * time.Second
	timer := time.NewTimer(duration)
	defer timer.Stop()
	for {
		select {
		case log := <-olc:
			Logs = append(Logs, *log)
			// 每10条记录到数据库
			if len(Logs) > 5 {
				common.DB.Create(&Logs)
				Logs = make([]model.OperationLog, 0)
				timer.Reset(duration) // 入库重置定时器
			}
		case <-timer.C: // 5s 自动同步一次
			if len(Logs) > 0 {
				common.DB.Create(&Logs)
				Logs = make([]model.OperationLog, 0)
			}
			timer.Reset(duration) // 入库重置定时器
		}
	}
}

// List 获取数据列表
func (s OperationLogService) List(req *request.OperationLogListReq) ([]*model.OperationLog, error) {
	var list []*model.OperationLog
	db := common.DB.Model(&model.OperationLog{}).Order("id DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	ip := strings.TrimSpace(req.Ip)
	if ip != "" {
		db = db.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	method := strings.TrimSpace(req.Method)
	if method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error

	return list, err
}

// Count 获取数据总数
func (s OperationLogService) Count() (count int64, err error) {
	err = common.DB.Model(&model.OperationLog{}).Count(&count).Error
	return count, err
}

// 获取单个用户
func (s OperationLogService) Find(filter map[string]any, data *model.OperationLog) error {
	return common.DB.Where(filter).First(&data).Error
}

// Exist 判断资源是否存在
func (s OperationLogService) Exist(filter map[string]any) bool {
	var dataObj model.OperationLog
	err := common.DB.Debug().Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Delete 删除资源
func (s OperationLogService) Delete(operationLogIds []uint) error {
	return common.DB.Where("id IN (?)", operationLogIds).Unscoped().Delete(&model.OperationLog{}).Error
}

// Clean 清空日志
func (s OperationLogService) Clean() error {
	return common.DB.Exec("DELETE FROM operation_logs").Error
}
