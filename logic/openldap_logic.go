package logic

import (
	"fmt"
	"strings"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/client/openldap"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/isql"
	"github.com/gin-gonic/gin"
)

type OpenLdapLogic struct {
}

// Get department information from LDAP
func (d *OpenLdapLogic) SyncOpenLdapDepts(c *gin.Context, req any) (data any, rspError any) {
	// 1. Get all departments
	depts, err := openldap.GetAllDepts()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get OpenLDAP department list: %s", err.Error())
		common.Log.Errorf("SyncOpenLdapDepts: %s", errMsg)
		return nil, tools.NewOperationI18nError("sync.dept_list_failed", i18n.Args{"provider": "OpenLDAP", "error": err.Error()})
	}
	if len(depts) == 0 {
		errMsg := "The department count is 0"
		common.Log.Errorf("SyncOpenLdapDepts: %s", errMsg)
		return nil, tools.NewOperationI18nError("sync.empty_depts", nil)
	}
	groups := make([]*model.Group, 0)
	for _, dept := range depts {
		groups = append(groups, &model.Group{
			GroupName:          dept.Name,
			Remark:             dept.Remark,
			SourceDeptId:       dept.Id,
			SourceDeptParentId: dept.ParentId,
			GroupDN:            dept.DN,
		})
	}
	// 2. Convert remote data to a tree
	deptTree := GroupListToTree("0", groups)

	// 3. Create entries from the tree
	err = d.addDepts(deptTree.Children)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create OpenLDAP departments: %s", err.Error())
		common.Log.Errorf("SyncOpenLdapDepts: %s", errMsg)
		return nil, err
	}

	common.Log.Infof("SyncOpenLdapDepts: OpenLDAP department sync succeeded")
	return nil, err
}

// Add the department
func (d OpenLdapLogic) addDepts(depts []*model.Group) error {
	for _, dept := range depts {
		err := d.AddDepts(dept)
		if err != nil {
			errMsg := fmt.Sprintf("DsyncOpenLdapDepts: failed to add department [%s]: %s", dept.GroupName, err.Error())
			common.Log.Errorf("%s", errMsg)
			return tools.NewOperationI18nError("sync.dept_add_failed", i18n.Args{"dept": dept.GroupName, "error": err.Error()})
		}
		if len(dept.Children) != 0 {
			err = d.addDepts(dept.Children)
			if err != nil {
				errMsg := fmt.Sprintf("DsyncOpenLdapDepts: failed to add child departments: %s", err.Error())
				common.Log.Errorf("%s", errMsg)
				return tools.NewOperationI18nError("sync.child_dept_add_failed", i18n.Args{"error": err.Error()})
			}
		}
	}
	return nil
}

// Add department data
func (d OpenLdapLogic) AddDepts(group *model.Group) error {
	// Check whether the department name exists. Use the unique LDAP DN to avoid incomplete synchronization
	if !isql.Group.Exist(tools.H{"group_dn": group.GroupDN}) {
		// The group now includes fields linked dynamically by Build. Add the remaining known field values to create the group
		group.Creator = "system"
		group.GroupType = strings.Split(strings.Split(group.GroupDN, ",")[0], "=")[0]
		parentid, err := d.getParentGroupID(group)
		if err != nil {
			return err
		}
		group.ParentId = parentid
		group.Source = "openldap"
		err = isql.Group.Add(group)
		if err != nil {
			return err
		}
	}
	return nil

}

// Add department data
func (d OpenLdapLogic) getParentGroupID(group *model.Group) (id uint, err error) {
	switch group.SourceDeptParentId {
	case "dingtalkroot":
		group.SourceDeptParentId = "dingtalk_1"
	case "feishuroot":
		group.SourceDeptParentId = "feishu_0"
	case "wecomroot":
		group.SourceDeptParentId = "wecom_1"
	}
	parentGroup := new(model.Group)
	err = isql.Group.Find(tools.H{"source_dept_id": group.SourceDeptParentId}, parentGroup)
	if err != nil {
		return id, tools.NewMySqlI18nError("sync.parent_dept_query_failed", i18n.Args{"error": fmt.Sprintf("%s,%s", err.Error(), group.GroupName)})
	}
	return parentGroup.ID, nil
}

