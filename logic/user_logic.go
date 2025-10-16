package logic

import (
	"fmt"
	"strings"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/model/response"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type UserLogic struct{}

// Add 添加数据
func (l UserLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if isql.User.Exist(tools.H{"username": r.Username}) {
		return nil, tools.NewValidatorError(fmt.Errorf("用户名已存在,请勿重复添加"))
	}
	if isql.User.Exist(tools.H{"mobile": r.Mobile}) {
		return nil, tools.NewValidatorError(fmt.Errorf("手机号已存在,请勿重复添加"))
	}
	if isql.User.Exist(tools.H{"job_number": r.JobNumber}) {
		return nil, tools.NewValidatorError(fmt.Errorf("工号已存在,请勿重复添加"))
	}
	if isql.User.Exist(tools.H{"mail": r.Mail}) {
		return nil, tools.NewValidatorError(fmt.Errorf("邮箱已存在,请勿重复添加"))
	}

	// 密码通过RSA解密
	// 密码不为空就解密
	if r.Password != "" {
		decodeData, err := tools.RSADecrypt([]byte(r.Password), config.Conf.System.RSAPrivateBytes)
		if err != nil {
			return nil, tools.NewValidatorError(fmt.Errorf("密码解密失败"))
		}
		r.Password = string(decodeData)
		if len(r.Password) < 6 {
			return nil, tools.NewValidatorError(fmt.Errorf("密码长度至少为6位"))
		}
	} else {
		r.Password = config.Conf.Ldap.UserInitPassword
	}

	// 当前登陆用户角色排序最小值（最高等级角色）以及当前登陆的用户
	currentRoleSortMin, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("获取当前登陆用户角色排序最小值失败"))
	}

	// 根据角色id获取角色
	if len(r.RoleIds) == 0 {
		r.RoleIds = []uint{2} // 默认添加为普通用户角色
	}

	roles, err := isql.Role.GetRolesByIds(r.RoleIds)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("根据角色ID获取角色信息失败"))
	}

	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := uint(funk.MinInt(reqRoleSorts).(int))

	// 如果登录用户的角色ID为1，亦即为管理员，则直接放行，保障管理员拥有最大权限
	if currentRoleSortMin != 1 {
		// 当前用户的角色排序最小值 需要小于 前端传来的角色排序最小值（用户不能创建比自己等级高的或者相同等级的用户）
		if currentRoleSortMin >= reqRoleSortMin {
			return nil, tools.NewValidatorError(fmt.Errorf("用户不能创建比自己等级高的或者相同等级的用户"))
		}
	}
	user := model.User{
		Username:      r.Username,
		Password:      r.Password,
		Nickname:      r.Nickname,
		GivenName:     r.GivenName,
		Mail:          r.Mail,
		JobNumber:     r.JobNumber,
		Mobile:        r.Mobile,
		Avatar:        r.Avatar,
		PostalAddress: r.PostalAddress,
		Departments:   r.Departments,
		Position:      r.Position,
		Introduction:  r.Introduction,
		Status:        r.Status,
		Creator:       ctxUser.Username,
		DepartmentId:  tools.SliceToString(r.DepartmentId, ","),
		Source:        r.Source,
		Roles:         roles,
		UserDN:        fmt.Sprintf("uid=%s,%s", r.Username, config.Conf.Ldap.UserDN),
	}

	if user.Source == "" {
		user.Source = "platform"
	}

	// 获取用户将要添加的分组
	groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "根据部门ID获取部门信息失败"+err.Error()))
	}

	err = CommonAddUser(&user, groups)
	if err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("%s", "添加用户失败"+err.Error()))
	}
	return nil, nil
}

// List 数据列表
func (l UserLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	users, err := isql.User.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "获取用户列表失败："+err.Error()))
	}

	rets := make([]model.User, 0)
	for _, user := range users {
		rets = append(rets, *user)
	}
	count, err := isql.User.ListCount(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "获取用户总数失败："+err.Error()))
	}

	return response.UserListRsp{
		Total: int(count),
		Users: rets,
	}, nil
}

