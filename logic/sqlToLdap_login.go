package logic

import (
	"fmt"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"
	"github.com/gin-gonic/gin"
)

type SqlLogic struct{}

// SyncSqlUsers 同步sql的用户信息到ldap
func (d *SqlLogic) SyncSqlUsers(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SyncSqlUserReq)
	if !ok {
		errMsg := "请求参数类型断言失败"
		common.Log.Errorf("SyncSqlUsers: %s", errMsg)
		return nil, ReqAssertErr
	}
	_ = c
	// 1.获取所有用户
	for _, id := range r.UserIds {
		filter := tools.H{"id": int(id)}
		if !isql.User.Exist(filter) {
			errMsg := fmt.Sprintf("用户ID[%d]不存在", id)
			common.Log.Errorf("SyncSqlUsers: %s", errMsg)
			return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
		}
	}
	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		errMsg := fmt.Sprintf("获取用户信息失败: %s", err.Error())
		common.Log.Errorf("SyncSqlUsers: %s", errMsg)
		return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
	}
	// 2.再将用户添加到ldap
	for i, user := range users {
		err = ildap.User.Add(&user)
		if err != nil {
			errMsg := fmt.Sprintf("向LDAP同步用户[%s]失败：%s", user.Username, err.Error())
			common.Log.Errorf("SyncSqlUsers: %s", errMsg)
			return nil, tools.NewLdapError(fmt.Errorf(errMsg))
		}
		// 获取用户将要添加的分组
		groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
		if err != nil {
			errMsg := fmt.Sprintf("用户[%s]根据部门ID获取部门信息失败: %s", user.Username, err.Error())
			common.Log.Errorf("SyncSqlUsers: %s", errMsg)
			return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
		}
		for _, group := range groups {
			//根据选择的部门，添加到部门内
			err = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
			if err != nil {
				errMsg := fmt.Sprintf("向Ldap添加用户[%s]到分组[%s]失败：%s", user.Username, group.GroupName, err.Error())
				common.Log.Errorf("SyncSqlUsers: %s", errMsg)
				return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
			}
		}
		err = isql.User.ChangeSyncState(int(user.ID), 1)
		if err != nil {
			errMsg := fmt.Sprintf("用户[%s]同步完毕之后更新状态失败：%s", user.Username, err.Error())
			common.Log.Errorf("SyncSqlUsers: %s", errMsg)
			return nil, tools.NewLdapError(fmt.Errorf(errMsg))
		}
		common.Log.Infof("SyncSqlUsers: 成功同步用户[%s] (%d/%d)", user.Username, i+1, len(users))
	}

	common.Log.Infof("SyncSqlUsers: SQL用户同步成功，共同步%d个用户", len(users))
	return nil, nil
}

// SyncSqlGroups 同步sql中的分组信息到ldap
func (d *SqlLogic) SyncSqlGroups(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SyncSqlGrooupsReq)
	if !ok {
		errMsg := "请求参数类型断言失败"
		common.Log.Errorf("SyncSqlGroups: %s", errMsg)
		return nil, ReqAssertErr
	}
	_ = c
	// 1.获取所有分组
	for _, id := range r.GroupIds {
		filter := tools.H{"id": int(id)}
		if !isql.Group.Exist(filter) {
			errMsg := fmt.Sprintf("分组ID[%d]不存在", id)
			common.Log.Errorf("SyncSqlGroups: %s", errMsg)
			return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
		}
	}
	groups, err := isql.Group.GetGroupByIds(r.GroupIds)
	if err != nil {
		errMsg := fmt.Sprintf("获取分组信息失败: %s", err.Error())
		common.Log.Errorf("SyncSqlGroups: %s", errMsg)
		return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
	}
	// 2.再将分组添加到ldap
	for _, group := range groups {
		err = ildap.Group.Add(group)
		if err != nil {
			errMsg := fmt.Sprintf("向LDAP同步分组[%s]失败：%s", group.GroupName, err.Error())
			common.Log.Errorf("SyncSqlGroups: %s", errMsg)
			return nil, tools.NewLdapError(fmt.Errorf(errMsg))
		}
		if len(group.Users) > 0 {
			for _, user := range group.Users {
				if user.UserDN == config.Conf.Ldap.AdminDN {
					continue
				}
				err = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
				if err != nil {
					errMsg := fmt.Sprintf("同步分组[%s]之后处理分组内的用户[%s]失败：%s", group.GroupName, user.Username, err.Error())
					common.Log.Errorf("SyncSqlGroups: %s", errMsg)
					return nil, tools.NewLdapError(fmt.Errorf(errMsg))
				}
			}
		}
		err = isql.Group.ChangeSyncState(int(group.ID), 1)
		if err != nil {
			errMsg := fmt.Sprintf("分组[%s]同步完毕之后更新状态失败：%s", group.GroupName, err.Error())
			common.Log.Errorf("SyncSqlGroups: %s", errMsg)
			return nil, tools.NewLdapError(fmt.Errorf(errMsg))
		}
	}

	common.Log.Infof("SyncSqlGroups: SQL分组同步成功，共同步%d个分组", len(groups))
	return nil, nil
}

// SearchGroupDiff 检索未同步到ldap中的分组
func SearchGroupDiff() (err error) {
	// 获取sql中的数据
	var sqlGroupList []*model.Group
	sqlGroupList, err = isql.Group.ListAll()
	if err != nil {
		return err
	}
	// 获取ldap中的数据
	var ldapGroupList []*model.Group
	ldapGroupList, err = ildap.Group.ListGroupDN()
	if err != nil {
		return err
	}
	// 比对两个系统中的数据
	groups := diffGroup(sqlGroupList, ldapGroupList)
	for _, group := range groups {
		if group.GroupDN == config.Conf.Ldap.BaseDN {
			continue
		}
		err = isql.Group.ChangeSyncState(int(group.ID), 2)
	}
	return
}

// SearchUserDiff 检索未同步到ldap中的用户
func SearchUserDiff() (err error) {
	// 获取sql中的数据
	var sqlUserList []*model.User
	sqlUserList, err = isql.User.ListAll()
	if err != nil {
		return err
	}
	// 获取ldap中的数据
	var ldapUserList []*model.User
	ldapUserList, err = ildap.User.ListUserDN()
	if err != nil {
		return err
	}
	// 比对两个系统中的数据
	users := diffUser(sqlUserList, ldapUserList)
	for _, user := range users {
		if user.UserDN == config.Conf.Ldap.AdminDN {
			continue
		}
		err = isql.User.ChangeSyncState(int(user.ID), 2)
	}
	return
}

// diffGroup 比较出sql中有但ldap中没有的group列表
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

// diffUser 比较出sql中有但ldap中没有的user列表
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
