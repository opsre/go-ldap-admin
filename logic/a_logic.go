package logic

import (
	"fmt"
	"math/rand"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
)

var (
	ReqAssertErr = tools.NewRspError(tools.SystemErr, fmt.Errorf("请求异常"))

	Api           = &ApiLogic{}
	User          = &UserLogic{}
	Group         = &GroupLogic{}
	Role          = &RoleLogic{}
	Menu          = &MenuLogic{}
	OperationLog  = &OperationLogLogic{}
	DingTalk      = &DingTalkLogic{}
	WeCom         = &WeComLogic{}
	FeiShu        = &FeiShuLogic{}
	OpenLdap      = &OpenLdapLogic{}
	Sql           = &SqlLogic{}
	Base          = &BaseLogic{}
	FieldRelation = &FieldRelationLogic{}

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// CommonAddGroup 标准创建分组
func CommonAddGroup(group *model.Group) error {
	// 先在ldap中创建组
	err := ildap.Group.Add(group)
	if err != nil {
		return err
	}

	// 然后在数据库中创建组
	err = isql.Group.Add(group)
	if err != nil {
		return err
	}

	// 默认创建分组之后，需要将admin添加到分组中
	adminInfo := new(model.User)
	err = isql.User.Find(tools.H{"id": 1}, adminInfo)
	if err != nil {
		return err
	}

	err = isql.Group.AddUserToGroup(group, []model.User{*adminInfo})
	if err != nil {
		return err
	}

	return nil
}

// CommonUpdateGroup 标准更新分组
func CommonUpdateGroup(oldGroup, newGroup *model.Group) error {
	//若配置了不允许修改分组名称，则不更新分组名称
	if !config.Conf.Ldap.GroupNameModify {
		newGroup.GroupName = oldGroup.GroupName
	}

	err := ildap.Group.Update(oldGroup, newGroup)
	if err != nil {
		return err
	}
	err = isql.Group.Update(newGroup)
	if err != nil {
		return err
	}
	return nil
}

// CommonAddUser 标准创建用户
func CommonAddUser(user *model.User, groups []*model.Group) error {
	// 用户信息的预置处理
	if user.Nickname == "" {
		user.Nickname = "佚名"
	}
	if user.GivenName == "" {
		user.GivenName = user.Nickname
	}
	if user.Introduction == "" {
		user.Introduction = user.Nickname
	}
	if user.Mail == "" {
		// 兼容
		if len(config.Conf.Ldap.DefaultEmailSuffix) > 0 {
			user.Mail = user.Username + "@" + config.Conf.Ldap.DefaultEmailSuffix
		} else {
			user.Mail = user.Username + "@eryajf.net"
		}
	}
	if user.JobNumber == "" {
		user.JobNumber = "0000"
	}
	if user.Departments == "" {
		user.Departments = "默认:研发中心"
	}
	if user.Position == "" {
		user.Position = "默认:打工人"
	}
	if user.PostalAddress == "" {
		user.PostalAddress = "默认:地球"
	}
	if user.Mobile == "" {
		user.Mobile = generateMobile()
	}

	// 先将用户添加到MySQL
	err := isql.User.Add(user)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("%s", "向MySQL创建用户失败："+err.Error()))
	}

	// 发送用户创建成功通知邮件
	if err := tools.SendUserCreationNotification(user.Username, user.Nickname, user.Mail, tools.NewParPasswd(user.Password)); err != nil {
		common.Log.Warnf("发送用户创建通知邮件失败，用户: %s, 邮箱: %s, 错误: %v", user.Username, user.Mail, err)
	}
	// 再将用户添加到ldap
	err = ildap.User.Add(user)
	if err != nil {
		return tools.NewLdapError(fmt.Errorf("%s", "AddUser向LDAP创建用户失败："+err.Error()))
	}

	// 处理用户归属的组
	for _, group := range groups {
		if group.GroupDN[:3] == "ou=" {
			continue
		}
		// 先将用户和部门信息维护到MySQL
		err := isql.Group.AddUserToGroup(group, []model.User{*user})
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("%s", "向MySQL添加用户到分组关系失败："+err.Error()))
		}
		//根据选择的部门，添加到部门内
		err = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("%s", "向Ldap添加用户到分组关系失败："+err.Error()))
		}
	}
	return nil
}

