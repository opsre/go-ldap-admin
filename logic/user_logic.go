package logic

import (
	"fmt"
	"strings"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/model/response"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type UserLogic struct{}

// Add data
func (l UserLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if isql.User.Exist(tools.H{"username": r.Username}) {
		return nil, tools.NewValidatorI18nError("user.username_exists", nil)
	}
	if isql.User.Exist(tools.H{"mobile": r.Mobile}) {
		return nil, tools.NewValidatorI18nError("user.mobile_exists", nil)
	}
	if isql.User.Exist(tools.H{"job_number": r.JobNumber}) {
		return nil, tools.NewValidatorI18nError("user.job_number_exists", nil)
	}
	if isql.User.Exist(tools.H{"mail": r.Mail}) {
		return nil, tools.NewValidatorI18nError("user.mail_exists", nil)
	}

	// Decrypt the password with RSA
	// Decrypt the password if it is not empty
	if r.Password != "" {
		decodeData, err := tools.RSADecrypt([]byte(r.Password), config.Conf.System.RSAPrivateBytes)
		if err != nil {
			return nil, tools.NewValidatorI18nError("user.password_decrypt_failed", nil)
		}
		r.Password = string(decodeData)
		if len(r.Password) < 6 {
			return nil, tools.NewValidatorI18nError("user.password_min_length", i18n.Args{"min": 6})
		}
	} else {
		r.Password = config.Conf.Ldap.UserInitPassword
	}

	// The currently logged-in user and their minimum role sort value (highest role level)
	currentRoleSortMin, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewValidatorI18nError("user.current_min_role_sort_failed", nil)
	}

	// Get the role by role ID
	if len(r.RoleIds) == 0 {
		r.RoleIds = []uint{2} // By default, assign the regular user role
	}

	roles, err := isql.Role.GetRolesByIds(r.RoleIds)
	if err != nil {
		return nil, tools.NewValidatorI18nError("user.role_info_failed", nil)
	}

	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// The minimum role sort value (highest role level) sent by the frontend
	reqRoleSortMin := uint(funk.MinInt(reqRoleSorts).(int))

	// If the logged-in user has role ID 1 (admin), allow the request so the admin retains full permissions
	if currentRoleSortMin != 1 {
		// The current user minimum role sort value must be less than the value sent by the frontend (users cannot create users at or above their own level)
		if currentRoleSortMin >= reqRoleSortMin {
			return nil, tools.NewValidatorI18nError("user.cannot_create_higher_or_equal", nil)
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

	// Get the groups to add the user to
	groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.department_info_failed", i18n.Args{"error": err.Error()})
	}

	err = CommonAddUser(&user, groups, i18n.LocaleFromContext(c))
	if err != nil {
		return nil, tools.NewOperationI18nError("user.create_failed", i18n.Args{"error": err.Error()})
	}
	return nil, nil
}

// List data
func (l UserLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	users, err := isql.User.List(r)
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.list_failed", i18n.Args{"error": err.Error()})
	}

	rets := make([]model.User, 0)
	for _, user := range users {
		rets = append(rets, *user)
	}
	count, err := isql.User.ListCount(r)
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.count_failed", i18n.Args{"error": err.Error()})
	}

	return response.UserListRsp{
		Total: int(count),
		Users: rets,
	}, nil
}

