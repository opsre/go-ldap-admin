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

// var Logs []model.OperationLog // This global variable needs locking across threads, so each thread maintains its own
// Process OperationLogChan and write logs to the database
func (s OperationLogService) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	// Runs only once when the thread starts
	Logs := make([]model.OperationLog, 0)
	// Automatically sync every 5 seconds
	duration := 5 * time.Second
	timer := time.NewTimer(duration)
	defer timer.Stop()
	for {
		select {
		case log := <-olc:
			Logs = append(Logs, *log)
			// Write to the database every 10 records
			if len(Logs) > 5 {
				common.DB.Create(&Logs)
				Logs = make([]model.OperationLog, 0)
				timer.Reset(duration) // Reset the timer after writing to the database
			}
		case <-timer.C: // Automatically sync every 5 seconds
			if len(Logs) > 0 {
				common.DB.Create(&Logs)
				Logs = make([]model.OperationLog, 0)
			}
			timer.Reset(duration) // Reset the timer after writing to the database
		}
	}
}

// List gets the data list
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

// Count gets the total number of records
func (s OperationLogService) Count() (count int64, err error) {
	err = common.DB.Model(&model.OperationLog{}).Count(&count).Error
	return count, err
}

// Get a single user
func (s OperationLogService) Find(filter map[string]any, data *model.OperationLog) error {
	return common.DB.Where(filter).First(&data).Error
}

// Exist checks whether the resource exists
func (s OperationLogService) Exist(filter map[string]any) bool {
	var dataObj model.OperationLog
	err := common.DB.Debug().Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Delete deletes a resource
func (s OperationLogService) Delete(operationLogIds []uint) error {
	return common.DB.Where("id IN (?)", operationLogIds).Unscoped().Delete(&model.OperationLog{}).Error
}

// Clean clears the logs
func (s OperationLogService) Clean() error {
	return common.DB.Exec("DELETE FROM operation_logs").Error
}
