package request

// GroupListReq is the resource list request struct
type GroupListReq struct {
	GroupName string `json:"groupName" form:"groupName"`
	Remark    string `json:"remark" form:"remark"`
	PageNum   int    `json:"pageNum" form:"pageNum"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
	SyncState uint   `json:"syncState" form:"syncState"`
}

// GroupListAllReq is the unpaginated resource list request struct
type GroupListAllReq struct {
	GroupName          string `json:"groupName" form:"groupName"`
	GroupType          string `json:"groupType" form:"groupType"`
	Remark             string `json:"remark" form:"remark"`
	Source             string `json:"source" form:"source"`
	SourceDeptId       string `json:"sourceDeptId"`
	SourceDeptParentId string `json:"SourceDeptParentId"`
}

// GroupAddReq is the resource creation request struct
type GroupAddReq struct {
	GroupType string `json:"groupType" validate:"required,min=1,max=20"`
	GroupName string `json:"groupName" validate:"required,min=1,max=128"`
	//Parent ID is required and must be greater than or equal to 0
	ParentId uint   `json:"parentId" validate:"omitempty,min=0"`
	Remark   string `json:"remark" validate:"min=0,max=128"` // Chinese description of the group
}

// DingTalkGroupAddReq is the DingTalk resource creation request struct
type DingGroupAddReq struct {
	GroupType string `json:"groupType" validate:"required,min=1,max=20"`
	GroupName string `json:"groupName" validate:"required,min=1,max=128"`
	//Parent ID is required and must be greater than or equal to 0
	ParentId           uint   `json:"parentId" validate:"omitempty,min=0"`
	Remark             string `json:"remark" validate:"min=0,max=128"` // Chinese description of the group
	SourceDeptId       string `json:"sourceDeptId"`
	Source             string `json:"source"`
	SourceDeptParentId string `json:"SourceDeptParentId"`
	SourceUserNum      int    `json:"sourceUserNum"`
}

// WeComGroupAddReq is the WeCom resource creation request struct
type WeComGroupAddReq struct {
	GroupType string `json:"groupType" validate:"required,min=1,max=20"`
	GroupName string `json:"groupName" validate:"required,min=1,max=128"`
	//Parent ID is required and must be greater than or equal to 0
	ParentId           uint   `json:"parentId" validate:"omitempty,min=0"`
	Remark             string `json:"remark" validate:"min=0,max=128"` // Chinese description of the group
	SourceDeptId       string `json:"sourceDeptId"`
	Source             string `json:"source"`
	SourceDeptParentId string `json:"SourceDeptParentId"`
	SourceUserNum      int    `json:"sourceUserNum"`
}

// GroupUpdateReq is the resource update request struct
type GroupUpdateReq struct {
	ID        uint   `json:"id" form:"id" validate:"required"`
	GroupName string `json:"groupName" validate:"required,min=1,max=128"`
	Remark    string `json:"remark" validate:"min=0,max=128"` // Chinese description of the group
}

// GroupDeleteReq is the resource deletion request struct
type GroupDeleteReq struct {
	GroupIds []uint `json:"groupIds" validate:"required"`
}

// GroupGetTreeReq is the resource tree request struct
type GroupGetTreeReq struct {
	GroupName string `json:"groupName" form:"groupName"`
	Remark    string `json:"remark" form:"remark"`
	PageNum   int    `json:"pageNum" form:"pageNum"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
}

type GroupAddUserReq struct {
	GroupID uint   `json:"groupId" validate:"required"`
	UserIds []uint `json:"userIds" validate:"required"`
}

type GroupRemoveUserReq struct {
	GroupID uint   `json:"groupId" validate:"required"`
	UserIds []uint `json:"userIds" validate:"required"`
}

// UserInGroupReq lists users in a group
type UserInGroupReq struct {
	GroupID  uint   `json:"groupId" form:"groupId" validate:"required"`
	Nickname string `json:"nickname" form:"nickname"`
}

// UserNoInGroupReq lists users outside a group
type UserNoInGroupReq struct {
	GroupID  uint   `json:"groupId" form:"groupId" validate:"required"`
	Nickname string `json:"nickname" form:"nickname"`
}

// SyncDingTalkDeptsReq syncs DingTalk department information
type SyncDingTalkDeptsReq struct {
}

// SyncWeComDeptsReq syncs WeCom department information
type SyncWeComDeptsReq struct {
}

// SyncFeiShuDeptsReq syncs Feishu department information
type SyncFeiShuDeptsReq struct {
}

// SyncOpenLdapDeptsReq syncs existing LDAP department information
type SyncOpenLdapDeptsReq struct {
}

// SyncOpenLdapDeptsReq syncs existing LDAP department information
type SyncSqlGrooupsReq struct {
	GroupIds []uint `json:"groupIds" validate:"required"`
}
