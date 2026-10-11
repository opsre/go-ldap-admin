package logic

import (
	"fmt"
	"math/rand"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
)

var (
	ReqAssertErr = tools.NewRspI18nError(tools.SystemErr, "common.request_invalid", nil)

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

// CommonAddGroup creates a group using the standard flow
func CommonAddGroup(group *model.Group) error {
	// First create the group in LDAP
	err := ildap.Group.Add(group)
	if err != nil {
		return err
	}

	// Then create the group in the database
	err = isql.Group.Add(group)
	if err != nil {
		return err
	}

	// By default, add admin to the group after creating it
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

// CommonUpdateGroup updates a group using the standard flow
func CommonUpdateGroup(oldGroup, newGroup *model.Group) error {
	//If group name changes are disabled in the configuration, do not update the group name
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

// CommonAddUser creates a user using the standard flow
func CommonAddUser(user *model.User, groups []*model.Group, locale ...string) error {
	// Preprocess user information
	if user.Nickname == "" {
		user.Nickname = "Anonymous"
	}
	if user.GivenName == "" {
		user.GivenName = user.Nickname
	}
	if user.Introduction == "" {
		user.Introduction = user.Nickname
	}
	if user.Mail == "" {
		// Compatibility handling
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
		user.Departments = "Default: R&D Center"
	}
	if user.Position == "" {
		user.Position = "Default: Staff"
	}
	if user.PostalAddress == "" {
		user.PostalAddress = "Default: Earth"
	}
	if user.Mobile == "" {
		user.Mobile = generateMobile()
	}

	// First add the user to MySQL
	err := isql.User.Add(user)
	if err != nil {
		return tools.NewMySqlI18nError("common_user.mysql_create_failed", i18n.Args{"error": err.Error()})
	}

	// Send the user creation success notification email
	if err := tools.SendUserCreationNotificationI18n(user.Username, user.Nickname, user.Mail, tools.NewParPasswd(user.Password), firstLocale(locale)); err != nil {
		common.Log.Warnf("Failed to send user creation notification email, user: %s, email: %s, error: %v", user.Username, user.Mail, err)
	}
	// Then add the user to LDAP
	err = ildap.User.Add(user)
	if err != nil {
		return tools.NewLdapI18nError("common_user.ldap_create_failed", i18n.Args{"error": err.Error()})
	}

	// Handle the groups the user belongs to
	for _, group := range groups {
		if group.GroupDN[:3] == "ou=" {
			continue
		}
		// First save the user and department information in MySQL
		err := isql.Group.AddUserToGroup(group, []model.User{*user})
		if err != nil {
			return tools.NewMySqlI18nError("common_user.mysql_add_group_relation_failed", i18n.Args{"error": err.Error()})
		}
		//Add the user to the selected departments
		err = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
		if err != nil {
			return tools.NewLdapI18nError("common_user.ldap_add_group_relation_failed", i18n.Args{"error": err.Error()})
		}
	}
	return nil
}

func firstLocale(locales []string) string {
	if len(locales) == 0 {
		return ""
	}
	return locales[0]
}

// CommonUpdateUser updates a user using the standard flow
func CommonUpdateUser(oldUser, newUser *model.User, groupId []uint) error {
	// Update the user
	if !config.Conf.Ldap.UserNameModify {
		newUser.Username = oldUser.Username
	}

	err := ildap.User.Update(oldUser.Username, newUser)
	if err != nil {
		return tools.NewLdapI18nError("common_user.ldap_update_failed", i18n.Args{"error": err.Error()})
	}

	err = isql.User.Update(newUser)
	if err != nil {
		return tools.NewMySqlI18nError("common_user.mysql_update_failed", i18n.Args{"error": err.Error()})
	}

	//Check whether department information changed; if so, update the corresponding database
	oldDeptIds := tools.StringToSlice(oldUser.DepartmentId, ",")
	addDeptIds, removeDeptIds := tools.ArrUintCmp(oldDeptIds, groupId)

	// First handle added departments
	addgroups, err := isql.Group.GetGroupByIds(addDeptIds)
	if err != nil {
		return tools.NewMySqlI18nError("user.department_info_failed", i18n.Args{"error": err.Error()})
	}
	for _, group := range addgroups {
		if group.GroupDN[:3] == "ou=" {
			continue
		}
		// First save the user and department information in MySQL
		err := isql.Group.AddUserToGroup(group, []model.User{*newUser})
		if err != nil {
			return tools.NewMySqlI18nError("common_user.mysql_add_group_relation_failed", i18n.Args{"error": err.Error()})
		}
		//Add the user to the selected departments
		err = ildap.Group.AddUserToGroup(group.GroupDN, newUser.UserDN)
		if err != nil {
			return tools.NewLdapI18nError("common_user.ldap_add_group_relation_failed", i18n.Args{"error": err.Error()})
		}
	}

	// Then handle removed departments
	removegroups, err := isql.Group.GetGroupByIds(removeDeptIds)
	if err != nil {
		return tools.NewMySqlI18nError("user.department_info_failed", i18n.Args{"error": err.Error()})
	}
	for _, group := range removegroups {
		if group.GroupDN[:3] == "ou=" {
			continue
		}
		err := isql.Group.RemoveUserFromGroup(group, []model.User{*newUser})
		if err != nil {
			return tools.NewMySqlI18nError("common_user.mysql_remove_group_relation_failed", i18n.Args{"error": err.Error()})
		}
		err = ildap.Group.RemoveUserFromGroup(group.GroupDN, newUser.UserDN)
		if err != nil {
			return tools.NewLdapI18nError("common_user.ldap_remove_group_relation_failed", i18n.Args{"error": err.Error()})
		}
	}
	return nil
}

// BuildGroupData assembles group data from data and dynamic fields
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

// BuildUserData assembles user data from data and dynamic fields
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

	// Check that username is not empty; it is required
	name := gjson.Get(string(output), fieldRelation["username"]).String()
	if len(name) == 0 {
		common.Log.Warnf("%s user has no username", output)
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

// ConvertDeptData converts department information to a local struct
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

// ConvertUserData converts user information to a local struct
func ConvertUserData(flag string, remoteData []map[string]any) (users []*model.User, err error) {
	for _, staff := range remoteData {
		groupIds, err := isql.Group.DeptIdsToGroupIds(staff["department_ids"].([]string))
		if err != nil {
			return nil, tools.NewMySqlI18nError("sync.dept_ids_convert_failed", i18n.Args{"error": err.Error()})
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
		//Start scheduled tasks
		_, err := c.AddFunc(config.Conf.DingTalk.DeptSyncTime, func() {
			DingTalk.SyncDingTalkDepts(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("Failed to start department sync schedule: %v", err)
		}
		//Run once daily at 1 a.m.
		_, err = c.AddFunc(config.Conf.DingTalk.UserSyncTime, func() {
			DingTalk.SyncDingTalkUsers(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("Failed to start user sync schedule: %v", err)
		}
	}
	if config.Conf.WeCom.EnableSync {
		_, err := c.AddFunc(config.Conf.WeCom.DeptSyncTime, func() {
			WeCom.SyncWeComDepts(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("Failed to start department sync schedule: %v", err)
		}
		//Run once daily at 1 a.m.
		_, err = c.AddFunc(config.Conf.WeCom.UserSyncTime, func() {
			WeCom.SyncWeComUsers(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("Failed to start user sync schedule: %v", err)
		}
	}
	if config.Conf.FeiShu.EnableSync {
		_, err := c.AddFunc(config.Conf.FeiShu.DeptSyncTime, func() {
			FeiShu.SyncFeiShuDepts(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("Failed to start department sync schedule: %v", err)
		}
		//Run once daily at 1 a.m.
		_, err = c.AddFunc(config.Conf.FeiShu.UserSyncTime, func() {
			FeiShu.SyncFeiShuUsers(nil, nil)
		})
		if err != nil {
			common.Log.Errorf("Failed to start user sync schedule: %v", err)
		}
	}

	// Automatically search for unsynchronized data
	_, err := c.AddFunc("0 */2 * * * *", func() {
		// During development, change the interval to 10 seconds
		// _, err := c.AddFunc("*/10 * * * * *", func() {
		_ = SearchGroupDiff()
		_ = SearchUserDiff()
	})
	if err != nil {
		common.Log.Errorf("Failed to start sync status check schedule: %v", err)
	}
	c.Start()
}

func GroupListToTree(rootId string, groupList []*model.Group) *model.Group {
	// Create an empty root node
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
