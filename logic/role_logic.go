package logic

import (
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/model/response"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type RoleLogic struct{}

// Add data
func (l RoleLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if isql.Role.Exist(tools.H{"name": r.Name}) {
		return nil, tools.NewValidatorI18nError("role.name_exists", nil)
	}

	// Get the current user highest role level
	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_max_sort_failed", i18n.Args{"error": err.Error()})
	}
	if minSort != 1 {
		return nil, tools.NewValidatorI18nError("role.no_update_permission", nil)
	}
	// Users cannot create roles at or above their own level
	if minSort >= r.Sort {
		return nil, tools.NewValidatorI18nError("role.cannot_create_higher_or_equal", nil)
	}

	role := model.Role{
		Name:    r.Name,
		Keyword: r.Keyword,
		Remark:  r.Remark,
		Status:  r.Status,
		Sort:    r.Sort,
		Creator: ctxUser.Username,
	}

	// Create the role
	err = isql.Role.Add(&role)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.create_failed", i18n.Args{"error": err.Error()})
	}
	return nil, nil
}

// List data
func (l RoleLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// Get the data list
	roles, err := isql.Role.List(r)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.list_failed", i18n.Args{"error": err.Error()})
	}

	count, err := isql.Role.Count()
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.count_failed", nil)
	}

	localizeRoles(c, roles)
	rets := make([]model.Role, 0)
	for _, role := range roles {
		rets = append(rets, *role)
	}

	return response.RoleListRsp{
		Total: count,
		Roles: rets,
	}, nil
}

// Update data
func (l RoleLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.ID}
	if !isql.Role.Exist(filter) {
		return nil, tools.NewValidatorI18nError("role.not_found", nil)
	}

	// The current user and their minimum role sort value (highest role level)
	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_max_sort_failed", i18n.Args{"error": err.Error()})
	}

	if minSort != 1 {
		return nil, tools.NewValidatorI18nError("role.no_update_permission", nil)
	}

	// Cannot update roles at or above your own level
	// Get the role information using the role ID in the path
	roles, _ := isql.Role.GetRolesByIds([]uint{r.ID})
	if len(roles) == 0 {
		return nil, tools.NewMySqlI18nError("role.info_not_found", nil)
	}

	if minSort >= roles[0].Sort {
		return nil, tools.NewValidatorI18nError("role.cannot_update_higher_or_equal", nil)
	}

	// Cannot raise the role level above the current user level
	if minSort >= r.Sort {
		return nil, tools.NewValidatorI18nError("role.cannot_set_higher_or_equal", nil)
	}
	oldData := new(model.Role)
	err = isql.Role.Find(filter, oldData)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}
	role := model.Role{
		Model:   oldData.Model,
		Name:    r.Name,
		Keyword: r.Keyword,
		Remark:  r.Remark,
		Status:  r.Status,
		Sort:    r.Sort,
		Creator: ctxUser.Username,
	}

	// Update the role
	err = isql.Role.Update(&role)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.update_failed", i18n.Args{"error": err.Error()})
	}

	// If the update succeeds and the role keyword changed, update the policy in Casbin
	if r.Keyword != roles[0].Keyword {
		// Get the policy
		rolePolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roles[0].Keyword)
		if len(rolePolicies) == 0 {
			return
		}
		rolePoliciesCopy := make([][]string, 0)
		// Replace the keyword
		for _, policy := range rolePolicies {
			policyCopy := make([]string, len(policy))
			copy(policyCopy, policy)
			rolePoliciesCopy = append(rolePoliciesCopy, policyCopy)
			policy[0] = r.Keyword
		}

		//gormadapter does not implement UpdatePolicies; wait for a GORM update---
		//isUpdated, _ := common.CasbinEnforcer.UpdatePolicies(rolePoliciesCopy, rolePolicies)
		//if !isUpdated {
		//	response.Fail(c, nil, "更新角色成功，但角色关键字关联的权限接口更新失败！")
		//	return
		//}

		// Add first and then delete here (deleting first and then adding causes an error)
		isAdded, _ := common.CasbinEnforcer.AddPolicies(rolePolicies)
		if !isAdded {
			return nil, tools.NewOperationI18nError("role.keyword_policy_update_failed", nil)
		}
		isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rolePoliciesCopy)
		if !isRemoved {
			return nil, tools.NewOperationI18nError("role.keyword_policy_update_failed", nil)
		}
		err := common.CasbinEnforcer.LoadPolicy()
		if err != nil {
			return nil, tools.NewOperationI18nError("role.keyword_policy_load_failed", nil)
		}

	}

	// There are two ways to handle user information caches after a successful role update. This uses the second to spread database load when a role has many users
	// 1. Update the information caches for users with this role using the method below
	// err = ur.UpdateUserInfoCacheByRoleId(uint(roleId))
	// 2. Clear the caches and let active users cache their latest information again
	isql.User.ClearUserInfoCache()

	return nil, nil
}

