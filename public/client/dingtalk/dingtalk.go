package dingtalk

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/zhaoyunxing92/dingtalk/v2/request"
)

// 官方文档地址： https://open.dingtalk.com/document/orgapp-server/obtain-the-department-list
// GetAllDepts 获取所有部门
func GetAllDepts() (ret []map[string]any, err error) {
	depts, err := InitDingTalkClient().FetchDeptList(1, true, "zh_CN")
	if err != nil {
		return ret, err
	}
	if len(config.Conf.DingTalk.DeptList) == 0 {

		ret = make([]map[string]any, 0)
		for _, dept := range depts.Dept {
			ele := make(map[string]any)
			ele["id"] = dept.Id
			ele["name"] = dept.Name
			ele["parentid"] = dept.ParentId
			ele["custom_name_pinyin"] = tools.ConvertToPinYin(dept.Name)
			ret = append(ret, ele)
		}
	} else {

		// 遍历配置的部门ID列表获取数据进行处理
		// 从取得的所有部门列表中将配置的部门ID筛选出来再去请求其子部门过滤为1和为配置值的部门ID
		ret = make([]map[string]any, 0)

		for _, dept := range depts.Dept {
			inset := false
			for _, dep_s := range config.Conf.DingTalk.DeptList {
				if strings.HasPrefix(dep_s, "^") {
					continue
				}
				setdepid, _ := strconv.Atoi(dep_s)
				if dept.Id == setdepid {
					inset = true
					break
				}
			}
			if dept.Id == 1 || inset {
				ele := make(map[string]any)
				ele["id"] = dept.Id
				ele["name"] = dept.Name
				ele["parentid"] = dept.ParentId
				ele["custom_name_pinyin"] = tools.ConvertToPinYin(dept.Name)
				ret = append(ret, ele)
			}
		}

		for _, dep_s := range config.Conf.DingTalk.DeptList {
			dept_id := dep_s

			if strings.HasPrefix(dep_s, "^") || dept_id == "1" {
				continue
			}
			depid, _ := strconv.Atoi(dept_id)
			depts, err := InitDingTalkClient().FetchDeptList(depid, true, "zh_CN")

			if err != nil {
				return ret, err
			}

			for _, dept := range depts.Dept {
				ele := make(map[string]any)
				ele["id"] = dept.Id
				ele["name"] = dept.Name
				ele["parentid"] = dept.ParentId
				ele["custom_name_pinyin"] = tools.ConvertToPinYin(dept.Name)
				ret = append(ret, ele)
			}
		}
	}
	return
}

// 官方文档地址： https://open.dingtalk.com/document/orgapp-server/queries-the-complete-information-of-a-department-user
// GetAllUsers 获取所有员工信息
func GetAllUsers() (ret []map[string]any, err error) {
	depts, err := GetAllDepts()
	if err != nil {
		return nil, err
	}
	for _, dept := range depts {
		r := request.DeptDetailUserInfo{
			DeptId:   dept["id"].(int),
			Cursor:   0,
			Size:     99,
			Language: "zh_CN",
		}
		for {
			//获取钉钉部门人员信息
			rsp, err := InitDingTalkClient().GetDeptDetailUserInfo(&r)
			if err != nil {
				return nil, err
			}
			for _, user := range rsp.Page.List {
				ele := make(map[string]any)
				ele["userid"] = user.UserId
				ele["unionid"] = user.UnionId
				ele["custom_name_pinyin"] = tools.ConvertToPinYin(user.Name)
				ele["name"] = user.Name
				ele["avatar"] = user.Avatar
				ele["mobile"] = user.Mobile
				ele["job_number"] = user.JobNumber
				ele["title"] = user.Title
				ele["work_place"] = user.WorkPlace
				ele["remark"] = user.Remark
				ele["leader"] = user.Leader
				ele["org_email"] = user.OrgEmail
				if user.OrgEmail != "" {
					ele["custom_nickname_org_email"] = strings.Split(user.OrgEmail, "@")[0]
				}
				ele["email"] = user.Email
				if user.Email != "" {
					ele["custom_nickname_email"] = strings.Split(user.Email, "@")[0]
				}
				// 部门ids
				var sourceDeptIds []string
				for _, deptId := range user.DeptIds {
					sourceDeptIds = append(sourceDeptIds, fmt.Sprintf("%s_%d", config.Conf.DingTalk.Flag, deptId))
				}
				ele["department_ids"] = sourceDeptIds
				ret = append(ret, ele)
			}
			if !rsp.Page.HasMore {
				break
			}
			r.Cursor = rsp.Page.NextCursor
		}
	}
	return
}

// 官方文档：https://open.dingtalk.com/document/orgapp-server/intelligent-personnel-query-company-turnover-list
// GetLeaveUserIds 获取离职人员ID列表
func GetLeaveUserIds() ([]string, error) {
	var ids []string
	ReqParm := struct {
		Cursor int `json:"cursor"`
		Size   int `json:"size"`
	}{
		Cursor: 0,
		Size:   50,
	}

	for {
		rsp, err := InitDingTalkClient().GetHrmResignEmployeeIds(ReqParm.Cursor, ReqParm.Size)
		if err != nil {
			return nil, err
		}
		ids = append(ids, rsp.Result.UserIds...)
		if rsp.Result.NextCursor == 0 {
			break
		}
		ReqParm.Cursor = rsp.Result.NextCursor
	}
	return ids, nil
}

// 官方文档：https://open.dingtalk.com/document/orgapp/query-the-details-of-employees-who-have-left-office
// GetLeaveUserIdsByDateRange 新接口根据时间范围获取离职人员ID列表
// GetHrmempLeaveRecordsKey    = "/v1.0/contact/empLeaveRecords"
func GetLeaveUserIdsDateRange(pushDays uint) ([]string, error) {
	var ids []string
	// 配置值为正数,往前推转为负数
	var leaveDays = int(0 - pushDays)
	startTime := time.Now().AddDate(0, 0, leaveDays).Format("2006-01-02T15:04:05Z")
	endTime := time.Now().Format("2006-01-02T15:04:05Z")
	ReqParm := struct {
		StartTime  string `json:"startTime"`
		EndTime    string `json:"endTime"`
		NextToken  string `json:"nextToken"`
		MaxResults int    `json:"maxResults"`
	}{
		StartTime:  startTime,
		EndTime:    endTime,
		NextToken:  "0",
		MaxResults: 50,
	}
	// 使用新的使用时间范围查询离职人员接口获取离职用户ID
	for {
		rsp, err := InitDingTalkClient().GetHrmEmpLeaveRecords(ReqParm.StartTime, ReqParm.EndTime, ReqParm.NextToken, ReqParm.MaxResults)
		if err != nil {
			return nil, err
		}
		for _, g := range rsp.Records {
			ids = append(ids, g.UserId)
		}

		if rsp.NextToken == "0" || rsp.NextToken == "" {
			break
		}
		ReqParm.NextToken = rsp.NextToken
	}
	return ids, nil
}
