package isql

import (
	"errors"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/common"

	"gorm.io/gorm"
)

type FieldRelationService struct{}

// Exist 判断资源是否存在
func (s FieldRelationService) Exist(filter map[string]any) bool {
	var dataObj model.FieldRelation
	err := common.DB.Debug().Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Count 获取资源总数
func (s FieldRelationService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.FieldRelation{}).Count(&count).Error
	return count, err
}

// Add 创建资源
func (s FieldRelationService) Add(fieldRelation *model.FieldRelation) error {
	return common.DB.Create(fieldRelation).Error
}

// Update 更新资源
func (s FieldRelationService) Update(fieldRelation *model.FieldRelation) error {
	return common.DB.Model(&model.FieldRelation{}).Where("id = ?", fieldRelation.ID).Updates(fieldRelation).Error
}

// Find 获取单个资源
func (s FieldRelationService) Find(filter map[string]any, data *model.FieldRelation) error {
	return common.DB.Where(filter).First(&data).Error
}

// List 获取数据列表
func (s FieldRelationService) List() (fieldRelations []*model.FieldRelation, err error) {
	err = common.DB.Find(&fieldRelations).Error
	return fieldRelations, err
}

// 批量删除资源
func (s FieldRelationService) Delete(fieldRelationIds []uint) error {
	return common.DB.Where("id IN (?)", fieldRelationIds).Unscoped().Delete(&model.FieldRelation{}).Error
}