// Delete data
func (l RoleLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// Get the highest-level role of the currently logged-in user
	minSort, _, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_max_sort_failed", i18n.Args{"error": err.Error()})
	}

	// Get the role information
	roles, _ := isql.Role.GetRolesByIds(r.RoleIds)
	if len(roles) == 0 {
		return nil, tools.NewMySqlI18nError("role.no_role_info", nil)
	}

	// Cannot delete roles at or above your own level
	for _, role := range roles {
		if minSort >= role.Sort {
			return nil, tools.NewValidatorI18nError("role.cannot_delete_higher_or_equal", nil)
		}
	}

	// Delete the role
	err = isql.Role.Delete(r.RoleIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.delete_failed", i18n.Args{"error": err.Error()})
	}

	// After deleting the role, clear caches so active users cache their latest information again
	isql.User.ClearUserInfoCache()
	return nil, nil
}

// GetMenuList gets the role menu list
func (l RoleLogic) GetMenuList(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleGetMenuListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	menus, err := isql.Role.GetRoleMenusById(r.RoleID)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.menu_list_failed", i18n.Args{"error": err.Error()})
	}
	localizeMenus(c, menus)
	return menus, nil
}

// GetApiList gets the role API list
func (l RoleLogic) GetApiList(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleGetApiListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	role := new(model.Role)
	err := isql.Role.Find(tools.H{"id": r.RoleID}, role)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.resource_failed", i18n.Args{"error": err.Error()})
	}
	localizeRole(c, role)

	policies := common.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)

	apis, err := isql.Api.ListAll()
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.resource_list_failed", i18n.Args{"error": err.Error()})
	}
	accessApis := make([]*model.Api, 0)

	for _, policy := range policies {
		path := policy[1]
		method := policy[2]
		for _, api := range apis {
			if path == api.Path && method == api.Method {
				accessApis = append(accessApis, api)
				break
			}
		}
	}
	localizeApis(c, accessApis)

	return accessApis, nil
}

// UpdateMenus updates role menus
func (l RoleLogic) UpdateMenus(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleUpdateMenusReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	roles, _ := isql.Role.GetRolesByIds([]uint{r.RoleID})
	if len(roles) == 0 {
		return nil, tools.NewMySqlI18nError("role.no_role_info", nil)
	}

	// The current user and their minimum role sort value (highest role level)
	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_max_sort_failed", i18n.Args{"error": err.Error()})
	}

	// (Non-admin) cannot update permission menus for roles at or above their own level
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			return nil, tools.NewValidatorI18nError("role.cannot_update_higher_or_equal_permissions", nil)
		}
	}

	// Get the permission menus owned by the current user
	ctxUserMenus, err := isql.Menu.GetUserMenusByUserId(ctxUser.ID)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_user_menu_list_failed", i18n.Args{"error": err.Error()})
	}

	// Get the IDs of permission menus owned by the current user
	ctxUserMenusIds := make([]uint, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.ID)
	}

	// The set of menus the user needs to modify
	reqMenus := make([]*model.Menu, 0)

	// (Non-admin) cannot assign a role more permission menus than the current user has
	if minSort != 1 {
		for _, id := range r.MenuIds {
			if !funk.Contains(ctxUserMenusIds, id) {
				return nil, tools.NewValidatorI18nError("role.no_menu_permission", i18n.Args{"id": id})
			}
		}

		for _, id := range r.MenuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.ID {
					reqMenus = append(reqMenus, menu)
					break
				}
			}
		}
	} else {
		// Admins can set these freely
		// Query menus by menuIds
		menus, err := isql.Menu.List()
		if err != nil {
			return nil, tools.NewValidatorI18nError("role.menu_all_list_failed", i18n.Args{"error": err.Error()})
		}
		for _, menuId := range r.MenuIds {
			for _, menu := range menus {
				if menuId == menu.ID {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}

	roles[0].Menus = reqMenus

	err = isql.Role.UpdateRoleMenus(roles[0])
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.update_menus_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// UpdateApis updates role APIs
func (l RoleLogic) UpdateApis(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleUpdateApisReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// Get the role information using the role ID in the path
	roles, _ := isql.Role.GetRolesByIds([]uint{r.RoleID})
	if len(roles) == 0 {
		return nil, tools.NewMySqlI18nError("role.no_role_info", nil)
	}

	// The current user and their minimum role sort value (highest role level)
	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_max_sort_failed", i18n.Args{"error": err.Error()})
	}

	// (Non-admin) cannot update permission menus for roles at or above their own level
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			return nil, tools.NewValidatorI18nError("role.cannot_update_higher_or_equal_permissions", nil)
		}
	}

	// Get the permission APIs owned by the current user
	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy := common.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	// Get the set of permission APIs assignable to the role ID in the path
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}

	// Get API details by apiID
	apis, err := isql.Api.GetApisById(r.ApiIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.api_info_failed", nil)
	}
	// Generate the role policies requested by the frontend
	reqRolePolicies := make([][]string, 0)
	for _, api := range apis {
		reqRolePolicies = append(reqRolePolicies, []string{
			roles[0].Keyword, api.Path, api.Method,
		})
	}

	// (Non-admin) cannot assign a role more permission APIs than the current user has
	if minSort != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				return nil, tools.NewValidatorI18nError("role.no_api_permission", i18n.Args{"path": reqPolicy[1], "method": reqPolicy[2]})
			}
		}
	}

	// Update the role permission APIs
	err = isql.Role.UpdateRoleApis(roles[0].Keyword, reqRolePolicies)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.update_apis_failed", nil)
	}
	return nil, nil
}
