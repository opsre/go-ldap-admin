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

// 当前用户信息缓存，避免频繁获取数据库
var userInfoCache = cache.New(24*time.Hour, 48*time.Hour)

// Add 添加资源
func (s UserService) Add(user *model.User) error {
	user.Password = tools.NewGenPasswd(user.Password)
	//result := common.DB.Create(user)
	//return user.ID, result.Error
	return common.DB.Create(user).Error
}

// List 获取数据列表
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

// ListCout 获取符合条件的数据列表条数
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

// List 获取数据列表
func (s UserService) ListAll() (list []*model.User, err error) {
	err = common.DB.Model(&model.User{}).Order("created_at DESC").Find(&list).Error

	return list, err
}

// Count 获取数据总数
func (s UserService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.User{}).Count(&count).Error
	return count, err
}

// Exist 判断资源是否存在
func (s UserService) Exist(filter map[string]any) bool {
	var dataObj model.User
	err := common.DB.Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Find 获取同名用户已入库的序号最大的用户信息
func (s UserService) FindTheSameUserName(username string, data *model.User) error {
	return common.DB.Where("username REGEXP ? ", fmt.Sprintf("^%s[0-9]{0,3}$", username)).Order("username desc").First(&data).Error
}

// Find 获取单个资源
func (s UserService) Find(filter map[string]any, data *model.User) error {
	return common.DB.Where(filter).Preload("Roles").First(&data).Error
}

// Update 更新资源
func (s UserService) Update(user *model.User) error {
	err := common.DB.Model(user).Updates(user).Error
	if err != nil {
		return err
	}
	err = common.DB.Model(user).Association("Roles").Replace(user.Roles)

	// 如果更新成功就更新用户信息缓存
	if err == nil {
		userDb := &model.User{}
		common.DB.Where("username = ?", user.Username).Preload("Roles").First(&userDb)
		userInfoCache.Set(user.Username, *userDb, cache.DefaultExpiration)
	}
	return err
}

// GetUserMinRoleSortsByIds 根据用户ID获取用户角色排序最小值
func (s UserService) GetUserMinRoleSortsByIds(ids []uint) ([]int, error) {
	// 根据用户ID获取用户信息
	var userList []model.User
	err := common.DB.Where("id IN (?)", ids).Preload("Roles").Find(&userList).Error
	if err != nil {
		return []int{}, err
	}
	if len(userList) == 0 {
		return []int{}, errors.New("未获取到任何用户信息")
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

// GetCurrentUserMinRoleSort  获取当前用户角色排序最小值（最高等级角色）以及当前用户信息
func (s UserService) GetCurrentUserMinRoleSort(c *gin.Context) (uint, model.User, error) {
	// 获取当前用户
	ctxUser, err := s.GetCurrentLoginUser(c)
	if err != nil {
		return 999, ctxUser, err
	}
	// 获取当前用户的所有角色
	currentRoles := ctxUser.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := uint(funk.MinInt(currentRoleSorts).(int))

	return currentRoleSortMin, ctxUser, nil
}

// Delete 批量删除
func (s UserService) Delete(ids []uint) error {
	// 用户和角色存在多对多关联关系
	var users []model.User
	for _, id := range ids {
		// 根据ID获取用户
		filter := tools.H{"id": id}

		user := new(model.User)
		err := s.Find(filter, user)
		if err != nil {
			return fmt.Errorf("获取用户信息失败，err: %v", err)
		}
		users = append(users, *user)
	}

	err := common.DB.Debug().Select("Roles").Unscoped().Delete(&users).Error
	if err != nil {
		return err
	} else {
		// 删除用户成功，则删除用户信息缓存
		for _, user := range users {
			userInfoCache.Delete(user.Username)
		}
	}

	// 删除用户在group的关联
	err = common.DB.Debug().Exec("DELETE FROM group_users WHERE user_id IN (?)", ids).Error
	if err != nil {
		return err
	}

	return err
}

// GetUserByIds 根据用户ID获取用户角色排序最小值
func (s UserService) GetUserByIds(ids []uint) ([]model.User, error) {
	// 根据用户ID获取用户信息
	var userList []model.User
	err := common.DB.Where("id IN (?)", ids).Preload("Roles").Find(&userList).Error
	return userList, err
}

// ChangePwd 更新密码
func (s UserService) ChangePwd(username string, hashNewPasswd string) error {
	err := common.DB.Model(&model.User{}).Where("username = ?", username).Update("password", hashNewPasswd).Error
	// 如果更新密码成功，则更新当前用户信息缓存
	// 先获取缓存
	cacheUser, found := userInfoCache.Get(username)
	if err == nil {
		if found {
			user := cacheUser.(model.User)
			user.Password = hashNewPasswd
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		} else {
			// 没有缓存就获取用户信息缓存
			var user model.User
			common.DB.Where("username = ?", username).Preload("Roles").First(&user)
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		}
	}

	return err
}

// ChangeStatus 更新状态
func (s UserService) ChangeStatus(id, status int) error {
	return common.DB.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// ChangeSyncState 更新用户的同步状态
func (s UserService) ChangeSyncState(id, status int) error {
	return common.DB.Model(&model.User{}).Where("id = ?", id).Update("sync_state", status).Error
}

// GetCurrentLoginUser 获取当前登录用户信息
// 需要缓存，减少数据库访问
func (s UserService) GetCurrentLoginUser(c *gin.Context) (model.User, error) {
	var newUser model.User
	ctxUser, exist := c.Get("user")
	if !exist {
		return newUser, errors.New("用户未登录")
	}
	u, _ := ctxUser.(model.User)

	// 先获取缓存
	cacheUser, found := userInfoCache.Get(u.Username)
	var user model.User
	var err error
	if found {
		user = cacheUser.(model.User)
		err = nil
	} else {
		// 缓存中没有就获取数据库
		user, err = s.GetUserById(u.ID)
		// 获取成功就缓存
		if err != nil {
			userInfoCache.Delete(u.Username)
		} else {
			userInfoCache.Set(u.Username, user, cache.DefaultExpiration)
		}
	}
	return user, err
}

// Login 登录
func (s UserService) Login(user *model.User) (*model.User, error) {
	// 根据用户名获取用户(正常状态:用户状态正常)
	var firstUser model.User
	// err := common.DB.
	// 	Where("username = ?", user.Username).
	// 	Preload("Roles").
	// 	First(&firstUser).Error
	// if err != nil {
	// 	return nil, errors.New("用户不存在")
	// }
	err := s.Find(tools.H{"username": user.Username}, &firstUser)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	// 判断用户的状态
	userStatus := firstUser.Status
	if userStatus != 1 {
		return nil, errors.New("用户被禁用")
	}

	// 判断用户拥有的所有角色的状态,全部角色都被禁用则不能登录
	// roles := firstUser.Roles
	// isValidate := false
	// for _, role := range roles {
	// 	// 有一个正常状态的角色就可以登录
	// 	if role.Status == 1 {
	// 		isValidate = true
	// 		break
	// 	}
	// }

	// if !isValidate {
	// 	return nil, errors.New("用户角色被禁用")
	// }

	if tools.NewParPasswd(firstUser.Password) != user.Password {
		return nil, errors.New("密码错误")
	}

	// 校验密码
	// err = tools.ComparePasswd(firstUser.Password, user.Password)
	// if err != nil {
	// 	return &firstUser, errors.New("密码错误")
	// }
	return &firstUser, nil
}

// ClearUserInfoCache 清理所有用户信息缓存
func (s UserService) ClearUserInfoCache() {
	userInfoCache.Flush()
}

// GetUserById 获取单个用户
func (us UserService) GetUserById(id uint) (model.User, error) {
	var user model.User
	err := common.DB.Where("id = ?", id).Preload("Roles").First(&user).Error
	return user, err
}
