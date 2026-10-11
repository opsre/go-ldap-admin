package logic

import (
	"fmt"
	"strings"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/client/wechat"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"
	"github.com/gin-gonic/gin"
)

type WeComLogic struct {
}

// Get department information from WeCom
func (d *WeComLogic) SyncWeComDepts(c *gin.Context, req any) (data any, rspError any) {
	// 1. Get all departments
	deptSource, err := wechat.GetAllDepts()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get WeCom department list: %s", err.Error())
		common.Log.Errorf("SyncWeComDepts: %s", errMsg)
		return nil, tools.NewOperationI18nError("sync.dept_list_failed", i18n.Args{"provider": "WeCom", "error": err.Error()})
	}
	depts, err := ConvertDeptData(config.Conf.WeCom.Flag, deptSource)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to convert WeCom department data: %s", err.Error())
		common.Log.Errorf("SyncWeComDepts: %s", errMsg)
		return nil, tools.NewOperationI18nError("sync.dept_convert_failed", i18n.Args{"provider": "WeCom", "error": err.Error()})
	}
	if len(depts) == 0 {
		errMsg := "The department count is 0"
		common.Log.Errorf("SyncWeComDepts: %s", errMsg)
		return nil, tools.NewOperationI18nError("sync.empty_depts", nil)
	}

	// 2. Convert remote data to a tree
	deptTree := GroupListToTree(fmt.Sprintf("%s_1", config.Conf.WeCom.Flag), depts)

	// 3. Create entries from the tree
	err = d.addDepts(deptTree.Children)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create WeCom departments: %s", err.Error())
		common.Log.Errorf("SyncWeComDepts: %s", errMsg)
		return nil, err
	}

	common.Log.Infof("SyncWeComDepts: WeCom department sync succeeded")
	return nil, err
}

// Add the department
func (d WeComLogic) addDepts(depts []*model.Group) error {
	for _, dept := range depts {
		err := d.AddDepts(dept)
		if err != nil {
			errMsg := fmt.Sprintf("DsyncWeComDepts: failed to add department [%s]: %s", dept.GroupName, err.Error())
			common.Log.Errorf("%s", errMsg)
			return tools.NewOperationI18nError("sync.dept_add_failed", i18n.Args{"dept": dept.GroupName, "error": err.Error()})
		}
		if len(dept.Children) != 0 {
			err = d.addDepts(dept.Children)
			if err != nil {
				errMsg := fmt.Sprintf("DsyncWeComDepts: failed to add child departments: %s", err.Error())
				common.Log.Errorf("%s", errMsg)
				return tools.NewOperationI18nError("sync.child_dept_add_failed", i18n.Args{"error": err.Error()})
			}
		}
	}
	return nil
}

// Add department data
func (d WeComLogic) AddDepts(group *model.Group) error {
	// Check whether the department name exists
	parentGroup := new(model.Group)
	err := isql.Group.Find(tools.H{"source_dept_id": group.SourceDeptParentId}, parentGroup)
	if err != nil {
		return tools.NewMySqlI18nError("sync.parent_dept_query_failed", i18n.Args{"error": err.Error()})
	}

	// The group now includes fields linked dynamically by Build. Add the remaining known field values to create the group
	group.Creator = "system"
	group.GroupType = "cn"
	group.ParentId = parentGroup.ID
	group.Source = config.Conf.WeCom.Flag
	group.GroupDN = fmt.Sprintf("cn=%s,%s", group.GroupName, parentGroup.GroupDN)

	if !isql.Group.Exist(tools.H{"group_dn": group.GroupDN}) {
		err = CommonAddGroup(group)
		if err != nil {
			return tools.NewOperationI18nError("sync.add_dept_failed", i18n.Args{"dept": group.GroupName, "error": err.Error()})
		}
	}
	return nil
}

