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

type ApiService struct{}

// List gets the data list
func (s ApiService) List(req *request.ApiListReq) ([]*model.Api, error) {
	var list []*model.Api
	db := common.DB.Model(&model.Api{}).Order("created_at DESC")

	method := strings.TrimSpace(req.Method)
	if method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		db = db.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error
	return list, err
}

// List gets the data list
func (s ApiService) ListAll() (list []*model.Api, err error) {
	err = common.DB.Model(&model.Api{}).Order("created_at DESC").Find(&list).Error

	return list, err
}

// Count gets the total number of records
func (s ApiService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.Api{}).Count(&count).Error
	return count, err
}

// Add a resource
func (s ApiService) Add(api *model.Api) error {
	return common.DB.Create(api).Error
}

// Update a resource
func (s ApiService) Update(api *model.Api) error {
	// Get API information by ID
	var oldApi model.Api
	err := common.DB.First(&oldApi, api.ID).Error
	if err != nil {
		return errors.New("failed to get API information by API ID")
	}
	err = common.DB.Model(api).Where("id = ?", api.ID).Updates(api).Error
	if err != nil {
		return err
	}
	// Update the Casbin policy when method or path changes
	if oldApi.Path != api.Path || oldApi.Method != api.Method {
		policies := common.CasbinEnforcer.GetFilteredPolicy(1, oldApi.Path, oldApi.Method)
		// Proceed only if the API exists in the Casbin policy
		if len(policies) > 0 {
			// Delete first
			isRemoved, _ := common.CasbinEnforcer.RemovePolicies(policies)
			if !isRemoved {
				return errors.New("failed to update permission API")
			}
			for _, policy := range policies {
				policy[1] = api.Path
				policy[2] = api.Method
			}
			// Add
			isAdded, _ := common.CasbinEnforcer.AddPolicies(policies)
			if !isAdded {
				return errors.New("failed to update permission API")
			}
			// Load the policy
			err := common.CasbinEnforcer.LoadPolicy()
			if err != nil {
				return errors.New("permission API updated, but failed to load permission API policy")
			} else {
				return err
			}
		}
	}
	return err
}

// Find a single resource
func (s ApiService) Find(filter map[string]any, data *model.Api) error {
	return common.DB.Where(filter).First(&data).Error
}

// Exist checks whether the resource exists
func (s ApiService) Exist(filter map[string]any) bool {
	var dataObj model.Api
	err := common.DB.Debug().Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Delete multiple resources
func (s ApiService) Delete(ids []uint) error {
	var apis []model.Api
	for _, id := range ids {
		// Get the user by ID
		api := new(model.Api)
		err := s.Find(tools.H{"id": id}, api)
		if err != nil {
			return fmt.Errorf("failed to get API information by ID: %v", err)
		}
		apis = append(apis, *api)
	}

	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.Api{}).Error
	// If deletion succeeds, delete the policy from Casbin
	if err == nil {
		for _, api := range apis {
			policies := common.CasbinEnforcer.GetFilteredPolicy(1, api.Path, api.Method)
			if len(policies) > 0 {
				isRemoved, _ := common.CasbinEnforcer.RemovePolicies(policies)
				if !isRemoved {
					return errors.New("failed to delete permission API")
				}
			}
		}
		// Reload the policies
		err := common.CasbinEnforcer.LoadPolicy()
		if err != nil {
			return errors.New("permission API deleted, but failed to load permission API policy")
		} else {
			return err
		}
	}
	return err
}

// // GetApiTree gets the API tree (grouped by the API Category field)
// func (s ApiService) GetApiTree() ([]*response.ApiTreeDto, error) {
// 	var apiList []*model.Api
// 	err := common.DB.Order("category").Order("created_at").Find(&apiList).Error
// 	// Get all categories
// 	var categoryList []string
// 	for _, api := range apiList {
// 		categoryList = append(categoryList, api.Category)
// 	}
// 	// Get deduplicated categories
// 	categoryUniq := funk.UniqString(categoryList)

// 	apiTree := make([]*response.ApiTreeDto, len(categoryUniq))

// 	for i, category := range categoryUniq {
// 		apiTree[i] = &response.ApiTreeDto{
// 			ID:       -i,
// 			Desc:     category,
// 			Category: category,
// 			Children: nil,
// 		}
// 		for _, api := range apiList {
// 			if category == api.Category {
// 				apiTree[i].Children = append(apiTree[i].Children, api)
// 			}
// 		}
// 	}

// 	return apiTree, err
// }

// GetApisById gets the API list by API ID
func (s ApiService) GetApisById(apiIds []uint) ([]*model.Api, error) {
	var apis []*model.Api
	err := common.DB.Where("id IN (?)", apiIds).Find(&apis).Error
	return apis, err
}
