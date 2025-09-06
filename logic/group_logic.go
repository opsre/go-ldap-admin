package logic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eryajf/go-ldap-admin/config"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/model/response"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
)

type GroupLogic struct{}

// Add 添加数据
func (l GroupLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 获取当前用户
	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户信息失败"))
	}

	group := model.Group{
		GroupType: r.GroupType,
		ParentId:  r.ParentId,
		GroupName: r.GroupName,
		Remark:    r.Remark,
		Creator:   ctxUser.Username,
		Source:    "platform", //默认是平台添加
	}

	if r.ParentId == 0 {
		group.SourceDeptId = "platform_0"
		group.SourceDeptParentId = "platform_0"
		group.GroupDN = fmt.Sprintf("%s=%s,%s", r.GroupType, r.GroupName, config.Conf.Ldap.BaseDN)
	} else {
		parentGroup := new(model.Group)
		err := isql.Group.Find(tools.H{"id": r.ParentId}, parentGroup)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("获取父级组信息失败"))
		}
		group.SourceDeptId = "platform_0"
		group.SourceDeptParentId = fmt.Sprintf("%s_%d", parentGroup.Source, r.ParentId)
		group.GroupDN = fmt.Sprintf("%s=%s,%s", r.GroupType, r.GroupName, parentGroup.GroupDN)
	}

	// 根据 group_dn 判断分组是否已存在
	if isql.Group.Exist(tools.H{"group_dn": group.GroupDN}) {
		return nil, tools.NewValidatorError(fmt.Errorf("该分组对应DN已存在"))
	}

	// 先在ldap中创建组
	err = ildap.Group.Add(&group)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("%s", "向LDAP创建分组失败"+err.Error()))
	}

	// 然后在数据库中创建组
	err = isql.Group.Add(&group)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("向MySQL创建分组失败"))
	}

	// 默认创建分组之后，需要将admin添加到分组中
	adminInfo := new(model.User)
	err = isql.User.Find(tools.H{"id": 1}, adminInfo)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}

	err = isql.Group.AddUserToGroup(&group, []model.User{*adminInfo})
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("添加用户到分组失败: %s", err.Error()))
	}

	return nil, nil
}

// List 数据列表
func (l GroupLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 获取数据列表
	groups, err := isql.Group.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组列表失败: %s", err.Error()))
	}

	rets := make([]model.Group, 0)
	for _, group := range groups {
		rets = append(rets, *group)
	}
	count, err := isql.Group.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组总数失败"))
	}

	return response.GroupListRsp{
		Total:  count,
		Groups: rets,
	}, nil
}

// GetTree 数据树
func (l GroupLogic) GetTree(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	var groups []*model.Group
	groups, err := isql.Group.ListTree(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "获取资源列表失败: "+err.Error()))
	}

	tree := isql.GenGroupTree(0, groups)

	return tree, nil
}

// Update 更新数据
func (l GroupLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": int(r.ID)}
	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("分组不存在"))
	}

	// 获取当前登陆用户
	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户失败"))
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

	//若配置了不允许修改分组名称，则不更新分组名称
	if !config.Conf.Ldap.GroupNameModify {
		newGroup.GroupName = oldGroup.GroupName
	}

	err = ildap.Group.Update(oldGroup, &newGroup)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("%s", "向LDAP更新分组失败："+err.Error()))
	}
	err = isql.Group.Update(&newGroup)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("向MySQL更新分组失败"))
	}
	return nil, nil
}

// Delete 删除数据
func (l GroupLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.GroupIds {
		filter := tools.H{"id": int(id)}
		if !isql.Group.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("有分组不存在"))
		}
	}

	groups, err := isql.Group.GetGroupByIds(r.GroupIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组列表失败: %s", err.Error()))
	}

	for _, group := range groups {
		// 判断存在子分组，不允许删除
		filter := tools.H{"parent_id": int(group.ID)}
		if isql.Group.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("存在子分组，请先删除子分组，再执行该分组的删除操作！"))
		}

		// 删除的时候先从ldap进行删除
		err = ildap.Group.Delete(group.GroupDN)
		if err != nil {
			return nil, tools.NewLdapError(fmt.Errorf("%s", "向LDAP删除分组失败："+err.Error()))
		}
	}

	// 从MySQL中删除
	err = isql.Group.Delete(groups)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除接口失败: %s", err.Error()))
	}

	return nil, nil
}