// Start user synchronization based on department information already synced to the database
func (d WeComLogic) SyncWeComUsers(c *gin.Context, req any) (data any, rspError any) {
	// 1. Get the WeCom user list
	staffSource, err := wechat.GetAllUsers()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get WeCom user list: %s", err.Error())
		common.Log.Errorf("SyncWeComUsers: %s", errMsg)
		return nil, tools.NewOperationI18nError("sync.user_list_failed", i18n.Args{"provider": "WeCom", "error": err.Error()})
	}
	staffs, err := ConvertUserData(config.Conf.WeCom.Flag, staffSource)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to convert WeCom user data: %s", err.Error())
		common.Log.Errorf("SyncWeComUsers: %s", errMsg)
		return nil, tools.NewOperationI18nError("sync.user_convert_failed", i18n.Args{"provider": "WeCom", "error": err.Error()})
	}
	if len(staffs) == 0 {
		errMsg := "The user count is 0"
		common.Log.Errorf("SyncWeComUsers: %s", errMsg)
		return nil, tools.NewOperationI18nError("sync.empty_users", nil)
	}
	// 2. Iterate over users and start writing them
	for i, staff := range staffs {
		// Write to the database
		err = d.AddUsers(staff)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to write user [%s]: %s", staff.Username, err.Error())
			common.Log.Errorf("SyncWeComUsers: %s", errMsg)
			return nil, tools.NewOperationI18nError("sync.user_write_failed", i18n.Args{"username": staff.Username, "error": err.Error()})
		}
		common.Log.Infof("SyncWeComUsers: synced user [%s] (%d/%d)", staff.Username, i+1, len(staffs))
	}

	// 3. Get the IDs of departed WeCom users
	// Get all MySQL users sourced from WeCom; those missing remotely have been deleted
	// This can be improved if WeCom later provides an API for listing departed users
	var res []*model.User
	users, err := isql.User.ListAll()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get MySQL user list: %s", err.Error())
		common.Log.Errorf("SyncWeComUsers: %s", errMsg)
		return nil, tools.NewMySqlI18nError("sync.mysql_user_list_failed", i18n.Args{"error": err.Error()})
	}
	for _, user := range users {
		if user.Source != config.Conf.WeCom.Flag {
			continue
		}
		in := true
		for _, staff := range staffs {
			if user.Username == staff.Username {
				in = false
				break
			}
		}
		if in {
			res = append(res, user)
		}
	}
	// 4. Iterate over the IDs and process them
	processedCount := 0
	for _, userTmp := range res {
		user := new(model.User)
		err = isql.User.Find(tools.H{"source_user_id": userTmp.SourceUserId, "status": 1}, user)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to query departed user [%s] in MySQL: %s", userTmp.Username, err.Error())
			common.Log.Errorf("SyncWeComUsers: %s", errMsg)
			return nil, tools.NewMySqlI18nError("sync.leave_user_query_failed", i18n.Args{"username": userTmp.Username, "error": err.Error()})
		}
		// First delete the user from LDAP
		err = ildap.User.Delete(user.UserDN)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to delete departed user [%s] in LDAP: %s", user.Username, err.Error())
			common.Log.Errorf("SyncWeComUsers: %s", errMsg)
			return nil, tools.NewLdapI18nError("sync.leave_user_ldap_delete_failed", i18n.Args{"username": user.Username, "error": err.Error()})
		}
		// Then update the user status in MySQL
		err = isql.User.ChangeStatus(int(user.ID), 2)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to update departed user [%s] status in MySQL: %s", user.Username, err.Error())
			common.Log.Errorf("SyncWeComUsers: %s", errMsg)
			return nil, tools.NewMySqlI18nError("sync.leave_user_status_update_failed", i18n.Args{"username": user.Username, "error": err.Error()})
		}
		processedCount++
		common.Log.Infof("SyncWeComUsers: processed departed user [%s]", user.Username)
	}

	common.Log.Infof("SyncWeComUsers: WeCom user sync complete, synced %d active users and processed %d departed users", len(staffs), processedCount)
	return nil, nil
}

// Add user data
func (d WeComLogic) AddUsers(user *model.User) error {
	// Get the role by role ID
	roles, err := isql.Role.GetRolesByIds([]uint{2})
	if err != nil {
		return tools.NewValidatorI18nError("sync.role_info_failed", i18n.Args{"error": err.Error()})
	}
	user.Creator = "system"
	user.Roles = roles
	user.Password = config.Conf.Ldap.UserInitPassword
	user.Source = config.Conf.WeCom.Flag
	user.UserDN = fmt.Sprintf("uid=%s,%s", user.Username, config.Conf.Ldap.UserDN)

	// Find the user by user_dn; create it if it does not exist
	if !isql.User.Exist(tools.H{"user_dn": user.UserDN}) {
		// Get the groups to add the user to
		groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
		if err != nil {
			return tools.NewMySqlI18nError("user.department_info_failed", i18n.Args{"error": err.Error()})
		}
		var deptTmp string
		for _, group := range groups {
			deptTmp = deptTmp + group.GroupName + ","
		}
		user.Departments = strings.TrimRight(deptTmp, ",")

		// Create the user
		err = CommonAddUser(user, groups)
		if err != nil {
			return tools.NewOperationI18nError("sync.add_user_failed", i18n.Args{"username": user.Username, "error": err.Error()})
		}
	} else {
		// This logic has not been tested in practice. Please report any issues
		if config.Conf.WeCom.IsUpdateSyncd {
			// First get the user information
			oldData := new(model.User)
			err = isql.User.Find(tools.H{"user_dn": user.UserDN}, oldData)
			if err != nil {
				return err
			}
			// Get the groups to add the user to
			groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
			if err != nil {
				return tools.NewMySqlI18nError("user.department_info_failed", i18n.Args{"error": err.Error()})
			}
			var deptTmp string
			for _, group := range groups {
				deptTmp = deptTmp + group.GroupName + ","
			}
			user.Model = oldData.Model
			user.Roles = oldData.Roles
			user.Creator = oldData.Creator
			user.Source = oldData.Source
			user.Password = oldData.Password
			user.UserDN = oldData.UserDN
			user.Departments = strings.TrimRight(deptTmp, ",")

			// Preprocess user information
			if user.Nickname == "" {
				user.Nickname = oldData.Nickname
			}
			if user.GivenName == "" {
				user.GivenName = user.Nickname
			}
			if user.Introduction == "" {
				user.Introduction = user.Nickname
			}
			if user.Mail == "" {
				user.Mail = oldData.Mail
			}
			if user.JobNumber == "" {
				user.JobNumber = oldData.JobNumber
			}
			if user.Departments == "" {
				user.Departments = oldData.Departments
			}
			if user.Position == "" {
				user.Position = oldData.Position
			}
			if user.PostalAddress == "" {
				user.PostalAddress = oldData.PostalAddress
			}
			if user.Mobile == "" {
				user.Mobile = oldData.Mobile
			}
			if err = CommonUpdateUser(oldData, user, tools.StringToSlice(user.DepartmentId, ",")); err != nil {
				return err
			}
		}
	}
	return nil
}
