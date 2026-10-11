package logic

import (
	"fmt"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"
	"github.com/gin-gonic/gin"
)

type SqlLogic struct{}

// SyncSqlUsers syncs SQL user information to LDAP
func (d *SqlLogic) SyncSqlUsers(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.SyncSqlUserReq)
	if !ok {
		errMsg := "Request parameter type assertion failed"
		common.Log.Errorf("SyncSqlUsers: %s", errMsg)
		return nil, ReqAssertErr
	}
	_ = c
	// 1. Get all users
	for _, id := range r.UserIds {
		filter := tools.H{"id": int(id)}
		if !isql.User.Exist(filter) {
			errMsg := fmt.Sprintf("User ID [%d] does not exist", id)
			common.Log.Errorf("SyncSqlUsers: %s", errMsg)
			return nil, tools.NewMySqlI18nError("sqlsync.user_id_not_found", i18n.Args{"id": id})
		}
	}
	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get user information: %s", err.Error())
		common.Log.Errorf("SyncSqlUsers: %s", errMsg)
		return nil, tools.NewMySqlI18nError("sqlsync.user_info_failed", i18n.Args{"error": err.Error()})
	}
	// 2. Then add the users to LDAP
	for i, user := range users {
		err = ildap.User.Add(&user)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to sync user [%s] to LDAP: %s", user.Username, err.Error())
			common.Log.Errorf("SyncSqlUsers: %s", errMsg)
			return nil, tools.NewLdapI18nError("sqlsync.user_ldap_sync_failed", i18n.Args{"username": user.Username, "error": err.Error()})
		}
		// Get the groups to add the user to
		groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get department information by department ID for user [%s]: %s", user.Username, err.Error())
			common.Log.Errorf("SyncSqlUsers: %s", errMsg)
			return nil, tools.NewMySqlI18nError("sqlsync.user_department_info_failed", i18n.Args{"username": user.Username, "error": err.Error()})
		}
		for _, group := range groups {
			//Add the user to the selected departments
			err = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to add user [%s] to LDAP group [%s]: %s", user.Username, group.GroupName, err.Error())
				common.Log.Errorf("SyncSqlUsers: %s", errMsg)
				return nil, tools.NewMySqlI18nError("sqlsync.user_add_group_failed", i18n.Args{"username": user.Username, "group": group.GroupName, "error": err.Error()})
			}
		}
		err = isql.User.ChangeSyncState(int(user.ID), 1)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to update status after syncing user [%s]: %s", user.Username, err.Error())
			common.Log.Errorf("SyncSqlUsers: %s", errMsg)
			return nil, tools.NewLdapI18nError("sqlsync.user_status_update_failed", i18n.Args{"username": user.Username, "error": err.Error()})
		}
		common.Log.Infof("SyncSqlUsers: synced user [%s] (%d/%d)", user.Username, i+1, len(users))
	}

	common.Log.Infof("SyncSqlUsers: SQL user sync succeeded, synced %d users", len(users))
	return nil, nil
}