// CommonUpdateUser 标准更新用户
func CommonUpdateUser(oldUser, newUser *model.User, groupId []uint) error {
	// 更新用户
	if !config.Conf.Ldap.UserNameModify {
		newUser.Username = oldUser.Username
	}

	err := ildap.User.Update(oldUser.Username, newUser)
	if err != nil {
		return tools.NewLdapError(fmt.Errorf("%s", "在LDAP更新用户失败："+err.Error()))
	}

	err = isql.User.Update(newUser)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("%s", "在MySQL更新用户失败："+err.Error()))
	}

	//判断部门信息是否有变化有变化则更新相应的数据库
	oldDeptIds := tools.StringToSlice(oldUser.DepartmentId, ",")
	addDeptIds, removeDeptIds := tools.ArrUintCmp(oldDeptIds, groupId)

	// 先处理添加的部门
	addgroups, err := isql.Group.GetGroupByIds(addDeptIds)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("%s", "根据部门ID获取部门信息失败"+err.Error()))
	}
	for _, group := range addgroups {
		if group.GroupDN[:3] == "ou=" {
			continue
		}
		// 先将用户和部门信息维护到MySQL
		err := isql.Group.AddUserToGroup(group, []model.User{*newUser})
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("%s", "向MySQL添加用户到分组关系失败："+err.Error()))
		}
		//根据选择的部门，添加到部门内
		err = ildap.Group.AddUserToGroup(group.GroupDN, newUser.UserDN)
		if err != nil {
			return tools.NewLdapError(fmt.Errorf("%s", "向Ldap添加用户到分组关系失败："+err.Error()))
		}
	}

	// 再处理删除的部门
	removegroups, err := isql.Group.GetGroupByIds(removeDeptIds)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("%s", "根据部门ID获取部门信息失败"+err.Error()))
	}
	for _, group := range removegroups {
		if group.GroupDN[:3] == "ou=" {
			continue
		}
		err := isql.Group.RemoveUserFromGroup(group, []model.User{*newUser})
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("%s", "在MySQL将用户从分组移除失败："+err.Error()))
		}
		err = ildap.Group.RemoveUserFromGroup(group.GroupDN, newUser.UserDN)
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("%s", "在ldap将用户从分组移除失败："+err.Error()))
		}
	}
	return nil
}

// BuildGroupData 根据数据与动态字段组装成分组数据
func BuildGroupData(flag string, remoteData map[string]any) (*model.Group, error) {
	output, err := json.Marshal(&remoteData)
	if err != nil {
		return nil, err
	}

	oldData := new(model.FieldRelation)
	err = isql.FieldRelation.Find(tools.H{"flag": flag + "_group"}, oldData)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}
	frs, err := tools.JsonToMap(string(oldData.Attributes))
	if err != nil {
		return nil, tools.NewOperationError(err)
	}

	g := &model.Group{}
	for system, remote := range frs {
		switch system {
		case "groupName":
			g.SetGroupName(gjson.Get(string(output), remote).String())
		case "remark":
			g.SetRemark(gjson.Get(string(output), remote).String())
		case "sourceDeptId":
			g.SetSourceDeptId(fmt.Sprintf("%s_%s", flag, gjson.Get(string(output), remote).String()))
		case "sourceDeptParentId":
			g.SetSourceDeptParentId(fmt.Sprintf("%s_%s", flag, gjson.Get(string(output), remote).String()))
		}
	}
	return g, nil
}

// BuildUserData 根据数据与动态字段组装成用户数据
func BuildUserData(flag string, remoteData map[string]any) (*model.User, error) {
	output, err := json.Marshal(&remoteData)
	if err != nil {
		return nil, err
	}

	fieldRelationSource := new(model.FieldRelation)
	err = isql.FieldRelation.Find(tools.H{"flag": flag + "_user"}, fieldRelationSource)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}
	fieldRelation, err := tools.JsonToMap(string(fieldRelationSource.Attributes))
	if err != nil {
		return nil, tools.NewOperationError(err)
	}

	// 校验username是否为空，username为必填项
	name := gjson.Get(string(output), fieldRelation["username"]).String()
	if len(name) == 0 {
		common.Log.Warnf("%s 该用户未填写username", output)
		return nil, nil
	}

	u := &model.User{}
	for system, remote := range fieldRelation {
		switch system {
		case "username":
			u.SetUserName(gjson.Get(string(output), remote).String())
		case "nickname":
			u.SetNickName(gjson.Get(string(output), remote).String())
		case "givenName":
			u.SetGivenName(gjson.Get(string(output), remote).String())
		case "mail":
			u.SetMail(gjson.Get(string(output), remote).String())
		case "jobNumber":
			u.SetJobNumber(gjson.Get(string(output), remote).String())
		case "mobile":
			u.SetMobile(gjson.Get(string(output), remote).String())
		case "avatar":
			u.SetAvatar(gjson.Get(string(output), remote).String())
		case "postalAddress":
			u.SetPostalAddress(gjson.Get(string(output), remote).String())
		case "position":
			u.SetPosition(gjson.Get(string(output), remote).String())
		case "introduction":
			u.SetIntroduction(gjson.Get(string(output), remote).String())
		case "sourceUserId":
			u.SetSourceUserId(fmt.Sprintf("%s_%s", flag, gjson.Get(string(output), remote).String()))
		case "sourceUnionId":
			u.SetSourceUnionId(fmt.Sprintf("%s_%s", flag, gjson.Get(string(output), remote).String()))
		}
	}
	return u, nil
}

