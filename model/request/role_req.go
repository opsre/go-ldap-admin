package request

// RoleAddReq is the resource creation request struct
type RoleAddReq struct {
	Name    string `json:"name" validate:"required,min=1,max=20"`
	Keyword string `json:"keyword" validate:"required,min=1,max=20"`
	Remark  string `json:"remark" validate:"min=0,max=100"`
	Status  uint   `json:"status" validate:"oneof=1 2"`
	Sort    uint   `json:"sort" validate:"gte=1,lte=999"`
}

// RoleListReq is the list request struct
type RoleListReq struct {
	Name     string `json:"name" form:"name"`
	Keyword  string `json:"keyword" form:"keyword"`
	Status   uint   `json:"status" form:"status"`
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}

// RoleUpdateReq is the resource update request struct
type RoleUpdateReq struct {
	ID      uint   `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required,min=1,max=20"`
	Keyword string `json:"keyword" validate:"required,min=1,max=20"`
	Remark  string `json:"remark" validate:"min=0,max=100"`
	Status  uint   `json:"status" validate:"oneof=1 2"`
	Sort    uint   `json:"sort" validate:"gte=1,lte=999"`
}

// RoleDeleteReq is the resource deletion request struct
type RoleDeleteReq struct {
	RoleIds []uint `json:"roleIds" validate:"required"`
}

// RoleGetTreeReq is the resource tree request struct
type RoleGetTreeReq struct {
}

// RoleGetMenuListReq is the role menu list request struct
type RoleGetMenuListReq struct {
	RoleID uint `json:"roleId" form:"roleId" validate:"required"`
}

// RoleGetApiListReq is the role API list request struct
type RoleGetApiListReq struct {
	RoleID uint `json:"roleId" form:"roleId" validate:"required"`
}

// RoleUpdateMenusReq is the role menu update request struct
type RoleUpdateMenusReq struct {
	RoleID  uint   `json:"roleId" validate:"required"`
	MenuIds []uint `json:"menuIds" validate:"required"`
}

// RoleUpdateApisReq is the role API update request struct
type RoleUpdateApisReq struct {
	RoleID uint   `json:"roleId" validate:"required"`
	ApiIds []uint `json:"apiIds" validate:"required"`
}
