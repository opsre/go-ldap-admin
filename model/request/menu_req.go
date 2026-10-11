package request

// MenuAddReq is the resource creation request struct
type MenuAddReq struct {
	Name       string `json:"name" validate:"required,min=1,max=50"`
	Title      string `json:"title" validate:"required,min=1,max=50"`
	Icon       string `json:"icon" validate:"min=0,max=50"`
	Path       string `json:"path" validate:"required,min=1,max=100"`
	Redirect   string `json:"redirect" validate:"min=0,max=100"`
	Component  string `json:"component" validate:"required,min=1,max=100"`
	Sort       uint   `json:"sort" validate:"gte=1,lte=999"`
	Status     uint   `json:"status" validate:"oneof=1 2"`
	Hidden     uint   `json:"hidden" validate:"oneof=1 2"`
	NoCache    uint   `json:"noCache" validate:"oneof=1 2"`
	AlwaysShow uint   `json:"alwaysShow" validate:"oneof=1 2"`
	Breadcrumb uint   `json:"breadcrumb" validate:"oneof=1 2"`
	ActiveMenu string `json:"activeMenu" validate:"min=0,max=100"`
	ParentId   uint   `json:"parentId"`
}

// MenuListReq is the list request struct
type MenuListReq struct {
}

// MenuUpdateReq is the resource update request struct
type MenuUpdateReq struct {
	ID         uint   `json:"id" validate:"required"`
	Name       string `json:"name" validate:"required,min=1,max=50"`
	Title      string `json:"title" validate:"required,min=1,max=50"`
	Icon       string `json:"icon" validate:"min=0,max=50"`
	Path       string `json:"path" validate:"required,min=1,max=100"`
	Redirect   string `json:"redirect" validate:"min=0,max=100"`
	Component  string `json:"component" validate:"min=0,max=100"`
	Sort       uint   `json:"sort" validate:"gte=1,lte=999"`
	Status     uint   `json:"status" validate:"oneof=1 2"`
	Hidden     uint   `json:"hidden" validate:"oneof=1 2"`
	NoCache    uint   `json:"noCache" validate:"oneof=1 2"`
	AlwaysShow uint   `json:"alwaysShow" validate:"oneof=1 2"`
	Breadcrumb uint   `json:"breadcrumb" validate:"oneof=1 2"`
	ActiveMenu string `json:"activeMenu" validate:"min=0,max=100"`
	ParentId   uint   `json:"parentId" validate:"gte=0"`
}

// MenuDeleteReq is the resource deletion request struct
type MenuDeleteReq struct {
	MenuIds []uint `json:"menuIds" validate:"required"`
}

// MenuGetTreeReq is the menu tree request struct
type MenuGetTreeReq struct {
}

// MenuGetAccessTreeReq gets the user menu tree
type MenuGetAccessTreeReq struct {
	ID uint `json:"id" form:"id"`
}