// AddUser 添加用户到分组
func (l GroupLogic) AddUser(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupAddUserReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("分组不存在"))
	}

	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取用户列表失败: %s", err.Error()))
	}

	group := new(model.Group)
	err = isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组失败: %s", err.Error()))
	}

	if group.GroupDN[:3] == "ou=" {
		return nil, tools.NewMySqlError(fmt.Errorf("ou类型的分组不能添加用户"))
	}

	// 先添加到MySQL
	err = isql.Group.AddUserToGroup(group, users)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("添加用户到分组失败: %s", err.Error()))
	}

	// 再往ldap添加
	for _, user := range users {
		err = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
		if err != nil {
			return nil, tools.NewLdapError(fmt.Errorf("%s", "向LDAP添加用户到分组失败"+err.Error()))
		}
	}

	for _, user := range users {
		oldData := new(model.User)
		err = isql.User.Find(tools.H{"id": user.ID}, oldData)
		if err != nil {
			return nil, tools.NewMySqlError(err)
		}
		newData := oldData
		// 添加新增的分组ID与部门
		newData.DepartmentId = oldData.DepartmentId + "," + strconv.Itoa(int(r.GroupID))
		newData.Departments = oldData.Departments + "," + group.GroupName
		err = l.updataUser(newData)
		if err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("%s", "处理用户的部门数据失败:"+err.Error()))
		}
	}

	return nil, nil
}

func (l GroupLogic) updataUser(newUser *model.User) error {
	err := isql.User.Update(newUser)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("%s", "在MySQL更新用户失败："+err.Error()))
	}
	return nil
}

// RemoveUser 移除用户
func (l GroupLogic) RemoveUser(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupRemoveUserReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("分组不存在"))
	}

	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取用户列表失败: %s", err.Error()))
	}

	group := new(model.Group)
	err = isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组失败: %s", err.Error()))
	}

	if group.GroupDN[:3] == "ou=" {
		return nil, tools.NewMySqlError(fmt.Errorf("ou类型的分组内没有用户"))
	}

	// 先操作ldap
	for _, user := range users {
		err := ildap.Group.RemoveUserFromGroup(group.GroupDN, user.UserDN)
		if err != nil {
			return nil, tools.NewLdapError(fmt.Errorf("%s", "将用户从ldap移除失败"+err.Error()))
		}
	}

	// 再操作MySQL
	err = isql.Group.RemoveUserFromGroup(group, users)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("将用户从MySQL移除失败: %s", err.Error()))
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
		// 删掉移除的分组名字
		for _, v := range strings.Split(oldData.Departments, ",") {
			if v != group.GroupName {
				newDepts = append(newDepts, v)
			}
		}
		// 删掉移除的分组id
		for _, v := range strings.Split(oldData.DepartmentId, ",") {
			if v != strconv.Itoa(int(r.GroupID)) {
				newDeptIds = append(newDeptIds, v)
			}
		}

		newData.Departments = strings.Join(newDepts, ",")
		newData.DepartmentId = strings.Join(newDeptIds, ",")
		err = l.updataUser(newData)
		if err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("%s", "处理用户的部门数据失败:"+err.Error()))
		}
	}

	return nil, nil
}

// UserInGroup 在分组内的用户
func (l GroupLogic) UserInGroup(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserInGroupReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("分组不存在"))
	}

	group := new(model.Group)
	err := isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组失败: %s", err.Error()))
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

// UserNoInGroup 不在分组内的用户
func (l GroupLogic) UserNoInGroup(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserNoInGroupReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !isql.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("分组不存在"))
	}

	group := new(model.Group)
	err := isql.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组失败: %s", err.Error()))
	}

	var userList []*model.User
	userList, err = isql.User.ListAll()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "获取资源列表失败: "+err.Error()))
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