// Start user synchronization based on department information already synced to the database
func (d OpenLdapLogic) SyncOpenLdapUsers(c *gin.Context, req any) (data any, rspError any) {
	// 1. Get the LDAP user list
	staffs, err := openldap.GetAllUsers()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get OpenLDAP user list: %s", err.Error())
		common.Log.Errorf("SyncOpenLdapUsers: %s", errMsg)
		return nil, tools.NewOperationI18nError("sync.user_list_failed", i18n.Args{"provider": "OpenLDAP", "error": err.Error()})
	}
	if len(staffs) == 0 {
		errMsg := "The user count is 0"
		common.Log.Errorf("SyncOpenLdapUsers: %s", errMsg)
		return nil, tools.NewOperationI18nError("sync.empty_users", nil)
	}
	// 2. Iterate over users and start writing them
	for i, staff := range staffs {
		groupIds, err := isql.Group.DeptIdsToGroupIds(staff.DepartmentIds)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to convert user [%s] department IDs to internal department IDs: %s", staff.Name, err.Error())
			common.Log.Errorf("SyncOpenLdapUsers: %s", errMsg)
			return nil, tools.NewMySqlI18nError("sync.user_dept_ids_convert_failed", i18n.Args{"username": staff.Name, "error": err.Error()})
		}
		// Get the role by role ID
		roles, err := isql.Role.GetRolesByIds([]uint{2})
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get role information for user [%s]: %s", staff.Name, err.Error())
			common.Log.Errorf("SyncOpenLdapUsers: %s", errMsg)
			return nil, tools.NewValidatorI18nError("sync.user_role_info_failed", i18n.Args{"username": staff.Name, "error": err.Error()})
		}
		// Write to the database
		err = d.AddUsers(&model.User{
			Username:      staff.Name,
			Nickname:      staff.DisplayName,
			GivenName:     staff.GivenName,
			Mail:          staff.Mail,
			JobNumber:     staff.EmployeeNumber,
			Mobile:        staff.Mobile,
			PostalAddress: staff.PostalAddress,
			Departments:   staff.BusinessCategory,
			Position:      staff.DepartmentNumber,
			Introduction:  staff.CN,
			Creator:       "system",
			Source:        "openldap",
			DepartmentId:  tools.SliceToString(groupIds, ","),
			SourceUserId:  staff.Name,
			SourceUnionId: staff.Name,
			Roles:         roles,
			UserDN:        staff.DN,
		})
		if err != nil {
			errMsg := fmt.Sprintf("Failed to write user [%s]: %s", staff.Name, err.Error())
			common.Log.Errorf("SyncOpenLdapUsers: %s", errMsg)
			return nil, tools.NewOperationI18nError("sync.user_write_failed", i18n.Args{"username": staff.Name, "error": err.Error()})
		}
		common.Log.Infof("SyncOpenLdapUsers: synced user [%s] (%d/%d)", staff.Name, i+1, len(staffs))
	}

	common.Log.Infof("SyncOpenLdapUsers: OpenLDAP user sync complete, synced %d users", len(staffs))
	return nil, nil
}

// Add user data
func (d OpenLdapLogic) AddUsers(user *model.User) error {
	// Find the user by user_dn; create it if it does not exist
	if !isql.User.Exist(tools.H{"user_dn": user.UserDN}) {
		if user.Departments == "" {
			user.Departments = "Default: R&D Center"
		}
		if user.GivenName == "" {
			user.GivenName = user.Nickname
		}
		if user.PostalAddress == "" {
			user.PostalAddress = "Default: Earth"
		}
		if user.Position == "" {
			user.Position = "Default: Engineer"
		}
		if user.Introduction == "" {
			user.Introduction = user.Nickname
		}
		if user.JobNumber == "" {
			user.JobNumber = "N/A"
		}
		// First add the user to MySQL
		err := isql.User.Add(user)
		if err != nil {
			common.Log.Errorf("AddUsers: failed to create user in MySQL: %s", err.Error())
			return tools.NewMySqlI18nError("openldap_sync.user_mysql_create_failed", i18n.Args{"error": err.Error()})
		}

		// Get the groups to add the user to
		groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
		if err != nil {
			common.Log.Errorf("AddUsers: failed to get department information by department ID: %s", err.Error())
			return tools.NewMySqlI18nError("openldap_sync.user_department_info_failed", i18n.Args{"error": err.Error()})
		}
		for _, group := range groups {
			if group.GroupDN[:3] == "ou=" {
				continue
			}
			// First save the user and department information in MySQL
			err := isql.Group.AddUserToGroup(group, []model.User{*user})
			if err != nil {
				common.Log.Errorf("AddUsers: failed to add user-group relation in MySQL: %s", err.Error())
				return tools.NewMySqlI18nError("openldap_sync.user_mysql_group_relation_failed", i18n.Args{"error": err.Error()})
			}
		}
		return nil
	}
	return nil
}
