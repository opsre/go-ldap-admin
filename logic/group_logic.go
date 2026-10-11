package logic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eryajf/go-ldap-admin/config"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/model/response"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
)

type GroupLogic struct{}

// Add data
func (l GroupLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// Get the current user
	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.current_user_failed", nil)
	}

	group := model.Group{
		GroupType: r.GroupType,
		ParentId:  r.ParentId,
		GroupName: r.GroupName,
		Remark:    r.Remark,
		Creator:   ctxUser.Username,
		Source:    "platform", //By default, the platform adds it
	}

	if r.ParentId == 0 {
		group.SourceDeptId = "platform_0"
		group.SourceDeptParentId = "platform_0"
		group.GroupDN = fmt.Sprintf("%s=%s,%s", r.GroupType, r.GroupName, config.Conf.Ldap.BaseDN)
	} else {
		parentGroup := new(model.Group)
		err := isql.Group.Find(tools.H{"id": r.ParentId}, parentGroup)
		if err != nil {
			return nil, tools.NewMySqlI18nError("group.parent_info_failed", nil)
		}
		group.SourceDeptId = "platform_0"
		group.SourceDeptParentId = fmt.Sprintf("%s_%d", parentGroup.Source, r.ParentId)
		group.GroupDN = fmt.Sprintf("%s=%s,%s", r.GroupType, r.GroupName, parentGroup.GroupDN)
	}

	// Check whether the group exists by group_dn
	if isql.Group.Exist(tools.H{"group_dn": group.GroupDN}) {
		return nil, tools.NewValidatorI18nError("group.dn_exists", nil)
	}

	// First create the group in LDAP
	err = ildap.Group.Add(&group)
	if err != nil {
		return nil, tools.NewLdapI18nError("group.ldap_create_failed", i18n.Args{"error": err.Error()})
	}

	// Then create the group in the database
	err = isql.Group.Add(&group)
	if err != nil {
		return nil, tools.NewLdapI18nError("group.mysql_create_failed", nil)
	}

	// By default, add admin to the group after creating it
	adminInfo := new(model.User)
	err = isql.User.Find(tools.H{"id": 1}, adminInfo)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}

	err = isql.Group.AddUserToGroup(&group, []model.User{*adminInfo})
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.add_user_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// List data
func (l GroupLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// Get the data list
	groups, err := isql.Group.List(r)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.list_failed", i18n.Args{"error": err.Error()})
	}

	rets := make([]model.Group, 0)
	for _, group := range groups {
		localizeGroup(c, group)
		rets = append(rets, *group)
	}
	count, err := isql.Group.Count()
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.count_failed", nil)
	}

	return response.GroupListRsp{
		Total:  count,
		Groups: rets,
	}, nil
}

// Get the data tree
func (l GroupLogic) GetTree(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	var groups []*model.Group
	groups, err := isql.Group.ListTree(r)
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.resource_list_failed", i18n.Args{"error": err.Error()})
	}

	tree := isql.GenGroupTree(0, groups)
	localizeGroups(c, tree)

	return tree, nil
}

// Update data
func (l GroupLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": int(r.ID)}
	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlI18nError("group.not_found", nil)
	}

	// Get the currently logged-in user
	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.current_user_failed", nil)
	}

	oldGroup := new(model.Group)
	err = isql.Group.Find(filter, oldGroup)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}

	newGroup := model.Group{
		Model:     oldGroup.Model,
		GroupName: r.GroupName,
		Remark:    r.Remark,
		Creator:   ctxUser.Username,
		GroupType: oldGroup.GroupType,
	}

	//If group name changes are disabled in the configuration, do not update the group name
	if !config.Conf.Ldap.GroupNameModify {
		newGroup.GroupName = oldGroup.GroupName
	}

	err = ildap.Group.Update(oldGroup, &newGroup)
	if err != nil {
		return nil, tools.NewLdapI18nError("group.ldap_update_failed", i18n.Args{"error": err.Error()})
	}
	err = isql.Group.Update(&newGroup)
	if err != nil {
		return nil, tools.NewLdapI18nError("group.mysql_update_failed", nil)
	}
	return nil, nil
}

// Delete data
func (l GroupLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.GroupIds {
		filter := tools.H{"id": int(id)}
		if !isql.Group.Exist(filter) {
			return nil, tools.NewMySqlI18nError("group.some_not_found", nil)
		}
	}

	groups, err := isql.Group.GetGroupByIds(r.GroupIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.list_failed", i18n.Args{"error": err.Error()})
	}

	for _, group := range groups {
		// Check for child groups; deletion is not allowed if any exist
		filter := tools.H{"parent_id": int(group.ID)}
		if isql.Group.Exist(filter) {
			return nil, tools.NewMySqlI18nError("group.has_children", nil)
		}

		// Delete from LDAP first
		err = ildap.Group.Delete(group.GroupDN)
		if err != nil {
			return nil, tools.NewLdapI18nError("group.ldap_delete_failed", i18n.Args{"error": err.Error()})
		}
	}

	// Delete from MySQL
	err = isql.Group.Delete(groups)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.delete_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// Add a user to a group
func (l GroupLogic) AddUser(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupAddUserReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlI18nError("group.not_found", nil)
	}

	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.user_list_failed", i18n.Args{"error": err.Error()})
	}

	group := new(model.Group)
	err = isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.get_failed", i18n.Args{"error": err.Error()})
	}

	if group.GroupDN[:3] == "ou=" {
		return nil, tools.NewMySqlI18nError("group.ou_cannot_add_user", nil)
	}

	// First add it to MySQL
	err = isql.Group.AddUserToGroup(group, users)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.add_user_failed", i18n.Args{"error": err.Error()})
	}

	// Then add it to LDAP
	for _, user := range users {
		err = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
		if err != nil {
			return nil, tools.NewLdapI18nError("group.ldap_add_user_failed", i18n.Args{"error": err.Error()})
		}
	}

	for _, user := range users {
		oldData := new(model.User)
		err = isql.User.Find(tools.H{"id": user.ID}, oldData)
		if err != nil {
			return nil, tools.NewMySqlError(err)
		}
		newData := oldData
		// Add the new group IDs and departments
		newData.DepartmentId = oldData.DepartmentId + "," + strconv.Itoa(int(r.GroupID))
		newData.Departments = oldData.Departments + "," + group.GroupName
		err = l.updataUser(newData)
		if err != nil {
			return nil, tools.NewOperationI18nError("group.user_department_failed", i18n.Args{"error": err.Error()})
		}
	}

	return nil, nil
}

