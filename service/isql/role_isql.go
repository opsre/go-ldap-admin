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

type RoleService struct{}

// Exist checks whether the resource exists
func (s RoleService) Exist(filter map[string]any) bool {
	var dataObj model.Role
	err := common.DB.Debug().Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// List gets the data list
func (s RoleService) List(req *request.RoleListReq) ([]*model.Role, error) {
	var list []*model.Role
	db := common.DB.Model(&model.Role{}).Order("created_at DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error
	return list, err
}

// Count gets the total number of resources
func (s RoleService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.Role{}).Count(&count).Error
	return count, err
}

// Add creates a resource
func (s RoleService) Add(role *model.Role) error {
	return common.DB.Create(role).Error
}

// Update a resource
func (s RoleService) Update(role *model.Role) error {
	return common.DB.Model(&model.Role{}).Where("id = ?", role.ID).Updates(role).Error
}

// Find a single resource
func (s RoleService) Find(filter map[string]any, data *model.Role) error {
	return common.DB.Where(filter).First(&data).Error
}

// Delete deletes a resource
func (s RoleService) Delete(roleIds []uint) error {
	var roles []*model.Role
	err := common.DB.Where("id IN (?)", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	err = common.DB.Select("Users", "Menus").Unscoped().Delete(&roles).Error
	// Delete the Casbin policy after successful deletion
	if err == nil {
		for _, role := range roles {
			roleKeyword := role.Keyword
			rmPolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
			if len(rmPolicies) > 0 {
				isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
				if !isRemoved {
					return errors.New("role deleted, but failed to delete related permission APIs")
				}
			}
		}

	}
	return err
}

// Delete gets a role by role ID
func (s RoleService) GetRolesByIds(roleIds []uint) ([]*model.Role, error) {
	var list []*model.Role
	err := common.DB.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// GetRoleMenusById gets the role permission menus
func (s RoleService) GetRoleMenusById(roleId uint) ([]*model.Menu, error) {
	var role model.Role
	err := common.DB.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// UpdateRoleMenus updates the role permission menus
func (s RoleService) UpdateRoleMenus(role *model.Role) error {
	return common.DB.Model(role).Association("Menus").Replace(role.Menus)
}

// UpdateRoleApis updates role permission APIs (delete all first, then add)
func (s RoleService) UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
	// First get the existing policies for the role ID in the path (to delete first)
	err := common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("failed to load role permission API policy")
	}
	rmPolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	if len(rmPolicies) > 0 {
		isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
		if !isRemoved {
			return errors.New("failed to update role permission APIs")
		}
	}
	isAdded, _ := common.CasbinEnforcer.AddPolicies(reqRolePolicies)
	if !isAdded {
		return errors.New("failed to update role permission APIs")
	}
	err = common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("role permission APIs updated, but failed to load role permission API policy")
	} else {
		return err
	}
}
