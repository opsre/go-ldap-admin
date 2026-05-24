package response

import "github.com/eryajf/go-ldap-admin/model"

type ApiTreeRsp struct {
	ID              int          `json:"ID"`
	Remark          string       `json:"remark"`
	RemarkDisplay   string       `json:"remarkDisplay,omitempty"`
	Category        string       `json:"category"`
	CategoryDisplay string       `json:"categoryDisplay,omitempty"`
	Children        []*model.Api `json:"children"`
}

type ApiListRsp struct {
	Total int64       `json:"total"`
	Apis  []model.Api `json:"apis"`
}
