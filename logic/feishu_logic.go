package logic

import (
	"fmt"
	"strings"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/client/feishu"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"
	"github.com/gin-gonic/gin"
)

type FeiShuLogic struct {
}

// 通过飞书获取部门信息
func (d *FeiShuLogic) SyncFeiShuDepts(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	// 1.获取所有部门
	deptSource, err := feishu.GetAllDepts()
	if err != nil {
		errMsg := fmt.Sprintf("获取飞书部门列表失败：%s", err.Error())
		common.Log.Errorf("SyncFeiShuDepts: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
	}
	depts, err := ConvertDeptData(config.Conf.FeiShu.Flag, deptSource)
	if err != nil {
		errMsg := fmt.Sprintf("转换飞书部门数据失败：%s", err.Error())
		common.Log.Errorf("SyncFeiShuDepts: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
	}

	// 2.将远程数据转换成树
	deptTree := GroupListToTree(fmt.Sprintf("%s_0", config.Conf.FeiShu.Flag), depts)

	// 3.根据树进行创建
	err = d.addDepts(deptTree.Children)
	if err != nil {
		errMsg := fmt.Sprintf("创建飞书部门失败：%s", err.Error())
		common.Log.Errorf("SyncFeiShuDepts: %s", errMsg)
		return nil, err
	}

	common.Log.Infof("SyncFeiShuDepts: 飞书部门同步成功")
	return nil, err
}

// 添加部门
func (d FeiShuLogic) addDepts(depts []*model.Group) error {
	for _, dept := range depts {
		err := d.AddDepts(dept)
		if err != nil {
			errMsg := fmt.Sprintf("DsyncFeiShuDepts添加部门[%s]失败: %s", dept.GroupName, err.Error())
			common.Log.Errorf("%s", errMsg)
			return tools.NewOperationError(fmt.Errorf(errMsg))
		}
		if len(dept.Children) != 0 {
			err = d.addDepts(dept.Children)
			if err != nil {
				errMsg := fmt.Sprintf("DsyncFeiShuDepts添加子部门失败: %s", err.Error())
				common.Log.Errorf("%s", errMsg)
				return tools.NewOperationError(fmt.Errorf(errMsg))
			}
		}
	}
	return nil
}

// AddGroup 添加部门数据
func (d FeiShuLogic) AddDepts(group *model.Group) error {
	// 查询当前分组父ID在MySQL中的数据信息
	parentGroup := new(model.Group)
	err := isql.Group.Find(tools.H{"source_dept_id": group.SourceDeptParentId}, parentGroup)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("查询父级部门失败：%s", err.Error()))
	}

	// 此时的 group 已经附带了Build后动态关联好的字段，接下来将一些确定性的其他字段值添加上，就可以创建这个分组了
	group.Creator = "system"
	group.GroupType = "cn"
	group.ParentId = parentGroup.ID
	group.Source = config.Conf.FeiShu.Flag
	group.GroupDN = fmt.Sprintf("cn=%s,%s", group.GroupName, parentGroup.GroupDN)

	if !isql.Group.Exist(tools.H{"group_dn": group.GroupDN}) {
		err = CommonAddGroup(group)
		if err != nil {
			return tools.NewOperationError(fmt.Errorf("添加部门: %s, 失败: %s", group.GroupName, err.Error()))
		}
	}
	return nil
}

