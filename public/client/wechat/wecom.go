package wechat

import (
	"fmt"
	"strings"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/wenerme/go-wecom/wecom"
)

// 官方文档： https://developer.work.weixin.qq.com/document/path/90208
// GetAllDepts 获取所有部门
func GetAllDepts() (ret []map[string]any, err error) {
	depts, err := InitWeComClient().ListDepartment(
		&wecom.ListDepartmentRequest{},
	)
	if err != nil {
		return nil, err
	}
	for _, dept := range depts.Department {
		ele := make(map[string]any)
		ele["name"] = dept.Name
		ele["custom_name_pinyin"] = tools.ConvertToPinYin(dept.Name)
		ele["id"] = dept.ID
		ele["name_en"] = dept.NameEn
		ele["parentid"] = dept.ParentID
		ret = append(ret, ele)
	}
	return ret, nil
}

// 官方文档： https://developer.work.weixin.qq.com/document/path/90201
// GetAllUsers 获取所有员工信息
func GetAllUsers() (ret []map[string]any, err error) {
	depts, err := GetAllDepts()
	if err != nil {
		return nil, err
	}
	for _, dept := range depts {
		users, err := InitWeComClient().ListUser(
			&wecom.ListUserRequest{
				DepartmentID: fmt.Sprintf("%d", dept["id"].(int)),
				FetchChild:   "1",
			},
		)
		if err != nil {
			return nil, err
		}
		for _, user := range users.UserList {
			ele := make(map[string]any)
			ele["name"] = user.Name
			ele["custom_name_pinyin"] = tools.ConvertToPinYin(user.Name)
			ele["userid"] = user.UserID
			ele["mobile"] = user.Mobile
			ele["position"] = user.Position
			ele["gender"] = user.Gender
			ele["email"] = user.Email
			if user.Email != "" {
				ele["custom_nickname_email"] = strings.Split(user.Email, "@")[0]
			}
			ele["biz_email"] = user.BizMail
			if user.BizMail != "" {
				ele["custom_nickname_biz_email"] = strings.Split(user.BizMail, "@")[0]
			}
			ele["avatar"] = user.Avatar
			ele["telephone"] = user.Telephone
			ele["alias"] = user.Alias
			ele["external_position"] = user.ExternalPosition
			ele["address"] = user.Address
			ele["open_userid"] = user.OpenUserID
			ele["main_department"] = user.MainDepartment
			ele["english_name"] = user.EnglishName
			// 部门ids
			var sourceDeptIds []string
			for _, deptId := range user.Department {
				sourceDeptIds = append(sourceDeptIds, fmt.Sprintf("%s_%d", config.Conf.WeCom.Flag, deptId))
			}
			ele["department_ids"] = sourceDeptIds
			ret = append(ret, ele)
		}
	}
	return ret, nil
}
