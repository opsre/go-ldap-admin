package request

// ApiListReq is the resource list request struct
type ApiListReq struct {
	Method   string `json:"method" form:"method"`
	Path     string `json:"path" form:"path"`
	Category string `json:"category" form:"category"`
	Creator  string `json:"creator" form:"creator"`
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}

// ApiAddReq is the resource creation request struct
type ApiAddReq struct {
	Method   string `json:"method" validate:"required,min=1,max=20"`
	Path     string `json:"path" validate:"required,min=1,max=100"`
	Category string `json:"category" validate:"required,min=1,max=50"`
	Remark   string `json:"remark" validate:"min=0,max=100"`
}

// ApiUpdateReq is the resource update request struct
type ApiUpdateReq struct {
	ID       uint   `json:"id" validate:"required"`
	Method   string `json:"method" validate:"min=1,max=20"`
	Path     string `json:"path" validate:"min=1,max=100"`
	Category string `json:"category" validate:"min=1,max=50"`
	Remark   string `json:"remark" validate:"min=0,max=100"`
}

// ApiDeleteReq is the resource deletion request struct
type ApiDeleteReq struct {
	ApiIds []uint `json:"apiIds" validate:"required"`
}

// ApiGetTreeReq is the resource tree request struct
type ApiGetTreeReq struct {
}