// Update 更新数据
func (l UserLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if !isql.User.Exist(tools.H{"id": r.ID}) {
		return nil, tools.NewMySqlError(fmt.Errorf("该记录不存在"))
	}

	// 获取当前登陆用户
	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户失败"))
	}

	// 获取当前登陆用户角色ID集合
	var currentRoleSorts []int
	for _, role := range ctxUser.Roles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}

	// 获取将要操作的用户角色ID集合
	var reqRoleSorts []int
	roles, _ := isql.Role.GetRolesByIds(r.RoleIds)
	if len(roles) == 0 {
		return nil, tools.NewValidatorError(fmt.Errorf("根据角色ID获取角色信息失败"))
	}
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}

	// 当前登陆用户角色排序最小值（最高等级角色）
	currentRoleSortMin := funk.MinInt(currentRoleSorts).(int)
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := funk.MinInt(reqRoleSorts).(int)

	// 如果登录用户的角色ID为1，亦即为管理员，则直接放行，保障管理员拥有最大权限
	if currentRoleSortMin != 1 {
		// 判断是更新自己还是更新别人,如果操作的ID与登陆用户的ID一致，则说明操作的是自己
		if int(r.ID) == int(ctxUser.ID) {
			// 不能更改自己的角色
			reqDiff, currentDiff := funk.Difference(reqRoleSorts, currentRoleSorts)
			if len(reqDiff.([]int)) > 0 || len(currentDiff.([]int)) > 0 {
				return nil, tools.NewValidatorError(fmt.Errorf("不能更改自己的角色"))
			}
		}

		// 如果是更新别人，操作者不能更新比自己角色等级高的或者相同等级的用户
		minRoleSorts, err := isql.User.GetUserMinRoleSortsByIds([]uint{uint(r.ID)}) // 根据userIdID获取用户角色排序最小值
		if err != nil || len(minRoleSorts) == 0 {
			return nil, tools.NewValidatorError(fmt.Errorf("根据用户ID获取用户角色排序最小值失败"))
		}
		if currentRoleSortMin >= minRoleSorts[0] || currentRoleSortMin >= reqRoleSortMin {
			return nil, tools.NewValidatorError(fmt.Errorf("用户不能更新比自己角色等级高的或者相同等级的用户"))
		}
	}

	// 先获取用户信息
	oldData := new(model.User)
	err = isql.User.Find(tools.H{"id": r.ID}, oldData)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}

	// 过滤掉前端会选择到的 请选择部门信息 这个选项
	var (
		depts   string
		deptids []uint
	)
	for _, v := range strings.Split(r.Departments, ",") {
		if v != "请选择部门信息" {
			depts += v + ","
		}
	}
	for _, j := range r.DepartmentId {
		if j != 0 {
			deptids = append(deptids, j)
		}
	}

	// 拼装新的用户信息
	user := model.User{
		Model:         oldData.Model,
		Username:      r.Username,
		Nickname:      r.Nickname,
		GivenName:     r.GivenName,
		Mail:          r.Mail,
		JobNumber:     r.JobNumber,
		Mobile:        r.Mobile,
		Avatar:        r.Avatar,
		PostalAddress: r.PostalAddress,
		Departments:   depts,
		Position:      r.Position,
		Introduction:  r.Introduction,
		Creator:       ctxUser.Username,
		DepartmentId:  tools.SliceToString(deptids, ","),
		Source:        oldData.Source,
		Roles:         roles,
		UserDN:        oldData.UserDN,
	}

	if err = CommonUpdateUser(oldData, &user, r.DepartmentId); err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("%s", "更新用户失败"+err.Error()))
	}

	return nil, nil
}

// Delete 删除数据
func (l UserLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.UserIds {
		filter := tools.H{"id": int(id)}
		if !isql.User.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("有用户不存在"))
		}
	}

	// 根据用户ID获取用户角色排序最小值
	roleMinSortList, err := isql.User.GetUserMinRoleSortsByIds(r.UserIds)
	if err != nil || len(roleMinSortList) == 0 {
		return nil, tools.NewValidatorError(fmt.Errorf("根据用户ID获取用户角色排序最小值失败"))
	}

	// 获取当前登陆用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("获取当前登陆用户角色排序最小值失败"))
	}

	// 不能删除自己
	if funk.Contains(r.UserIds, ctxUser.ID) {
		return nil, tools.NewValidatorError(fmt.Errorf("用户不能删除自己"))
	}

	// 不能删除比自己(登陆用户)角色排序低(等级高)的用户
	for _, sort := range roleMinSortList {
		if int(minSort) > sort {
			return nil, tools.NewValidatorError(fmt.Errorf("用户不能删除比自己角色等级高的用户"))
		}
	}

	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "获取用户信息失败: "+err.Error()))
	}

	// 先将用户从ldap中删除
	for _, user := range users {
		err := ildap.User.Delete(user.UserDN)
		if err != nil {
			return nil, tools.NewLdapError(fmt.Errorf("%s", "在LDAP删除用户失败"+err.Error()))
		}
	}

	// 再将用户从MySQL中删除
	err = isql.User.Delete(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "在MySQL删除用户失败: "+err.Error()))
	}

	return nil, nil
}

