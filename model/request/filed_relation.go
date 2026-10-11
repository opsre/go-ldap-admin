package request

// FieldRelationListReq is the resource list request struct
type FieldRelationListReq struct {
}

// FieldRelationAddReq is the resource creation request struct
type FieldRelationAddReq struct {
	Flag       string            `json:"flag" validate:"required,min=1,max=20"`
	Attributes map[string]string `json:"attributes" validate:"required,gt=0"`
}

// FieldRelationUpdateReq is the resource update request struct
type FieldRelationUpdateReq struct {
	ID         uint              `json:"id" validate:"required"`
	Flag       string            `json:"flag" validate:"required,min=1,max=20"`
	Attributes map[string]string `json:"attributes" validate:"required,gt=0"`
}

// FieldRelationDeleteReq is the resource deletion request struct
type FieldRelationDeleteReq struct {
	FieldRelationIds []uint `json:"fieldRelationIds" validate:"required"`
}