// Update data
func (l UserLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if !isql.User.Exist(tools.H{"id": r.ID}) {
		return nil, tools.NewMySqlI18nError("legacy.common.record_not_found", nil)
	}

	// Get the currently logged-in user
	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.current_user_failed", nil)
	}

	// Get the role IDs of the currently logged-in user
	var currentRoleSorts []int
	for _, role := range ctxUser.Roles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}

	// Get the role IDs of the user to be operated on
	var reqRoleSorts []int
	roles, _ := isql.Role.GetRolesByIds(r.RoleIds)
	if len(roles) == 0 {
		return nil, tools.NewValidatorI18nError("user.role_info_failed", nil)
	}
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}

	// The minimum role sort value (highest role level) of the currently logged-in user
	currentRoleSortMin := funk.MinInt(currentRoleSorts).(int)
	// The minimum role sort value (highest role level) sent by the frontend
	reqRoleSortMin := funk.MinInt(reqRoleSorts).(int)

	// If the logged-in user has role ID 1 (admin), allow the request so the admin retains full permissions
	if currentRoleSortMin != 1 {
		// Check whether updating self or another user; matching IDs means it is the current user
		if int(r.ID) == int(ctxUser.ID) {
			// Cannot change your own role
			reqDiff, currentDiff := funk.Difference(reqRoleSorts, currentRoleSorts)
			if len(reqDiff.([]int)) > 0 || len(currentDiff.([]int)) > 0 {
				return nil, tools.NewValidatorI18nError("user.cannot_change_self_role", nil)
			}
		}

		// When updating another user, the operator cannot update a user at or above their own role level
		minRoleSorts, err := isql.User.GetUserMinRoleSortsByIds([]uint{uint(r.ID)}) // Get the minimum role sort value for a user by userIdID
		if err != nil || len(minRoleSorts) == 0 {
			return nil, tools.NewValidatorI18nError("user.min_role_sort_by_id_failed", nil)
		}
		if currentRoleSortMin >= minRoleSorts[0] || currentRoleSortMin >= reqRoleSortMin {
			return nil, tools.NewValidatorI18nError("user.cannot_update_higher_or_equal", nil)
		}
	}

	// First get the user information
	oldData := new(model.User)
	err = isql.User.Find(tools.H{"id": r.ID}, oldData)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}

	// Filter out the frontend option "Please select department information"
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

	// Assemble the new user information
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
		return nil, tools.NewOperationI18nError("user.update_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// Delete data
func (l UserLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.UserIds {
		filter := tools.H{"id": int(id)}
		if !isql.User.Exist(filter) {
			return nil, tools.NewMySqlI18nError("legacy.user.some_not_found", nil)
		}
	}

	// Get the minimum role sort value by user ID
	roleMinSortList, err := isql.User.GetUserMinRoleSortsByIds(r.UserIds)
	if err != nil || len(roleMinSortList) == 0 {
		return nil, tools.NewValidatorI18nError("user.min_role_sort_by_id_failed", nil)
	}

	// Get the current user and their minimum role sort value (highest role level)
	minSort, ctxUser, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewValidatorI18nError("user.current_min_role_sort_failed", nil)
	}

	// Cannot delete yourself
	if funk.Contains(r.UserIds, ctxUser.ID) {
		return nil, tools.NewValidatorI18nError("user.cannot_delete_self", nil)
	}

	// Cannot delete users with a lower role sort value (higher role level) than yourself
	for _, sort := range roleMinSortList {
		if int(minSort) > sort {
			return nil, tools.NewValidatorI18nError("user.cannot_delete_higher", nil)
		}
	}

	users, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.info_failed", i18n.Args{"error": err.Error()})
	}

	// First delete the user from LDAP
	for _, user := range users {
		err := ildap.User.Delete(user.UserDN)
		if err != nil {
			return nil, tools.NewLdapI18nError("user.ldap_delete_failed", i18n.Args{"error": err.Error()})
		}
	}

	// Then delete the user from MySQL
	err = isql.User.Delete(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.mysql_delete_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// ChangePwd changes the password
func (l UserLogic) ChangePwd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserChangePwdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// The password from the frontend is RSA-encrypted; decrypt it first
	// Decrypt the password with RSA
	decodeOldPassword, err := tools.RSADecrypt([]byte(r.OldPassword), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, tools.NewValidatorI18nError("legacy.user.old_password_parse_failed", nil)
	}
	decodeNewPassword, err := tools.RSADecrypt([]byte(r.NewPassword), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, tools.NewValidatorI18nError("legacy.user.new_password_parse_failed", nil)
	}
	r.OldPassword = string(decodeOldPassword)
	r.NewPassword = string(decodeNewPassword)
	// Get the current user
	user, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("legacy.common.current_user_failed", nil)
	}
	// Get the user's actual password
	// correctPasswd := user.Password
	// Check whether the password in the frontend request matches the actual password
	// err = tools.ComparePasswd(correctPasswd, r.OldPassword)
	if tools.NewParPasswd(user.Password) != r.OldPassword {
		return nil, tools.NewValidatorI18nError("legacy.user.old_password_incorrect", nil)
	}
	// To update an LDAP password, specify the user DN and new password
	err = ildap.User.ChangePwd(user.UserDN, "", r.NewPassword)
	if err != nil {
		return nil, tools.NewLdapI18nError("user.ldap_password_update_failed", i18n.Args{"error": err.Error()})
	}

	// Update the password
	err = isql.User.ChangePwd(user.Username, tools.NewGenPasswd(r.NewPassword))
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.mysql_password_update_failed", i18n.Args{"error": err.Error()})
	}

	return nil, nil
}