// SyncSqlGroups syncs SQL group information to LDAP
func (d *SqlLogic) SyncSqlGroups(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.SyncSqlGrooupsReq)
	if !ok {
		errMsg := "Request parameter type assertion failed"
		common.Log.Errorf("SyncSqlGroups: %s", errMsg)
		return nil, ReqAssertErr
	}
	_ = c
	// 1. Get all groups
	for _, id := range r.GroupIds {
		filter := tools.H{"id": int(id)}
		if !isql.Group.Exist(filter) {
			errMsg := fmt.Sprintf("Group ID [%d] does not exist", id)
			common.Log.Errorf("SyncSqlGroups: %s", errMsg)
			return nil, tools.NewMySqlI18nError("sqlsync.group_id_not_found", i18n.Args{"id": id})
		}
	}
	groups, err := isql.Group.GetGroupByIds(r.GroupIds)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get group information: %s", err.Error())
		common.Log.Errorf("SyncSqlGroups: %s", errMsg)
		return nil, tools.NewMySqlI18nError("sqlsync.group_info_failed", i18n.Args{"error": err.Error()})
	}
	// 2. Then add the groups to LDAP
	for _, group := range groups {
		err = ildap.Group.Add(group)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to sync group [%s] to LDAP: %s", group.GroupName, err.Error())
			common.Log.Errorf("SyncSqlGroups: %s", errMsg)
			return nil, tools.NewLdapI18nError("sqlsync.group_ldap_sync_failed", i18n.Args{"group": group.GroupName, "error": err.Error()})
		}
		if len(group.Users) > 0 {
			for _, user := range group.Users {
				if user.UserDN == config.Conf.Ldap.AdminDN {
					continue
				}
				err = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
				if err != nil {
					errMsg := fmt.Sprintf("Failed to process user [%s] in group [%s] after syncing the group: %s", user.Username, group.GroupName, err.Error())
					common.Log.Errorf("SyncSqlGroups: %s", errMsg)
					return nil, tools.NewLdapI18nError("sqlsync.group_user_sync_failed", i18n.Args{"group": group.GroupName, "username": user.Username, "error": err.Error()})
				}
			}
		}
		err = isql.Group.ChangeSyncState(int(group.ID), 1)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to update status after syncing group [%s]: %s", group.GroupName, err.Error())
			common.Log.Errorf("SyncSqlGroups: %s", errMsg)
			return nil, tools.NewLdapI18nError("sqlsync.group_status_update_failed", i18n.Args{"group": group.GroupName, "error": err.Error()})
		}
	}

	common.Log.Infof("SyncSqlGroups: SQL group sync succeeded, synced %d groups", len(groups))
	return nil, nil
}

// SearchGroupDiff finds groups not synced to LDAP
func SearchGroupDiff() (err error) {
	// Get the data from SQL
	var sqlGroupList []*model.Group
	sqlGroupList, err = isql.Group.ListAll()
	if err != nil {
		return err
	}
	// Get the data from LDAP
	var ldapGroupList []*model.Group
	ldapGroupList, err = ildap.Group.ListGroupDN()
	if err != nil {
		return err
	}
	// Compare the data in the two systems
	groups := diffGroup(sqlGroupList, ldapGroupList)
	for _, group := range groups {
		if group.GroupDN == config.Conf.Ldap.BaseDN {
			continue
		}
		err = isql.Group.ChangeSyncState(int(group.ID), 2)
	}
	return
}

// SearchUserDiff finds users not synced to LDAP
func SearchUserDiff() (err error) {
	// Get the data from SQL
	var sqlUserList []*model.User
	sqlUserList, err = isql.User.ListAll()
	if err != nil {
		return err
	}
	// Get the data from LDAP
	var ldapUserList []*model.User
	ldapUserList, err = ildap.User.ListUserDN()
	if err != nil {
		return err
	}
	// Compare the data in the two systems
	users := diffUser(sqlUserList, ldapUserList)
	for _, user := range users {
		if user.UserDN == config.Conf.Ldap.AdminDN {
			continue
		}
		err = isql.User.ChangeSyncState(int(user.ID), 2)
	}
	return
}

// diffGroup compares groups present in SQL but absent from LDAP
func diffGroup(sqlGroup, ldapGroup []*model.Group) (rst []*model.Group) {
	var tmp = make(map[string]struct{}, 0)

	for _, v := range ldapGroup {
		tmp[v.GroupDN] = struct{}{}
	}

	for _, v := range sqlGroup {
		if _, ok := tmp[v.GroupDN]; !ok {
			rst = append(rst, v)
		}
	}
	return
}

// diffUser compares users present in SQL but absent from LDAP
func diffUser(sqlUser, ldapUser []*model.User) (rst []*model.User) {
	var tmp = make(map[string]struct{}, len(sqlUser))

	for _, v := range ldapUser {
		tmp[v.UserDN] = struct{}{}
	}

	for _, v := range sqlUser {
		if _, ok := tmp[v.UserDN]; !ok {
			rst = append(rst, v)
		}
	}
	return
}