// ChangePwd 修改密码
func (l UserLogic) ChangePwd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserChangePwdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// 前端传来的密码是rsa加密的,先解密
	// 密码通过RSA解密
	decodeOldPassword, err := tools.RSADecrypt([]byte(r.OldPassword), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("原密码解析失败"))
	}
	decodeNewPassword, err := tools.RSADecrypt([]byte(r.NewPassword), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("新密码解析失败"))
	}
	r.OldPassword = string(decodeOldPassword)
	r.NewPassword = string(decodeNewPassword)
	// 获取当前用户
	user, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户失败"))
	}
	// 获取用户的真实正确密码
	// correctPasswd := user.Password
	// 判断前端请求的密码是否等于真实密码
	// err = tools.ComparePasswd(correctPasswd, r.OldPassword)
	// if err != nil {
	// 	return nil, tools.NewValidatorError(fmt.Errorf("原密码错误"))
	// }
	if tools.NewParPasswd(user.Password) != r.OldPassword {
		return nil, tools.NewValidatorError(fmt.Errorf("原密码错误"))
	}

	err = tools.CheckPasswdStrength(r.NewPassword)
	if err != nil {
		return nil, tools.NewValidatorError(err)
	}

	// ldap更新密码时可以直接指定用户DN和新密码即可更改成功
	err = ildap.User.ChangePwd(user.UserDN, "", r.NewPassword)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("%s", "在LDAP更新密码失败"+err.Error()))
	}

	// 更新密码
	err = isql.User.ChangePwd(user.Username, tools.NewGenPasswd(r.NewPassword))
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "在MySQL更新密码失败: "+err.Error()))
	}

	return nil, nil
}

// ResetPassword 重置用户密码
func (l UserLogic) ResetPassword(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserResetPasswordReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 检查用户是否存在
	if !isql.User.Exist(tools.H{"username": r.Username}) {
		return nil, tools.NewValidatorError(fmt.Errorf("用户不存在"))
	}

	// 获取用户信息
	user := new(model.User)
	err := isql.User.Find(tools.H{"username": r.Username}, user)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取用户信息失败: %s", err.Error()))
	}

	// 生成随机密码
	newPassword := tools.GenerateRandomPassword()

	// 在LDAP中更新密码
	err = ildap.User.ChangePwd(user.UserDN, "", newPassword)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("在LDAP更新密码失败: %s", err.Error()))
	}

	// 在MySQL中更新密码
	err = isql.User.ChangePwd(user.Username, tools.NewGenPasswd(newPassword))
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("在MySQL更新密码失败: %s", err.Error()))
	}

	// 发送密码重置通知邮件
	if err := tools.SendPasswordResetNotification(user.Username, user.Nickname, user.Mail, newPassword); err != nil {
		common.Log.Warnf("发送密码重置通知邮件失败，用户: %s, 邮箱: %s, 错误: %v", user.Username, user.Mail, err)
	}

	return map[string]string{"newPassword": newPassword}, nil
}

// ChangeUserStatus 修改用户状态
func (l UserLogic) ChangeUserStatus(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserChangeUserStatusReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// 校验工作
	filter := tools.H{"id": r.ID}
	if !isql.User.Exist(filter) {
		return nil, tools.NewValidatorError(fmt.Errorf("该用户不存在"))
	}
	user := new(model.User)
	err := isql.User.Find(filter, user)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "在MySQL查询用户失败: "+err.Error()))
	}

	if r.Status == 1 && r.Status == user.Status {
		return nil, tools.NewValidatorError(fmt.Errorf("用户已经是在职状态"))
	}
	if r.Status == 2 && r.Status == user.Status {
		return nil, tools.NewValidatorError(fmt.Errorf("用户已经是离职状态"))
	}

	// 获取当前登录用户，只有管理员才能够将用户状态改变
	// 获取当前登陆用户角色排序最小值（最高等级角色）以及当前用户
	minSort, _, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewValidatorError(fmt.Errorf("获取当前登陆用户角色排序最小值失败"))
	}

	if int(minSort) != 1 {
		return nil, tools.NewValidatorError(fmt.Errorf("只有管理员才能更改用户状态"))
	}

	if r.Status == 2 {
		err = ildap.User.Delete(user.UserDN)
		if err != nil {
			return nil, tools.NewLdapError(fmt.Errorf("%s", "在LDAP删除用户失败"+err.Error()))
		}
	} else {
		err = ildap.User.Add(user)
		if err != nil {
			return nil, tools.NewLdapError(fmt.Errorf("%s", "在LDAP添加用户失败"+err.Error()))
		}
	}
	err = isql.User.ChangeStatus(int(r.ID), int(r.Status))
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "在MySQL更新用户状态失败: "+err.Error()))
	}
	return nil, nil
}

// GetUserInfo 获取用户信息
func (l UserLogic) GetUserInfo(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserGetUserInfoReq)
	if !ok {
		return nil, ReqAssertErr
	}

	_ = c
	_ = r

	user, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "获取当前用户信息失败: "+err.Error()))
	}
	return user, nil
}