// ResetPassword resets the user password
func (l UserLogic) ResetPassword(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserResetPasswordReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// Check whether the user exists
	if !isql.User.Exist(tools.H{"username": r.Username}) {
		return nil, tools.NewValidatorI18nError("user.not_found", nil)
	}

	// Get the user information
	user := new(model.User)
	err := isql.User.Find(tools.H{"username": r.Username}, user)
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.info_failed", i18n.Args{"error": err.Error()})
	}

	// Generate a random password
	newPassword := tools.GenerateRandomPassword()

	// Update the password in LDAP
	err = ildap.User.ChangePwd(user.UserDN, "", newPassword)
	if err != nil {
		return nil, tools.NewLdapI18nError("user.ldap_password_update_failed", i18n.Args{"error": err.Error()})
	}

	// Update the password in MySQL
	err = isql.User.ChangePwd(user.Username, tools.NewGenPasswd(newPassword))
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.mysql_password_update_failed", i18n.Args{"error": err.Error()})
	}

	// Send the password reset notification email
	if err := tools.SendPasswordResetNotificationI18n(user.Username, user.Nickname, user.Mail, newPassword, i18n.LocaleFromContext(c)); err != nil {
		common.Log.Warnf("Failed to send password reset notification email, user: %s, email: %s, error: %v", user.Username, user.Mail, err)
	}

	return map[string]string{"newPassword": newPassword}, nil
}

// ChangeUserStatus changes the user status
func (l UserLogic) ChangeUserStatus(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserChangeUserStatusReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// Validation
	filter := tools.H{"id": r.ID}
	if !isql.User.Exist(filter) {
		return nil, tools.NewValidatorI18nError("legacy.user.not_found", nil)
	}
	user := new(model.User)
	err := isql.User.Find(filter, user)
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.mysql_query_failed", i18n.Args{"error": err.Error()})
	}

	if r.Status == 1 && r.Status == user.Status {
		return nil, tools.NewValidatorI18nError("legacy.user.already_active", nil)
	}
	if r.Status == 2 && r.Status == user.Status {
		return nil, tools.NewValidatorI18nError("legacy.user.already_inactive", nil)
	}

	// Get the current user; only admins can change user status
	// Get the current user and their minimum role sort value (highest role level)
	minSort, _, err := isql.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		return nil, tools.NewValidatorI18nError("user.current_min_role_sort_failed", nil)
	}

	if int(minSort) != 1 {
		return nil, tools.NewValidatorI18nError("legacy.user.only_admin_change_status", nil)
	}

	if r.Status == 2 {
		err = ildap.User.Delete(user.UserDN)
		if err != nil {
			return nil, tools.NewLdapI18nError("user.ldap_delete_failed", i18n.Args{"error": err.Error()})
		}
	} else {
		err = ildap.User.Add(user)
		if err != nil {
			return nil, tools.NewLdapI18nError("user.ldap_add_failed", i18n.Args{"error": err.Error()})
		}
	}
	err = isql.User.ChangeStatus(int(r.ID), int(r.Status))
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.mysql_status_update_failed", i18n.Args{"error": err.Error()})
	}
	return nil, nil
}

// GetUserInfo gets user information
func (l UserLogic) GetUserInfo(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.UserGetUserInfoReq)
	if !ok {
		return nil, ReqAssertErr
	}

	_ = c
	_ = r

	user, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlI18nError("user.current_info_failed", i18n.Args{"error": err.Error()})
	}
	return user, nil
}
