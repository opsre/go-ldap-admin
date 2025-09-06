package feishu

import (
	"context"
	"fmt"
	"strings"

	"github.com/chyroc/lark"
	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/public/tools"
)

// 官方文档： https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/department/children
// GetAllDepts 获取所有部门
func GetAllDepts() (ret []map[string]any, err error) {
	var (
		fetchChild       = true
		pageSize   int64 = 50
		pageToken        = ""
		// DeptID     lark.DepartmentIDType = "department_id"
	)

	if len(config.Conf.FeiShu.DeptList) == 0 {
		req := lark.GetDepartmentListReq{
			// DepartmentIDType: &DeptID,
			PageToken:    &pageToken,
			FetchChild:   &fetchChild,
			PageSize:     &pageSize,
			DepartmentID: "0",
		}
		for {
			res, _, err := InitFeiShuClient().Contact.GetDepartmentList(context.TODO(), &req)
			if err != nil {
				return nil, err
			}

			for _, dept := range res.Items {
				ele := make(map[string]any)
				ele["name"] = dept.Name
				ele["custom_name_pinyin"] = tools.ConvertToPinYin(dept.Name)
				ele["parent_department_id"] = dept.ParentDepartmentID
				ele["department_id"] = dept.DepartmentID
				ele["open_department_id"] = dept.OpenDepartmentID
				ele["leader_user_id"] = dept.LeaderUserID
				ele["unit_ids"] = dept.UnitIDs
				ret = append(ret, ele)
			}
			if !res.HasMore {
				break
			}
			pageToken = res.PageToken
		}
	} else {
		//使用dept-list来一个一个添加部门，开头为^的不添加子部门
		isInDeptList := func(id string) bool {
			for _, v := range config.Conf.FeiShu.DeptList {
				v = strings.TrimPrefix(v, "^")
				if id == v {
					return true
				}
			}
			return false
		}
		dep_append_norepeat := func(ret []map[string]any, dept map[string]any) []map[string]any {
			for _, v := range ret {
				if v["open_department_id"] == dept["open_department_id"] {
					return ret
				}
			}
			return append(ret, dept)
		}
		for _, dep_s := range config.Conf.FeiShu.DeptList {
			dept_id := dep_s
			no_add_children := false
			if strings.HasPrefix(dep_s, "^") {
				no_add_children = true
				dept_id = dep_s[1:]
			}
			req := lark.GetDepartmentReq{
				DepartmentID: dept_id,
			}
			res, _, err := InitFeiShuClient().Contact.GetDepartment(context.TODO(), &req)
			if err != nil {
				return nil, err
			}
			ele := make(map[string]any)

			ele["name"] = res.Department.Name
			ele["custom_name_pinyin"] = tools.ConvertToPinYin(res.Department.Name)
			if isInDeptList(res.Department.ParentDepartmentID) {
				ele["parent_department_id"] = res.Department.ParentDepartmentID
			} else {
				ele["parent_department_id"] = "0"
			}
			ele["department_id"] = res.Department.DepartmentID
			ele["open_department_id"] = res.Department.OpenDepartmentID
			ele["leader_user_id"] = res.Department.LeaderUserID
			ele["unit_ids"] = res.Department.UnitIDs
			ret = dep_append_norepeat(ret, ele)

			if !no_add_children {
				pageToken = ""
				req := lark.GetDepartmentListReq{
					// DepartmentIDType: &DeptID,
					PageToken:    &pageToken,
					FetchChild:   &fetchChild,
					PageSize:     &pageSize,
					DepartmentID: dept_id,
				}
				for {
					res, _, err := InitFeiShuClient().Contact.GetDepartmentList(context.TODO(), &req)
					if err != nil {
						return nil, err
					}

					for _, dept := range res.Items {
						ele := make(map[string]any)
						ele["name"] = dept.Name
						ele["custom_name_pinyin"] = tools.ConvertToPinYin(dept.Name)
						ele["parent_department_id"] = dept.ParentDepartmentID
						ele["department_id"] = dept.DepartmentID
						ele["open_department_id"] = dept.OpenDepartmentID
						ele["leader_user_id"] = dept.LeaderUserID
						ele["unit_ids"] = dept.UnitIDs
						ret = dep_append_norepeat(ret, ele)
					}
					if !res.HasMore {
						break
					}
					pageToken = res.PageToken
				}
			}
		}
	}
	return
}

// 官方文档： https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/user/find_by_department
// GetAllUsers 获取所有员工信息
func GetAllUsers() (ret []map[string]any, err error) {
	var (
		pageSize  int64 = 50
		pageToken       = ""
		// deptidtype lark.DepartmentIDType = "department_id"
	)
	depts, err := GetAllDepts()
	if err != nil {
		return nil, err
	}

	deptids := make([]string, 0)
	// deptids = append(deptids, "0")
	for _, dept := range depts {
		deptids = append(deptids, dept["open_department_id"].(string))
	}

	for _, deptid := range deptids {
		req := lark.GetUserListReq{
			PageSize:  &pageSize,
			PageToken: &pageToken,
			// DepartmentIDType: &deptidtype,
			DepartmentID: deptid,
		}
		for {
			res, _, err := InitFeiShuClient().Contact.GetUserList(context.Background(), &req)
			if err != nil {
				return nil, err
			}
			for _, user := range res.Items {
				ele := make(map[string]any)
				ele["name"] = user.Name
				ele["custom_name_pinyin"] = tools.ConvertToPinYin(user.Name)
				ele["union_id"] = user.UnionID
				ele["user_id"] = user.UserID
				ele["open_id"] = user.OpenID
				ele["en_name"] = user.EnName
				ele["nickname"] = user.Nickname
				if user.Email != "" {
					ele["custom_nickname_email"] = strings.Split(user.Email, "@")[0]
				}
				if user.EnterpriseEmail != "" {
					ele["custom_nickname_enterprise_email"] = strings.Split(user.EnterpriseEmail, "@")[0]
				}
				ele["email"] = user.Email
				ele["mobile"] = user.Mobile
				ele["gender"] = user.Gender
				ele["avatar"] = user.Avatar.AvatarOrigin
				ele["city"] = user.City
				ele["country"] = user.Country
				ele["work_station"] = user.WorkStation
				ele["join_time"] = user.JoinTime
				ele["employee_no"] = user.EmployeeNo
				ele["enterprise_email"] = user.EnterpriseEmail
				ele["job_title"] = user.JobTitle
				// 部门ids
				var sourceDeptIds []string
				for _, deptId := range user.DepartmentIDs {
					sourceDeptIds = append(sourceDeptIds, fmt.Sprintf("%s_%s", config.Conf.FeiShu.Flag, deptId))
				}
				ele["department_ids"] = sourceDeptIds
				ret = append(ret, ele)
			}
			if !res.HasMore {
				pageToken = ""
				break
			}
			pageToken = res.PageToken
		}
	}
	return
}

// 官方文档： https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/ehr/ehr-v1/employee/list
// GetLeaveUserIds 获取离职人员ID列表
func GetLeaveUserIds() ([]string, error) {
	var ids []string
	users, _, err := InitFeiShuClient().EHR.GetEHREmployeeList(context.TODO(), &lark.GetEHREmployeeListReq{
		Status:     []int64{5},
		UserIDType: lark.IDTypePtr(lark.IDTypeUnionID), // 只查询unionID
	})
	if err != nil {
		return nil, err
	}
	for _, user := range users.Items {
		ids = append(ids, user.UserID)
	}
	return ids, nil
}