// 根据现有数据库同步到的部门信息，开启用户同步
func (d FeiShuLogic) SyncFeiShuUsers(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	// 1.获取飞书用户列表
	staffSource, err := feishu.GetAllUsers()
	if err != nil {
		errMsg := fmt.Sprintf("获取飞书用户列表失败：%s", err.Error())
		common.Log.Errorf("SyncFeiShuUsers: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
	}
	staffs, err := ConvertUserData(config.Conf.FeiShu.Flag, staffSource)
	if err != nil {
		errMsg := fmt.Sprintf("转换飞书用户数据失败：%s", err.Error())
		common.Log.Errorf("SyncFeiShuUsers: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
	}
	// 2.遍历用户，开始写入
	for i, staff := range staffs {
		// 入库
		err = d.AddUsers(staff)
		if err != nil {
			errMsg := fmt.Sprintf("写入用户[%s]失败：%s", staff.Username, err.Error())
			common.Log.Errorf("SyncFeiShuUsers: %s", errMsg)
			return nil, tools.NewOperationError(fmt.Errorf(errMsg))
		}
		common.Log.Infof("SyncFeiShuUsers: 成功同步用户[%s] (%d/%d)", staff.Username, i+1, len(staffs))
	}

	// 3.获取飞书已离职用户id列表
	userIds, err := feishu.GetLeaveUserIds()
	if err != nil {
		errMsg := fmt.Sprintf("获取飞书离职用户列表失败：%s", err.Error())
		common.Log.Errorf("SyncFeiShuUsers: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
	}
	// 4.遍历id，开始处理
	processedCount := 0
	for _, uid := range userIds {
		if isql.User.Exist(
			tools.H{
				"status":          1, //只处理1在职的
				"source_union_id": fmt.Sprintf("%s_%s", config.Conf.FeiShu.Flag, uid),
			}) {
			user := new(model.User)
			err = isql.User.Find(tools.H{"source_union_id": fmt.Sprintf("%s_%s", config.Conf.FeiShu.Flag, uid)}, user)
			if err != nil {
				errMsg := fmt.Sprintf("在MySQL查询离职用户[%s]失败: %s", uid, err.Error())
				common.Log.Errorf("SyncFeiShuUsers: %s", errMsg)
				return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
			}
			// 先从ldap删除用户
			err = ildap.User.Delete(user.UserDN)
			if err != nil {
				errMsg := fmt.Sprintf("在LDAP删除离职用户[%s]失败: %s", user.Username, err.Error())
				common.Log.Errorf("SyncFeiShuUsers: %s", errMsg)
				return nil, tools.NewLdapError(fmt.Errorf(errMsg))
			}
			// 然后更新MySQL中用户状态
			err = isql.User.ChangeStatus(int(user.ID), 2)
			if err != nil {
				errMsg := fmt.Sprintf("在MySQL更新离职用户[%s]状态失败: %s", user.Username, err.Error())
				common.Log.Errorf("SyncFeiShuUsers: %s", errMsg)
				return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
			}
			processedCount++
			common.Log.Infof("SyncFeiShuUsers: 成功处理离职用户[%s]", user.Username)
		}
	}

	common.Log.Infof("SyncFeiShuUsers: 飞书用户同步完成，共同步%d个在职用户，处理%d个离职用户", len(staffs), processedCount)
	return nil, nil
}

// AddUser 添加用户数据
func (d FeiShuLogic) AddUsers(user *model.User) error {
	// 根据角色id获取角色
	roles, err := isql.Role.GetRolesByIds([]uint{2})
	if err != nil {
		return tools.NewValidatorError(fmt.Errorf("根据角色ID获取角色信息失败:%s", err.Error()))
	}
	user.Roles = roles
	user.Creator = "system"
	user.Source = config.Conf.FeiShu.Flag
	user.Password = config.Conf.Ldap.UserInitPassword
	user.UserDN = fmt.Sprintf("uid=%s,%s", user.Username, config.Conf.Ldap.UserDN)

	// 根据 user_dn 查询用户,不存在则创建
	if !isql.User.Exist(tools.H{"user_dn": user.UserDN}) {
		// 获取用户将要添加的分组
		groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("根据部门ID获取部门信息失败" + err.Error()))
		}
		var deptTmp string
		for _, group := range groups {
			deptTmp = deptTmp + group.GroupName + ","
		}
		user.Departments = strings.TrimRight(deptTmp, ",")

		// 添加用户
		err = CommonAddUser(user, groups)
		if err != nil {
			return tools.NewOperationError(fmt.Errorf("添加用户: %s, 失败: %s", user.Username, err.Error()))
		}
	} else {
		// 此处逻辑未经实际验证，如在使用中有问题，请反馈
		if config.Conf.FeiShu.IsUpdateSyncd {
			// 先获取用户信息
			oldData := new(model.User)
			err = isql.User.Find(tools.H{"user_dn": user.UserDN}, oldData)
			if err != nil {
				return err
			}
			// 获取用户将要添加的分组
			groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
			if err != nil {
				return tools.NewMySqlError(fmt.Errorf("根据部门ID获取部门信息失败" + err.Error()))
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

			// 用户信息的预置处理
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
