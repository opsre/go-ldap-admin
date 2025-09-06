package isql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"

	"gorm.io/gorm"
)

type GroupService struct{}

// List 获取数据列表
func (s GroupService) List(req *request.GroupListReq) ([]*model.Group, error) {
	var list []*model.Group
	db := common.DB.Model(&model.Group{}).Order("created_at DESC")

	groupName := strings.TrimSpace(req.GroupName)
	if groupName != "" {
		db = db.Where("group_name LIKE ?", fmt.Sprintf("%%%s%%", groupName))
	}
	groupRemark := strings.TrimSpace(req.Remark)
	if groupRemark != "" {
		db = db.Where("remark LIKE ?", fmt.Sprintf("%%%s%%", groupRemark))
	}
	syncState := req.SyncState
	if syncState != 0 {
		db = db.Where("sync_state = ?", syncState)
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Preload("Users").Find(&list).Error
	return list, err
}

// List 获取数据列表
func (s GroupService) ListTree(req *request.GroupListReq) ([]*model.Group, error) {
	var list []*model.Group
	db := common.DB.Model(&model.Group{}).Order("created_at DESC")

	groupName := strings.TrimSpace(req.GroupName)
	if groupName != "" {
		db = db.Where("group_name LIKE ?", fmt.Sprintf("%%%s%%", groupName))
	}
	groupRemark := strings.TrimSpace(req.Remark)
	if groupRemark != "" {
		db = db.Where("remark LIKE ?", fmt.Sprintf("%%%s%%", groupRemark))
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error
	return list, err
}

// List 获取数据列表
func (s GroupService) ListAll() (list []*model.Group, err error) {
	err = common.DB.Model(&model.Group{}).Order("created_at DESC").Find(&list).Error

	return list, err
}

// GenGroupTree 生成分组树
func GenGroupTree(parentId uint, groups []*model.Group) []*model.Group {
	tree := make([]*model.Group, 0)

	for _, g := range groups {
		if g.ParentId == parentId {
			children := GenGroupTree(g.ID, groups)
			g.Children = children
			tree = append(tree, g)
		}
	}
	return tree
}

// Count 获取数据总数
func (s GroupService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.Group{}).Count(&count).Error
	return count, err
}

// Add 添加资源
func (s GroupService) Add(data *model.Group) error {
	return common.DB.Create(data).Error
}

// Update 更新资源
func (s GroupService) Update(dataObj *model.Group) error {
	return common.DB.Model(dataObj).Where("id = ?", dataObj.ID).Updates(dataObj).Error
}

// ChangeSyncState 更新分组的同步状态
func (s GroupService) ChangeSyncState(id, status int) error {
	return common.DB.Model(&model.Group{}).Where("id = ?", id).Update("sync_state", status).Error
}

// Find 获取单个资源
func (s GroupService) Find(filter map[string]any, data *model.Group, args ...any) error {
	return common.DB.Where(filter, args).Preload("Users").First(&data).Error
}

// Exist 判断资源是否存在
func (s GroupService) Exist(filter map[string]any) bool {
	var dataObj model.Group
	err := common.DB.Debug().Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Delete 批量删除
func (s GroupService) Delete(groups []*model.Group) error {
	return common.DB.Debug().Select("Users").Unscoped().Delete(&groups).Error
}

// GetApisById 根据接口ID获取接口列表
func (s GroupService) GetGroupByIds(ids []uint) (datas []*model.Group, err error) {
	err = common.DB.Where("id IN (?)", ids).Preload("Users").Find(&datas).Error
	return datas, err
}

// AddUserToGroup 添加用户到分组
func (s GroupService) AddUserToGroup(group *model.Group, users []model.User) error {
	return common.DB.Model(&group).Association("Users").Append(users)
}

// RemoveUserFromGroup 将用户从分组移除
func (s GroupService) RemoveUserFromGroup(group *model.Group, users []model.User) error {
	return common.DB.Model(&group).Association("Users").Delete(users)
}

// DeptIdsToGroupIds 将企业IM部门id转换为MySQL分组id
func (s GroupService) DeptIdsToGroupIds(ids []string) (groupIds []uint, err error) {
	var tempGroups []model.Group
	err = common.DB.Model(&model.Group{}).Where("source_dept_id IN (?)", ids).Find(&tempGroups).Error
	if err != nil {
		return nil, err
	}
	var tempGroupIds []uint
	for _, g := range tempGroups {
		tempGroupIds = append(tempGroupIds, g.ID)
	}
	return tempGroupIds, nil
}
