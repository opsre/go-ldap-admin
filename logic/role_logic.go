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

// Add 添加数据
func (l RoleLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if isql.Role.Exist(tools.H{"name": r.Name}) {
		return nil, tools.NewValidatorI18nError("role.name_exists", nil)
	}

	// 获取当前用户最高角色等级
	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_max_sort_failed", i18n.Args{"error": err.Error()})
	}
	if minSort != 1 {
		return nil, tools.NewValidatorI18nError("role.no_update_permission", nil)
	}
	// 用户不能创建比自己等级高或相同等级的角色
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

	// 创建角色
	err = isql.Role.Add(&role)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.create_failed", i18n.Args{"error": err.Error()})
	}
	return nil, nil
}

// List 数据列表
func (l RoleLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 获取数据列表
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

// Update 更新数据
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

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_max_sort_failed", i18n.Args{"error": err.Error()})
	}

	if minSort != 1 {
		return nil, tools.NewValidatorI18nError("role.no_update_permission", nil)
	}

	// 不能更新比自己角色等级高或相等的角色
	// 根据path中的角色ID获取该角色信息
	roles, _ := isql.Role.GetRolesByIds([]uint{r.ID})
	if len(roles) == 0 {
		return nil, tools.NewMySqlI18nError("role.info_not_found", nil)
	}

	if minSort >= roles[0].Sort {
		return nil, tools.NewValidatorI18nError("role.cannot_update_higher_or_equal", nil)
	}

	// 不能把角色等级更新得比当前用户的等级高
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

	// 更新角色
	err = isql.Role.Update(&role)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.update_failed", i18n.Args{"error": err.Error()})
	}

	// 如果更新成功，且更新了角色的keyword, 则更新casbin中policy
	if r.Keyword != roles[0].Keyword {
		// 获取policy
		rolePolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roles[0].Keyword)
		if len(rolePolicies) == 0 {
			return
		}
		rolePoliciesCopy := make([][]string, 0)
		// 替换keyword
		for _, policy := range rolePolicies {
			policyCopy := make([]string, len(policy))
			copy(policyCopy, policy)
			rolePoliciesCopy = append(rolePoliciesCopy, policyCopy)
			policy[0] = r.Keyword
		}

		//gormadapter未实现UpdatePolicies方法，等gorm更新---
		//isUpdated, _ := common.CasbinEnforcer.UpdatePolicies(rolePoliciesCopy, rolePolicies)
		//if !isUpdated {
		//	response.Fail(c, nil, "更新角色成功，但角色关键字关联的权限接口更新失败！")
		//	return
		//}

		// 这里需要先新增再删除（先删除再增加会出错）
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

	// 更新角色成功处理用户信息缓存有两种做法:（这里使用第二种方法，因为一个角色下用户数量可能很多，第二种方法可以分散数据库压力）
	// 1.可以帮助用户更新拥有该角色的用户信息缓存,使用下面方法
	// err = ur.UpdateUserInfoCacheByRoleId(uint(roleId))
	// 2.直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	isql.User.ClearUserInfoCache()

	return nil, nil
}

// Delete 删除数据
func (l RoleLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 获取当前登陆用户最高等级角色
	minSort, _, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_max_sort_failed", i18n.Args{"error": err.Error()})
	}

	// 获取角色信息
	roles, _ := isql.Role.GetRolesByIds(r.RoleIds)
	if len(roles) == 0 {
		return nil, tools.NewMySqlI18nError("role.no_role_info", nil)
	}

	// 不能删除比自己角色等级高或相等的角色
	for _, role := range roles {
		if minSort >= role.Sort {
			return nil, tools.NewValidatorI18nError("role.cannot_delete_higher_or_equal", nil)
		}
	}

	// 删除角色
	err = isql.Role.Delete(r.RoleIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.delete_failed", i18n.Args{"error": err.Error()})
	}

	// 删除角色成功直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	isql.User.ClearUserInfoCache()
	return nil, nil
}

// GetMenuList 获取角色菜单列表
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

// GetApiList 获取角色接口列表
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

// UpdateMenus 更新角色菜单
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

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_max_sort_failed", i18n.Args{"error": err.Error()})
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限菜单
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			return nil, tools.NewValidatorI18nError("role.cannot_update_higher_or_equal_permissions", nil)
		}
	}

	// 获取当前用户所拥有的权限菜单
	ctxUserMenus, err := isql.Menu.GetUserMenusByUserId(ctxUser.ID)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_user_menu_list_failed", i18n.Args{"error": err.Error()})
	}

	// 获取当前用户所拥有的权限菜单ID
	ctxUserMenusIds := make([]uint, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.ID)
	}

	// 用户需要修改的菜单集合
	reqMenus := make([]*model.Menu, 0)

	// (非管理员)不能把角色的权限菜单设置的比当前用户所拥有的权限菜单多
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
		// 管理员随意设置
		// 根据menuIds查询查询菜单
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

// UpdateApis 更新角色接口
func (l RoleLogic) UpdateApis(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.RoleUpdateApisReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 根据path中的角色ID获取该角色信息
	roles, _ := isql.Role.GetRolesByIds([]uint{r.RoleID})
	if len(roles) == 0 {
		return nil, tools.NewMySqlI18nError("role.no_role_info", nil)
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.current_max_sort_failed", i18n.Args{"error": err.Error()})
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限菜单
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			return nil, tools.NewValidatorI18nError("role.cannot_update_higher_or_equal_permissions", nil)
		}
	}

	// 获取当前用户所拥有的权限接口
	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy := common.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	// 得到path中的角色ID对应角色能够设置的权限接口集合
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}

	// 根据apiID获取接口详情
	apis, err := isql.Api.GetApisById(r.ApiIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.api_info_failed", nil)
	}
	// 生成前端想要设置的角色policies
	reqRolePolicies := make([][]string, 0)
	for _, api := range apis {
		reqRolePolicies = append(reqRolePolicies, []string{
			roles[0].Keyword, api.Path, api.Method,
		})
	}

	// (非管理员)不能把角色的权限接口设置的比当前用户所拥有的权限接口多
	if minSort != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				return nil, tools.NewValidatorI18nError("role.no_api_permission", i18n.Args{"path": reqPolicy[1], "method": reqPolicy[2]})
			}
		}
	}

	// 更新角色的权限接口
	err = isql.Role.UpdateRoleApis(roles[0].Keyword, reqRolePolicies)
	if err != nil {
		return nil, tools.NewMySqlI18nError("role.update_apis_failed", nil)
	}
	return nil, nil
}
