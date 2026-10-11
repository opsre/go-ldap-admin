package isql

import (
	"errors"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/common"

	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type MenuService struct{}

// Exist checks whether the resource exists
func (s MenuService) Exist(filter map[string]any) bool {
	var dataObj model.Menu
	err := common.DB.Debug().Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Count gets the total number of resources
func (s MenuService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.Menu{}).Count(&count).Error
	return count, err
}

// Add creates a resource
func (s MenuService) Add(menu *model.Menu) error {
	return common.DB.Create(menu).Error
}

// Update a resource
func (s MenuService) Update(menu *model.Menu) error {
	return common.DB.Model(&model.Menu{}).Where("id = ?", menu.ID).Updates(menu).Error
}

// Find a single resource
func (s MenuService) Find(filter map[string]any, data *model.Menu) error {
	return common.DB.Where(filter).First(&data).Error
}

// List gets the data list
func (s MenuService) List() (menus []*model.Menu, err error) {
	err = common.DB.Order("sort").Find(&menus).Error
	return menus, err
}

// List gets the data list
func (s MenuService) ListUserMenus(roleIds []uint) (menus []*model.Menu, err error) {
	err = common.DB.Where("id IN (select menu_id as id from role_menus where role_id IN (?))", roleIds).Order("sort").Find(&menus).Error
	return menus, err
}

// Delete multiple resources
func (s MenuService) Delete(menuIds []uint) error {
	return common.DB.Where("id IN (?)", menuIds).Select("Roles").Unscoped().Delete(&model.Menu{}).Error
}

// GetUserMenusByUserId gets the permission (accessible) menu list for a user ID
func (s MenuService) GetUserMenusByUserId(userId uint) ([]*model.Menu, error) {
	// Get the user
	var user model.User
	err := common.DB.Where("id = ?", userId).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	// Get the role
	roles := user.Roles
	// The combined menus for all roles
	allRoleMenus := make([]*model.Menu, 0)
	for _, role := range roles {
		var userRole model.Role
		err := common.DB.Where("id = ?", role.ID).Preload("Menus").First(&userRole).Error
		if err != nil {
			return nil, err
		}
		// Get the role menus
		menus := userRole.Menus
		allRoleMenus = append(allRoleMenus, menus...)
	}

	// Deduplicate the combined menus for all roles
	allRoleMenusId := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusId = append(allRoleMenusId, int(menu.ID))
	}
	allRoleMenusIdUniq := funk.UniqInt(allRoleMenusId)
	allRoleMenusUniq := make([]*model.Menu, 0)
	for _, id := range allRoleMenusIdUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.ID) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}

	// Get menus with status 1
	accessMenus := make([]*model.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == 1 {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, err
}

// GenMenuTree generates the menu tree
func GenMenuTree(parentId uint, menus []*model.Menu) []*model.Menu {
	tree := make([]*model.Menu, 0)

	for _, m := range menus {
		if m.ParentId == parentId {
			children := GenMenuTree(m.ID, menus)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// // GetMenuTree gets the menu tree
// func (s MenuService) GetMenuTree() ([]*model.Menu, error) {
// 	var menus []*model.Menu
// 	err := common.DB.Order("sort").Find(&menus).Error
// 	// parentId 0 indicates a root menu
// 	return GenMenuTree(0, menus), err
// }
