package request

// OperationLogListReq is the operation log request struct
type OperationLogListReq struct {
	Username string `json:"username" form:"username"`
	Ip       string `json:"ip" form:"ip"`
	Path     string `json:"path" form:"path"`
	Method   string `json:"method" form:"method"`
	Status   int    `json:"status" form:"status"`
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}

// OperationLogDeleteReq is the bulk operation log deletion request struct
type OperationLogDeleteReq struct {
	OperationLogIds []uint `json:"operationLogIds" validate:"required"`
}