func (l GroupLogic) updataUser(newUser *model.User) error {
	err := isql.User.Update(newUser)
	if err != nil {
		return tools.NewMySqlI18nError("group.mysql_user_update_failed", i18n.Args{"error": err.Error()})
	}
	return nil
}

// Remove a user
func (l GroupLogic) RemoveUser(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupRemoveUserReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlI18nError("group.not_found", nil)
	}

	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.user_list_failed", i18n.Args{"error": err.Error()})
	}

	group := new(model.Group)
	err = isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.get_failed", i18n.Args{"error": err.Error()})
	}

	if group.GroupDN[:3] == "ou=" {
		return nil, tools.NewMySqlI18nError("group.ou_has_no_user", nil)
	}

	// First update LDAP
	for _, user := range users {
		err := ildap.Group.RemoveUserFromGroup(group.GroupDN, user.UserDN)
		if err != nil {
			return nil, tools.NewLdapI18nError("group.ldap_remove_user_failed", i18n.Args{"error": err.Error()})
		}
	}

	// Then update MySQL
	err = isql.Group.RemoveUserFromGroup(group, users)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.mysql_remove_user_failed", i18n.Args{"error": err.Error()})
	}

	for _, user := range users {
		oldData := new(model.User)
		err = isql.User.Find(tools.H{"id": user.ID}, oldData)
		if err != nil {
			return nil, tools.NewMySqlError(err)
		}
		newData := oldData

		var newDepts []string
		var newDeptIds []string
		// Remove the names of removed groups
		for _, v := range strings.Split(oldData.Departments, ",") {
			if v != group.GroupName {
				newDepts = append(newDepts, v)
			}
		}
		// Remove the IDs of removed groups
		for _, v := range strings.Split(oldData.DepartmentId, ",") {
			if v != strconv.Itoa(int(r.GroupID)) {
				newDeptIds = append(newDeptIds, v)
			}
		}

		newData.Departments = strings.Join(newDepts, ",")
		newData.DepartmentId = strings.Join(newDeptIds, ",")
		err = l.updataUser(newData)
		if err != nil {
			return nil, tools.NewOperationI18nError("group.user_department_failed", i18n.Args{"error": err.Error()})
		}
	}

	return nil, nil
}

// UserInGroup lists users in the group
func (l GroupLogic) UserInGroup(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserInGroupReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlI18nError("group.not_found", nil)
	}

	group := new(model.Group)
	err := isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.get_failed", i18n.Args{"error": err.Error()})
	}

	rets := make([]response.Guser, 0)

	for _, user := range group.Users {
		if r.Nickname != "" && !strings.Contains(user.Nickname, r.Nickname) {
			continue
		}
		rets = append(rets, response.Guser{
			UserId:       int64(user.ID),
			UserName:     user.Username,
			NickName:     user.Nickname,
			Mail:         user.Mail,
			JobNumber:    user.JobNumber,
			Mobile:       user.Mobile,
			Introduction: user.Introduction,
		})
	}

	return response.GroupUsers{
		GroupId:     int64(group.ID),
		GroupName:   group.GroupName,
		GroupRemark: group.Remark,
		UserList:    rets,
	}, nil
}

// UserNoInGroup lists users outside the group
func (l GroupLogic) UserNoInGroup(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserNoInGroupReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlI18nError("group.not_found", nil)
	}

	group := new(model.Group)
	err := isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlI18nError("group.get_failed", i18n.Args{"error": err.Error()})
	}

	var userList []*model.User
	userList, err = isql.User.ListAll()
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.resource_list_failed", i18n.Args{"error": err.Error()})
	}

	rets := make([]response.Guser, 0)
	for _, user := range userList {
		in := true
		for _, groupUser := range group.Users {
			if user.Username == groupUser.Username {
				in = false
				break
			}
		}
		if in {
			if r.Nickname != "" && !strings.Contains(user.Nickname, r.Nickname) {
				continue
			}
			rets = append(rets, response.Guser{
				UserId:       int64(user.ID),
				UserName:     user.Username,
				NickName:     user.Nickname,
				Mail:         user.Mail,
				JobNumber:    user.JobNumber,
				Mobile:       user.Mobile,
				Introduction: user.Introduction,
			})
		}
	}

	return response.GroupUsers{
		GroupId:     int64(group.ID),
		GroupName:   group.GroupName,
		GroupRemark: group.Remark,
		UserList:    rets,
	}, nil
}
