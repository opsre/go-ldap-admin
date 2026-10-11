package common

import (
	"errors"
	"fmt"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/tools"

	"github.com/thoas/go-funk"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Initialize MySQL data
func InitData() {
	// Whether to initialize data
	if !config.Conf.System.InitData {
		return
	}

	// 1. Insert role data
	newRoles := make([]*model.Role, 0)
	roles := []*model.Role{
		{
			Model:   gorm.Model{ID: 1},
			Name:    "Administrator",
			Keyword: "admin",
			Remark:  "",
			Sort:    1,
			Status:  1,
			Creator: "system",
		},
		{
			Model:   gorm.Model{ID: 2},
			Name:    "User",
			Keyword: "user",
			Remark:  "",
			Sort:    3,
			Status:  1,
			Creator: "system",
		},
		{
			Model:   gorm.Model{ID: 3},
			Name:    "Guest",
			Keyword: "guest",
			Remark:  "",
			Sort:    5,
			Status:  1,
			Creator: "system",
		},
	}

	for _, role := range roles {
		err := DB.First(&role, role.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRoles = append(newRoles, role)
		}
	}

	if len(newRoles) > 0 {
		err := DB.Create(&newRoles).Error
		if err != nil {
			Log.Errorf("Failed to write system role data: %v", err)
		}
	}

	// 2. Insert menus
	newMenus := make([]model.Menu, 0)
	var uint0 uint = 0
	var uint1 uint = 1
	var uint4 uint = 5
	var uint8 uint = 9
	componentStr := "component"
	systemRoleStr := "/system/role"
	personnelManageStr := "/personnel/user"
	userStr := "user"
	peopleStr := "people"
	groupStr := "peoples"
	fieldRelationStr := "el-icon-s-tools"
	roleStr := "eye-open"
	treeTableStr := "tree-table"
	treeStr := "tree"
	exampleStr := "example"
	logOperationStr := "/log/operation-log"
	documentationStr := "documentation"
	menus := []model.Menu{
		{
			Model:     gorm.Model{ID: 1},
			Name:      "UserManage",
			Title:     "Personnel",
			Icon:      userStr,
			Path:      "/personnel",
			Component: "Layout",
			Redirect:  personnelManageStr,
			Sort:      5,
			ParentId:  uint0,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 2},
			Name:      "User",
			Title:     "Users",
			Icon:      peopleStr,
			Path:      "user",
			Component: "/personnel/user/index",
			Sort:      6,
			ParentId:  uint1,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 3},
			Name:      "Group",
			Title:     "Groups",
			Icon:      groupStr,
			Path:      "group",
			Component: "/personnel/group/index",
			Sort:      7,
			ParentId:  uint1,
			NoCache:   1,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 4},
			Name:      "FieldRelation",
			Title:     "Field Relations",
			Icon:      fieldRelationStr,
			Path:      "fieldRelation",
			Component: "/personnel/fieldRelation/index",
			Sort:      8,
			ParentId:  uint1,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 5},
			Name:      "System",
			Title:     "System",
			Icon:      componentStr,
			Path:      "/system",
			Component: "Layout",
			Redirect:  systemRoleStr,
			Sort:      9,
			ParentId:  uint0,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 6},
			Name:      "Role",
			Title:     "Roles",
			Icon:      roleStr,
			Path:      "role",
			Component: "/system/role/index",
			Sort:      10,
			ParentId:  uint4,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 7},
			Name:      "Menu",
			Title:     "Menus",
			Icon:      treeTableStr,
			Path:      "menu",
			Component: "/system/menu/index",
			Sort:      13,
			ParentId:  uint4,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 8},
			Name:      "Api",
			Title:     "APIs",
			Icon:      treeStr,
			Path:      "api",
			Component: "/system/api/index",
			Sort:      14,
			ParentId:  uint4,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 9},
			Name:      "Log",
			Title:     "Logs",
			Icon:      exampleStr,
			Path:      "/log",
			Component: "Layout",
			Redirect:  logOperationStr,
			Sort:      20,
			ParentId:  uint0,
			Roles:     roles[:2],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 10},
			Name:      "OperationLog",
			Title:     "Operation Logs",
			Icon:      documentationStr,
			Path:      "operation-log",
			Component: "/log/operation-log/index",
			Sort:      21,
			ParentId:  uint8,
			Roles:     roles[:2],
			Creator:   "system",
		},
	}
	for _, menu := range menus {
		err := DB.First(&menu, menu.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newMenus = append(newMenus, menu)
		}
	}
	if len(newMenus) > 0 {
		err := DB.Create(&newMenus).Error
		if err != nil {
			Log.Errorf("Failed to write system menu data: %v", err)
		}
	}

	// 3. Insert users
	newUsers := make([]*model.User, 0)
	users := []*model.User{
		{
			Model:         gorm.Model{ID: 1},
			Username:      "admin",
			Password:      tools.NewGenPasswd(config.Conf.Ldap.AdminPass),
			Nickname:      "Administrator",
			GivenName:     "Go LDAP Admin",
			Mail:          "admin@" + config.Conf.Ldap.DefaultEmailSuffix,
			JobNumber:     "0000",
			Mobile:        "18888888888",
			Avatar:        "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			PostalAddress: "Earth",
			Departments:   "R&D Center",
			Position:      "Employee",
			Introduction:  "Go LDAP Admin administrator",
			Status:        1,
			Creator:       "system",
			Roles:         roles[:1],
			UserDN:        config.Conf.Ldap.AdminDN,
		},
	}

	for _, user := range users {
		err := DB.First(&user, user.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) > 0 {
		err := DB.Create(&newUsers).Error
		if err != nil {
			Log.Errorf("Failed to write user data: %v", err)
		}
	}

	// 4. Insert APIs
	apis := []model.Api{
		{
			Method:   "GET",
			Path:     "/base/ping",
			Category: "base",
			Remark:   "Test endpoint",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/base/encryptpwd",
			Category: "base",
			Remark:   "Encrypt password",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/base/decryptpwd",
			Category: "base",
			Remark:   "Decrypt password",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/base/config",
			Category: "base",
			Remark:   "Get system configuration",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/base/version",
			Category: "base",
			Remark:   "Get system version",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/base/dashboard",
			Category: "base",
			Remark:   "Get dashboard data",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/base/login",
			Category: "base",
			Remark:   "User login",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/base/logout",
			Category: "base",
			Remark:   "User logout",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/base/refreshToken",
			Category: "base",
			Remark:   "Refresh JWT token",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/base/sendcode",
			Category: "base",
			Remark:   "Send verification code by email",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/base/changePwd",
			Category: "base",
			Remark:   "Change password by email",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/user/info",
			Category: "user",
			Remark:   "Get current user information",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/user/list",
			Category: "user",
			Remark:   "Get user list",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/changePwd",
			Category: "user",
			Remark:   "Change user password",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/resetPassword",
			Category: "user",
			Remark:   "Reset user password",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/add",
			Category: "user",
			Remark:   "Create user",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/update",
			Category: "user",
			Remark:   "Update user",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/delete",
			Category: "user",
			Remark:   "Delete users",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/changeUserStatus",
			Category: "user",
			Remark:   "Change user status",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/syncDingTalkUsers",
			Category: "user",
			Remark:   "Sync users from DingTalk",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/syncWeComUsers",
			Category: "user",
			Remark:   "Sync users from WeCom",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/syncFeiShuUsers",
			Category: "user",
			Remark:   "Sync users from Feishu",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/syncOpenLdapUsers",
			Category: "user",
			Remark:   "Sync users from OpenLDAP",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/syncSqlUsers",
			Category: "user",
			Remark:   "Sync database users to LDAP",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/group/list",
			Category: "group",
			Remark:   "Get group list",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/group/tree",
			Category: "group",
			Remark:   "Get group tree",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/group/add",
			Category: "group",
			Remark:   "Create group",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/group/update",
			Category: "group",
			Remark:   "Update group",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/group/delete",
			Category: "group",
			Remark:   "Delete groups",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/group/adduser",
			Category: "group",
			Remark:   "Add user to group",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/group/removeuser",
			Category: "group",
			Remark:   "Remove user from group",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/group/useringroup",
			Category: "group",
			Remark:   "Get users in group",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/group/usernoingroup",
			Category: "group",
			Remark:   "Get users outside group",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/group/syncDingTalkDepts",
			Category: "group",
			Remark:   "Sync departments from DingTalk",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/group/syncWeComDepts",
			Category: "group",
			Remark:   "Sync departments from WeCom",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/group/syncFeiShuDepts",
			Category: "group",
			Remark:   "Sync departments from Feishu",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/group/syncOpenLdapDepts",
			Category: "group",
			Remark:   "Sync departments from OpenLDAP",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/group/syncSqlGroups",
			Category: "group",
			Remark:   "Sync database groups to LDAP",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/role/list",
			Category: "role",
			Remark:   "Get role list",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/role/add",
			Category: "role",
			Remark:   "Create role",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/role/update",
			Category: "role",
			Remark:   "Update role",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/role/getmenulist",
			Category: "role",
			Remark:   "Get role menus",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/role/updatemenus",
			Category: "role",
			Remark:   "Update role menus",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/role/getapilist",
			Category: "role",
			Remark:   "Get role APIs",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/role/updateapis",
			Category: "role",
			Remark:   "Update role APIs",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/role/delete",
			Category: "role",
			Remark:   "Delete roles",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/menu/tree",
			Category: "menu",
			Remark:   "Get menu tree",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/menu/access/tree",
			Category: "menu",
			Remark:   "Get user menu tree",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/menu/add",
			Category: "menu",
			Remark:   "Create menu",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/menu/update",
			Category: "menu",
			Remark:   "Update menu",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/menu/delete",
			Category: "menu",
			Remark:   "Delete menus",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/api/list",
			Category: "api",
			Remark:   "Get API list",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/api/tree",
			Category: "api",
			Remark:   "Get API tree",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/api/add",
			Category: "api",
			Remark:   "Create API",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/api/update",
			Category: "api",
			Remark:   "Update API",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/api/delete",
			Category: "api",
			Remark:   "Delete APIs",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/fieldrelation/list",
			Category: "fieldrelation",
			Remark:   "Get field relation list",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/fieldrelation/add",
			Category: "fieldrelation",
			Remark:   "Create field relation",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/fieldrelation/update",
			Category: "fieldrelation",
			Remark:   "Update field relation",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/fieldrelation/delete",
			Category: "fieldrelation",
			Remark:   "Delete field relations",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/log/operation/list",
			Category: "log",
			Remark:   "Get operation log list",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/log/operation/delete",
			Category: "log",
			Remark:   "Delete operation logs",
			Creator:  "system",
		},
		{
			Method:   "DELETE",
			Path:     "/log/operation/clean",
			Category: "log",
			Remark:   "Clear operation logs",
			Creator:  "system",
		},
	}

	// 5. Assign roles to menus
	newApi := make([]model.Api, 0)
	newRoleCasbin := make([]model.RoleCasbin, 0)
	for i, api := range apis {
		api.ID = uint(i + 1)
		err := DB.First(&api, api.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newApi = append(newApi, api)

			// Admins have all API permissions
			newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
				Keyword: roles[0].Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})

			// Non-admins have basic permissions
			basePaths := []string{
				"/base/login",
				"/base/logout",
				"/base/refreshToken",
				"/base/sendcode",
				"/base/changePwd",
				"/base/dashboard",
				"/user/info",
				"/user/changePwd",
				"/menu/access/tree",
				"/log/operation/list",
			}

			if funk.ContainsString(basePaths, api.Path) {
				newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
					Keyword: roles[1].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
		}
	}

	if len(newApi) > 0 {
		if err := DB.Create(&newApi).Error; err != nil {
			Log.Errorf("Failed to write API data: %v", err)
		}
	}

	if len(newRoleCasbin) > 0 {
		rules := make([][]string, 0)
		for _, c := range newRoleCasbin {
			rules = append(rules, []string{
				c.Keyword, c.Path, c.Method,
			})
		}
		isAdd, err := CasbinEnforcer.AddPolicies(rules)
		if !isAdd {
			Log.Errorf("Failed to write Casbin data: %v", err)
		}
	}

	// 6. Insert groups
	newGroups := make([]model.Group, 0)
	groups := []model.Group{
		{
			Model:              gorm.Model{ID: 1},
			GroupName:          "root",
			Remark:             "Base",
			Creator:            "system",
			GroupType:          "",
			ParentId:           0,
			SourceDeptId:       "0",
			Source:             "openldap",
			SourceDeptParentId: "0",
			GroupDN:            config.Conf.Ldap.BaseDN,
		},
		{
			Model:              gorm.Model{ID: 2},
			GroupName:          config.Conf.DingTalk.Flag + "root",
			Remark:             "DingTalk root department",
			Creator:            "system",
			GroupType:          "ou",
			ParentId:           1,
			SourceDeptId:       fmt.Sprintf("%s_%d", config.Conf.DingTalk.Flag, 1),
			Source:             config.Conf.DingTalk.Flag,
			SourceDeptParentId: fmt.Sprintf("%s_%d", config.Conf.DingTalk.Flag, 0),
			GroupDN:            fmt.Sprintf("ou=%s,%s", config.Conf.DingTalk.Flag+"root", config.Conf.Ldap.BaseDN),
		},
		{
			Model:              gorm.Model{ID: 3},
			GroupName:          "wecomroot",
			Remark:             "WeCom root department",
			Creator:            "system",
			GroupType:          "ou",
			ParentId:           1,
			SourceDeptId:       fmt.Sprintf("%s_%d", config.Conf.WeCom.Flag, 1),
			Source:             config.Conf.WeCom.Flag,
			SourceDeptParentId: fmt.Sprintf("%s_%d", config.Conf.WeCom.Flag, 0),
			GroupDN:            fmt.Sprintf("ou=%s,%s", config.Conf.WeCom.Flag+"root", config.Conf.Ldap.BaseDN),
		},
		{
			Model:              gorm.Model{ID: 4},
			GroupName:          config.Conf.FeiShu.Flag + "root",
			Remark:             "Feishu root department",
			Creator:            "system",
			GroupType:          "ou",
			ParentId:           1,
			SourceDeptId:       fmt.Sprintf("%s_%d", config.Conf.FeiShu.Flag, 0),
			Source:             config.Conf.FeiShu.Flag,
			SourceDeptParentId: fmt.Sprintf("%s_%d", config.Conf.FeiShu.Flag, 0),
			GroupDN:            fmt.Sprintf("ou=%s,%s", config.Conf.FeiShu.Flag+"root", config.Conf.Ldap.BaseDN),
		},
		{
			Model:              gorm.Model{ID: 5},
			GroupName:          "group",
			Remark:             "Default group",
			Creator:            "system",
			GroupType:          "cn",
			ParentId:           1,
			SourceDeptId:       "platform_0",
			Source:             "platform",
			SourceDeptParentId: "openldap_1",
			GroupDN:            fmt.Sprintf("cn=%s,%s", "group", config.Conf.Ldap.BaseDN),
		},
	}

	for _, group := range groups {
		err := DB.First(&group, group.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newGroups = append(newGroups, group)
		}
	}
	if len(newGroups) > 0 {
		err := DB.Create(&newGroups).Error
		if err != nil {
			Log.Errorf("Failed to write group data: %v", err)
		}
	}

	// 7. Insert field relations
	filedRelation := []model.FieldRelation{
		{
			Flag:       "dingtalk_group",
			Attributes: datatypes.JSON(`{"groupName":"custom_name_pinyin","remark":"name","sourceDeptId":"id","sourceDeptParentId":"parentid"}`),
		},
		{
			Flag:       "dingtalk_user",
			Attributes: datatypes.JSON(`{"avatar":"avatar","givenName":"name","introduction":"remark","jobNumber":"job_number","mail":"email","mobile":"mobile","nickname":"name","position":"title","postalAddress":"work_place","sourceUnionId":"unionid","sourceUserId":"userid","username":"custom_name_pinyin"}`),
		},
		{
			Flag:       "feishu_group",
			Attributes: datatypes.JSON(`{"groupName":"custom_name_pinyin","remark":"name","sourceDeptId":"open_department_id","sourceDeptParentId":"parent_department_id"}`),
		},
		{
			Flag:       "feishu_user",
			Attributes: datatypes.JSON(`{"avatar":"avatar","givenName":"name","introduction":"name","jobNumber":"employee_no","mail":"email","mobile":"mobile","nickname":"name","position":"job_title","postalAddress":"work_station","sourceUnionId":"union_id","sourceUserId":"user_id","username":"custom_name_pinyin"}`),
		},
		{
			Flag:       "wecom_group",
			Attributes: datatypes.JSON(`{"groupName":"custom_name_pinyin","remark":"name","sourceDeptId":"id","sourceDeptParentId":"parentid"}`),
		},
		{
			Flag:       "wecom_user",
			Attributes: datatypes.JSON(`{"avatar":"avatar","givenName":"alias","introduction":"name","jobNumber":"mobile","mail":"email","mobile":"mobile","nickname":"name","position":"external_position","postalAddress":"address","sourceUnionId":"userid","sourceUserId":"userid","username":"custom_name_pinyin"}`),
		},
	}

	newFieldRelations := make([]model.FieldRelation, 0)
	for i, newFieldRelation := range filedRelation {
		newFieldRelation.ID = uint(i + 1)
		err := DB.First(&newFieldRelation, newFieldRelation.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newFieldRelations = append(newFieldRelations, newFieldRelation)
		}
	}

	if len(newFieldRelations) > 0 {
		if err := DB.Create(&newFieldRelations).Error; err != nil {
			Log.Errorf("Failed to write field relation data: %v", err)
		}
	}
}
