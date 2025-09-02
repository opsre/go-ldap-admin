package logic

import (
	"fmt"
	"strings"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/client/wechat"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"
	"github.com/gin-gonic/gin"
)

type WeComLogic struct {
}

// 通过企业微信获取部门信息
func (d *WeComLogic) SyncWeComDepts(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	// 1.获取所有部门
	deptSource, err := wechat.GetAllDepts()
	if err != nil {
		errMsg := fmt.Sprintf("获取企业微信部门列表失败：%s", err.Error())
		common.Log.Errorf("SyncWeComDepts: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
	}
	depts, err := ConvertDeptData(config.Conf.WeCom.Flag, deptSource)
	if err != nil {
		errMsg := fmt.Sprintf("转换企业微信部门数据失败：%s", err.Error())
		common.Log.Errorf("SyncWeComDepts: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
	}

	// 2.将远程数据转换成树
	deptTree := GroupListToTree(fmt.Sprintf("%s_1", config.Conf.WeCom.Flag), depts)

	// 3.根据树进行创建
	err = d.addDepts(deptTree.Children)
	if err != nil {
		errMsg := fmt.Sprintf("创建企业微信部门失败：%s", err.Error())
		common.Log.Errorf("SyncWeComDepts: %s", errMsg)
		return nil, err
	}

	common.Log.Infof("SyncWeComDepts: 企业微信部门同步成功")
	return nil, err
}

// 添加部门
func (d WeComLogic) addDepts(depts []*model.Group) error {
	for _, dept := range depts {
		err := d.AddDepts(dept)
		if err != nil {
			errMsg := fmt.Sprintf("DsyncWeComDepts添加部门[%s]失败: %s", dept.GroupName, err.Error())
			common.Log.Errorf("%s", errMsg)
			return tools.NewOperationError(fmt.Errorf(errMsg))
		}
		if len(dept.Children) != 0 {
			err = d.addDepts(dept.Children)
			if err != nil {
				errMsg := fmt.Sprintf("DsyncWeComDepts添加子部门失败: %s", err.Error())
				common.Log.Errorf("%s", errMsg)
				return tools.NewOperationError(fmt.Errorf(errMsg))
			}
		}
	}
	return nil
}

// AddGroup 添加部门数据
func (d WeComLogic) AddDepts(group *model.Group) error {
	// 判断部门名称是否存在
	parentGroup := new(model.Group)
	err := isql.Group.Find(tools.H{"source_dept_id": group.SourceDeptParentId}, parentGroup)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("查询父级部门失败：%s", err.Error()))
	}

	// 此时的 group 已经附带了Build后动态关联好的字段，接下来将一些确定性的其他字段值添加上，就可以创建这个分组了
	group.Creator = "system"
	group.GroupType = "cn"
	group.ParentId = parentGroup.ID
	group.Source = config.Conf.WeCom.Flag
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
func (d WeComLogic) SyncWeComUsers(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	// 1.获取企业微信用户列表
	staffSource, err := wechat.GetAllUsers()
	if err != nil {
		errMsg := fmt.Sprintf("获取企业微信用户列表失败：%s", err.Error())
		common.Log.Errorf("SyncWeComUsers: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
	}
	staffs, err := ConvertUserData(config.Conf.WeCom.Flag, staffSource)
	if err != nil {
		errMsg := fmt.Sprintf("转换企业微信用户数据失败：%s", err.Error())
		common.Log.Errorf("SyncWeComUsers: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
	}
	// 2.遍历用户，开始写入
	for i, staff := range staffs {
		// 入库
		err = d.AddUsers(staff)
		if err != nil {
			errMsg := fmt.Sprintf("写入用户[%s]失败：%s", staff.Username, err.Error())
			common.Log.Errorf("SyncWeComUsers: %s", errMsg)
			return nil, tools.NewOperationError(fmt.Errorf(errMsg))
		}
		common.Log.Infof("SyncWeComUsers: 成功同步用户[%s] (%d/%d)", staff.Username, i+1, len(staffs))
	}

	// 3.获取企业微信已离职用户id列表
	// 拿到MySQL所有用户数据(来源为 wecom的用户)，远程没有的，则说明被删除了
	// 如果以后企业微信透出了已离职用户列表的接口，则这里可以进行改进
	var res []*model.User
	users, err := isql.User.ListAll()
	if err != nil {
		errMsg := fmt.Sprintf("获取MySQL用户列表失败：%s", err.Error())
		common.Log.Errorf("SyncWeComUsers: %s", errMsg)
		return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
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
	// 4.遍历id，开始处理
	processedCount := 0
	for _, userTmp := range res {
		user := new(model.User)
		err = isql.User.Find(tools.H{"source_user_id": userTmp.SourceUserId, "status": 1}, user)
		if err != nil {
			errMsg := fmt.Sprintf("在MySQL查询离职用户[%s]失败: %s", userTmp.Username, err.Error())
			common.Log.Errorf("SyncWeComUsers: %s", errMsg)
			return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
		}
		// 先从ldap删除用户
		err = ildap.User.Delete(user.UserDN)
		if err != nil {
			errMsg := fmt.Sprintf("在LDAP删除离职用户[%s]失败: %s", user.Username, err.Error())
			common.Log.Errorf("SyncWeComUsers: %s", errMsg)
			return nil, tools.NewLdapError(fmt.Errorf(errMsg))
		}
		// 然后更新MySQL中用户状态
		err = isql.User.ChangeStatus(int(user.ID), 2)
		if err != nil {
			errMsg := fmt.Sprintf("在MySQL更新离职用户[%s]状态失败: %s", user.Username, err.Error())
			common.Log.Errorf("SyncWeComUsers: %s", errMsg)
			return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
		}
		processedCount++
		common.Log.Infof("SyncWeComUsers: 成功处理离职用户[%s]", user.Username)
	}
	
	common.Log.Infof("SyncWeComUsers: 企业微信用户同步完成，共同步%d个在职用户，处理%d个离职用户", len(staffs), processedCount)
	return nil, nil
}

// AddUser 添加用户数据
func (d WeComLogic) AddUsers(user *model.User) error {
	// 根据角色id获取角色
	roles, err := isql.Role.GetRolesByIds([]uint{2})
	if err != nil {
		return tools.NewValidatorError(fmt.Errorf("根据角色ID获取角色信息失败:%s", err.Error()))
	}
	user.Creator = "system"
	user.Roles = roles
	user.Password = config.Conf.Ldap.UserInitPassword
	user.Source = config.Conf.WeCom.Flag
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

		// 创建用户
		err = CommonAddUser(user, groups)
		if err != nil {
			return tools.NewOperationError(fmt.Errorf("添加用户: %s, 失败: %s", user.Username, err.Error()))
		}
	} else {
		// 此处逻辑未经实际验证，如在使用中有问题，请反馈
		if config.Conf.WeCom.IsUpdateSyncd {
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
