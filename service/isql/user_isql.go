package isql

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type UserService struct{}

// Cache current user information to avoid frequent database reads
var userInfoCache = cache.New(24*time.Hour, 48*time.Hour)

// Add a resource
func (s UserService) Add(user *model.User) error {
	user.Password = tools.NewGenPasswd(user.Password)
	//result := common.DB.Create(user)
	//return user.ID, result.Error
	return common.DB.Create(user).Error
}

// List gets the data list
func (s UserService) List(req *request.UserListReq) ([]*model.User, error) {
	var list []*model.User
	db := common.DB.Model(&model.User{}).Order("id DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	departmentId := req.DepartmentId
	if len(departmentId) > 0 {
		db = db.Where("department_id = ?", departmentId)
	}
	givenName := strings.TrimSpace(req.GivenName)
	if givenName != "" {
		db = db.Where("given_name LIKE ?", fmt.Sprintf("%%%s%%", givenName))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	syncState := req.SyncState
	if syncState != 0 {
		db = db.Where("sync_state = ?", syncState)
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Preload("Roles").Find(&list).Debug().Error
	return list, err
}

// ListCout gets the number of records matching the conditions
func (s UserService) ListCount(req *request.UserListReq) (int64, error) {
	var count int64
	db := common.DB.Model(&model.User{}).Order("id DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	departmentId := req.DepartmentId
	if len(departmentId) > 0 {
		db = db.Where("department_id = ?", departmentId)
	}
	givenName := strings.TrimSpace(req.GivenName)
	if givenName != "" {
		db = db.Where("given_name LIKE ?", fmt.Sprintf("%%%s%%", givenName))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	syncState := req.SyncState
	if syncState != 0 {
		db = db.Where("sync_state = ?", syncState)
	}

	err := db.Count(&count).Error
	return count, err
}

// List gets the data list
func (s UserService) ListAll() (list []*model.User, err error) {
	err = common.DB.Model(&model.User{}).Order("created_at DESC").Find(&list).Error

	return list, err
}

// Count gets the total number of records
func (s UserService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.User{}).Count(&count).Error
	return count, err
}

// Exist checks whether the resource exists
func (s UserService) Exist(filter map[string]any) bool {
	var dataObj model.User
	err := common.DB.Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Find gets the stored user with the highest sequence number among users with the same name
func (s UserService) FindTheSameUserName(username string, data *model.User) error {
	return common.DB.Where("username REGEXP ? ", fmt.Sprintf("^%s[0-9]{0,3}$", username)).Order("username desc").First(&data).Error
}

// Find a single resource
func (s UserService) Find(filter map[string]any, data *model.User) error {
	return common.DB.Where(filter).Preload("Roles").First(&data).Error
}

// Update a resource
func (s UserService) Update(user *model.User) error {
	err := common.DB.Model(user).Updates(user).Error
	if err != nil {
		return err
	}
	err = common.DB.Model(user).Association("Roles").Replace(user.Roles)

	// Update the user information cache if the update succeeds
	if err == nil {
		userDb := &model.User{}
		common.DB.Where("username = ?", user.Username).Preload("Roles").First(&userDb)
		userInfoCache.Set(user.Username, *userDb, cache.DefaultExpiration)
	}
	return err
}

// GetUserMinRoleSortsByIds gets the minimum role sort value by user IDs
func (s UserService) GetUserMinRoleSortsByIds(ids []uint) ([]int, error) {
	// Get user information by user ID
	var userList []model.User
	err := common.DB.Where("id IN (?)", ids).Preload("Roles").Find(&userList).Error
	if err != nil {
		return []int{}, err
	}
	if len(userList) == 0 {
		return []int{}, errors.New("no user information found")
	}
	var roleMinSortList []int
	for _, user := range userList {
		roles := user.Roles
		var roleSortList []int
		for _, role := range roles {
			roleSortList = append(roleSortList, int(role.Sort))
		}
		roleMinSort := funk.MinInt(roleSortList).(int)
		roleMinSortList = append(roleMinSortList, roleMinSort)
	}
	return roleMinSortList, nil
}

// GetCurrentUserMinRoleSort gets the current user and their minimum role sort value (highest role level)
func (s UserService) GetCurrentUserMinRoleSort(c *gin.Context) (uint, model.User, error) {
	// Get the current user
	ctxUser, err := s.GetCurrentLoginUser(c)
	if err != nil {
		return 999, ctxUser, err
	}
	// Get all roles of the current user
	currentRoles := ctxUser.Roles
	// Get the current user role sort values and compare them with those sent by the frontend
	var currentRoleSorts []int
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}
	// The minimum role sort value (highest role level) of the current user
	currentRoleSortMin := uint(funk.MinInt(currentRoleSorts).(int))

	return currentRoleSortMin, ctxUser, nil
}

// Delete multiple resources
func (s UserService) Delete(ids []uint) error {
	// Users and roles have a many-to-many relationship
	var users []model.User
	for _, id := range ids {
		// Get the user by ID
		filter := tools.H{"id": id}

		user := new(model.User)
		err := s.Find(filter, user)
		if err != nil {
			return fmt.Errorf("failed to get user information: %v", err)
		}
		users = append(users, *user)
	}

	err := common.DB.Debug().Select("Roles").Unscoped().Delete(&users).Error
	if err != nil {
		return err
	} else {
		// Delete the user information cache after successfully deleting the user
		for _, user := range users {
			userInfoCache.Delete(user.Username)
		}
	}

	// Delete the user-group associations
	err = common.DB.Debug().Exec("DELETE FROM group_users WHERE user_id IN (?)", ids).Error
	if err != nil {
		return err
	}

	return err
}

// GetUserByIds gets the minimum role sort value by user IDs
func (s UserService) GetUserByIds(ids []uint) ([]model.User, error) {
	// Get user information by user ID
	var userList []model.User
	err := common.DB.Where("id IN (?)", ids).Preload("Roles").Find(&userList).Error
	return userList, err
}

// ChangePwd updates the password
func (s UserService) ChangePwd(username string, hashNewPasswd string) error {
	err := common.DB.Model(&model.User{}).Where("username = ?", username).Update("password", hashNewPasswd).Error
	// Update the current user information cache if the password update succeeds
	// First get the cache
	cacheUser, found := userInfoCache.Get(username)
	if err == nil {
		if found {
			user := cacheUser.(model.User)
			user.Password = hashNewPasswd
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		} else {
			// If there is no cache, get the user information and cache it
			var user model.User
			common.DB.Where("username = ?", username).Preload("Roles").First(&user)
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		}
	}

	return err
}

// ChangeStatus updates the status
func (s UserService) ChangeStatus(id, status int) error {
	return common.DB.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// ChangeSyncState updates the user sync status
func (s UserService) ChangeSyncState(id, status int) error {
	return common.DB.Model(&model.User{}).Where("id = ?", id).Update("sync_state", status).Error
}

// GetCurrentLoginUser gets the current logged-in user information
// Cache this to reduce database access
func (s UserService) GetCurrentLoginUser(c *gin.Context) (model.User, error) {
	var newUser model.User
	ctxUser, exist := c.Get("user")
	if !exist {
		return newUser, errors.New("用户未登录")
	}
	u, _ := ctxUser.(model.User)

	// First get the cache
	cacheUser, found := userInfoCache.Get(u.Username)
	var user model.User
	var err error
	if found {
		user = cacheUser.(model.User)
		err = nil
	} else {
		// If it is not cached, get it from the database
		user, err = s.GetUserById(u.ID)
		// Cache it after a successful lookup
		if err != nil {
			userInfoCache.Delete(u.Username)
		} else {
			userInfoCache.Set(u.Username, user, cache.DefaultExpiration)
		}
	}
	return user, err
}

// Login
func (s UserService) Login(user *model.User) (*model.User, error) {
	// Get the user by username (normal status: the user is active)
	var firstUser model.User
	// err := common.DB.
	// 	Where("username = ?", user.Username).
	// 	Preload("Roles").
	// 	First(&firstUser).Error
	err := s.Find(tools.H{"username": user.Username}, &firstUser)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	// Check the user status
	userStatus := firstUser.Status
	if userStatus != 1 {
		return nil, errors.New("用户被禁用")
	}

	// Check the status of all user roles; login is denied if all roles are disabled
	// roles := firstUser.Roles
	// isValidate := false
	// for _, role := range roles {
	// 	// Login is allowed if at least one role is active
	// 	if role.Status == 1 {
	// 		isValidate = true
	// 		break
	// 	}
	// }

	if tools.NewParPasswd(firstUser.Password) != user.Password {
		return nil, errors.New("密码错误")
	}

	return &firstUser, nil
}

// ClearUserInfoCache clears all user information caches
func (s UserService) ClearUserInfoCache() {
	userInfoCache.Flush()
}

// GetUserById gets a single user
func (us UserService) GetUserById(id uint) (model.User, error) {
	var user model.User
	err := common.DB.Where("id = ?", id).Preload("Roles").First(&user).Error
	return user, err
}
