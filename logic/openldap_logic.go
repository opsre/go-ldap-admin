package logic

import (
	"fmt"
	"strings"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/client/openldap"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/isql"
	"github.com/gin-gonic/gin"
)

type OpenLdapLogic struct {
}

//通过ldap获取部门信息
func (d *OpenLdapLogic) SyncOpenLdapDepts(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	// 1.获取所有部门
	depts, err := openldap.GetAllDepts()
	if err != nil {
		errMsg := fmt.Sprintf("获取OpenLDAP部门列表失败：%s", err.Error())
		common.Log.Errorf("SyncOpenLdapDepts: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
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
	// 2.将远程数据转换成树
	deptTree := GroupListToTree("0", groups)

	// 3.根据树进行创建
	err = d.addDepts(deptTree.Children)
	if err != nil {
		errMsg := fmt.Sprintf("创建OpenLDAP部门失败：%s", err.Error())
		common.Log.Errorf("SyncOpenLdapDepts: %s", errMsg)
		return nil, err
	}

	common.Log.Infof("SyncOpenLdapDepts: OpenLDAP部门同步成功")
	return nil, err
}

// 添加部门
func (d OpenLdapLogic) addDepts(depts []*model.Group) error {
	for _, dept := range depts {
		err := d.AddDepts(dept)
		if err != nil {
			errMsg := fmt.Sprintf("DsyncOpenLdapDepts添加部门[%s]失败: %s", dept.GroupName, err.Error())
			common.Log.Errorf("%s", errMsg)
			return tools.NewOperationError(fmt.Errorf(errMsg))
		}
		if len(dept.Children) != 0 {
			err = d.addDepts(dept.Children)
			if err != nil {
				errMsg := fmt.Sprintf("DsyncOpenLdapDepts添加子部门失败: %s", err.Error())
				common.Log.Errorf("%s", errMsg)
				return tools.NewOperationError(fmt.Errorf(errMsg))
			}
		}
	}
	return nil
}

// AddGroup 添加部门数据
func (d OpenLdapLogic) AddDepts(group *model.Group) error {
	// 判断部门名称是否存在,此处使用ldap中的唯一值dn,以免出现数据同步不全的问题
	if !isql.Group.Exist(tools.H{"group_dn": group.GroupDN}) {
		// 此时的 group 已经附带了Build后动态关联好的字段，接下来将一些确定性的其他字段值添加上，就可以创建这个分组了
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

// AddGroup 添加部门数据
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
		return id, tools.NewMySqlError(fmt.Errorf("查询父级部门失败：%s,%s", err.Error(), group.GroupName))
	}
	return parentGroup.ID, nil
}

//根据现有数据库同步到的部门信息，开启用户同步
func (d OpenLdapLogic) SyncOpenLdapUsers(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	// 1.获取ldap用户列表
	staffs, err := openldap.GetAllUsers()
	if err != nil {
		errMsg := fmt.Sprintf("获取OpenLDAP用户列表失败：%s", err.Error())
		common.Log.Errorf("SyncOpenLdapUsers: %s", errMsg)
		return nil, tools.NewOperationError(fmt.Errorf(errMsg))
	}
	// 2.遍历用户，开始写入
	for i, staff := range staffs {
		groupIds, err := isql.Group.DeptIdsToGroupIds(staff.DepartmentIds)
		if err != nil {
			errMsg := fmt.Sprintf("将用户[%s]的部门ids转换为内部部门id失败：%s", staff.Name, err.Error())
			common.Log.Errorf("SyncOpenLdapUsers: %s", errMsg)
			return nil, tools.NewMySqlError(fmt.Errorf(errMsg))
		}
		// 根据角色id获取角色
		roles, err := isql.Role.GetRolesByIds([]uint{2})
		if err != nil {
			errMsg := fmt.Sprintf("获取用户[%s]的角色信息失败：%s", staff.Name, err.Error())
			common.Log.Errorf("SyncOpenLdapUsers: %s", errMsg)
			return nil, tools.NewValidatorError(fmt.Errorf(errMsg))
		}
		// 入库
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
			errMsg := fmt.Sprintf("写入用户[%s]失败：%s", staff.Name, err.Error())
			common.Log.Errorf("SyncOpenLdapUsers: %s", errMsg)
			return nil, tools.NewOperationError(fmt.Errorf(errMsg))
		}
		common.Log.Infof("SyncOpenLdapUsers: 成功同步用户[%s] (%d/%d)", staff.Name, i+1, len(staffs))
	}
	
	common.Log.Infof("SyncOpenLdapUsers: OpenLDAP用户同步完成，共同步%d个用户", len(staffs))
	return nil, nil
}

// AddUser 添加用户数据
func (d OpenLdapLogic) AddUsers(user *model.User) error {
	// 根据 user_dn 查询用户,不存在则创建
	if !isql.User.Exist(tools.H{"user_dn": user.UserDN}) {
		if user.Departments == "" {
			user.Departments = "默认:研发中心"
		}
		if user.GivenName == "" {
			user.GivenName = user.Nickname
		}
		if user.PostalAddress == "" {
			user.PostalAddress = "默认:地球"
		}
		if user.Position == "" {
			user.Position = "默认:技术"
		}
		if user.Introduction == "" {
			user.Introduction = user.Nickname
		}
		if user.JobNumber == "" {
			user.JobNumber = "未启用"
		}
		// 先将用户添加到MySQL
		err := isql.User.Add(user)
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("向MySQL创建用户失败：" + err.Error()))
		}

		// 获取用户将要添加的分组
		groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("根据部门ID获取部门信息失败" + err.Error()))
		}
		for _, group := range groups {
			if group.GroupDN[:3] == "ou=" {
				continue
			}
			// 先将用户和部门信息维护到MySQL
			err := isql.Group.AddUserToGroup(group, []model.User{*user})
			if err != nil {
				return tools.NewMySqlError(fmt.Errorf("向MySQL添加用户到分组关系失败：" + err.Error()))
			}
		}
		return nil
	}
	return nil
}