// ConvertDeptData 将部门信息转成本地结构体
func ConvertDeptData(flag string, remoteData []map[string]any) (groups []*model.Group, err error) {
	for _, dept := range remoteData {
		group, err := BuildGroupData(flag, dept)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return
}

// ConvertUserData 将用户信息转成本地结构体
func ConvertUserData(flag string, remoteData []map[string]any) (users []*model.User, err error) {
	for _, staff := range remoteData {
		groupIds, err := isql.Group.DeptIdsToGroupIds(staff["department_ids"].([]string))
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("将部门ids转换为内部部门id失败：%s", err.Error()))
		}
		user, err := BuildUserData(flag, staff)
		if err != nil {
			return nil, err
		}
		if user != nil {
			user.DepartmentId = tools.SliceToString(groupIds, ",")
			users = append(users, user)
		}
	}
	return
}

func InitCron() {
	c := cron.New(cron.WithSeconds())

	if config.Conf.DingTalk.EnableSync {
		//启动定时任务
		_, err := c.AddFunc(config.Conf.DingTalk.DeptSyncTime, func() {
			DingTalk.SyncDingTalkDepts(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("启动同步部门的定时任务失败: %v", err)
		}
		//每天凌晨1点执行一次
		_, err = c.AddFunc(config.Conf.DingTalk.UserSyncTime, func() {
			DingTalk.SyncDingTalkUsers(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("启动同步用户的定时任务失败: %v", err)
		}
	}
	if config.Conf.WeCom.EnableSync {
		_, err := c.AddFunc(config.Conf.WeCom.DeptSyncTime, func() {
			WeCom.SyncWeComDepts(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("启动同步部门的定时任务失败: %v", err)
		}
		//每天凌晨1点执行一次
		_, err = c.AddFunc(config.Conf.WeCom.UserSyncTime, func() {
			WeCom.SyncWeComUsers(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("启动同步用户的定时任务失败: %v", err)
		}
	}
	if config.Conf.FeiShu.EnableSync {
		_, err := c.AddFunc(config.Conf.FeiShu.DeptSyncTime, func() {
			FeiShu.SyncFeiShuDepts(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("启动同步部门的定时任务失败: %v", err)
		}
		//每天凌晨1点执行一次
		_, err = c.AddFunc(config.Conf.FeiShu.UserSyncTime, func() {
			FeiShu.SyncFeiShuUsers(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("启动同步用户的定时任务失败: %v", err)
		}
	}

	// 自动检索未同步数据
	_, err := c.AddFunc("0 */2 * * * *", func() {
		// 开发调试时调整为10秒执行一次
		// _, err := c.AddFunc("*/10 * * * * *", func() {
		_ = SearchGroupDiff()
		_ = SearchUserDiff()
	})
	if err != nil {
		common.Log.Errorf("启动同步任务状态检查任务失败: %v", err)
	}
	c.Start()
}

func GroupListToTree(rootId string, groupList []*model.Group) *model.Group {
	// 创建空根节点
	rootGroup := &model.Group{SourceDeptId: rootId}
	rootGroup.Children = groupListToTree(rootGroup, groupList)
	return rootGroup
}

func groupListToTree(rootGroup *model.Group, list []*model.Group) []*model.Group {
	children := make([]*model.Group, 0)
	for _, group := range list {
		if group.SourceDeptParentId == rootGroup.SourceDeptId {
			children = append(children, group)
		}
	}
	for _, group := range children {
		group.Children = groupListToTree(group, list)
	}
	return children
}

func generateMobile() string {
	randNum := rand.Intn(9000000000) + 1000000000
	randNum = randNum + 10000000000
	if isql.User.Exist(tools.H{"mobile": randNum}) {
		return generateMobile()
	}
	return fmt.Sprintf("%v", randNum)
}
