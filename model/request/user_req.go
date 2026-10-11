package request

// UserAddReq is the resource creation request struct
type UserAddReq struct {
	Username      string `json:"username" validate:"required,min=2,max=50"`
	Password      string `json:"password"`
	Nickname      string `json:"nickname" validate:"required,min=0,max=50"`
	GivenName     string `json:"givenName" validate:"min=0,max=50"`
	Mail          string `json:"mail" validate:"required,min=0,max=100"`
	JobNumber     string `json:"jobNumber" validate:"required,min=0,max=20"`
	PostalAddress string `json:"postalAddress" validate:"min=0,max=255"`
	Departments   string `json:"departments" validate:"min=0,max=512"`
	Position      string `json:"position" validate:"min=0,max=128"`
	Mobile        string `json:"mobile" validate:"required,checkMobile"`
	Avatar        string `json:"avatar"`
	Introduction  string `json:"introduction" validate:"min=0,max=255"`
	Status        uint   `json:"status" validate:"oneof=1 2"`
	DepartmentId  []uint `json:"departmentId" validate:"required"`
	Source        string `json:"source" validate:"min=0,max=50"`
	RoleIds       []uint `json:"roleIds" validate:"required"`
}

// DingUserAddReq is the DingTalk user creation request struct
type DingUserAddReq struct {
	Username      string `json:"username" validate:"required,min=2,max=50"`
	Password      string `json:"password"`
	Nickname      string `json:"nickname" validate:"required,min=0,max=50"`
	GivenName     string `json:"givenName" validate:"min=0,max=50"`
	Mail          string `json:"mail" validate:"required,min=0,max=100"`
	JobNumber     string `json:"jobNumber" validate:"required,min=0,max=20"`
	PostalAddress string `json:"postalAddress" validate:"min=0,max=255"`
	Departments   string `json:"departments" validate:"min=0,max=512"`
	Position      string `json:"position" validate:"min=0,max=128"`
	Mobile        string `json:"mobile" validate:"required,checkMobile"`
	Avatar        string `json:"avatar"`
	Introduction  string `json:"introduction" validate:"min=0,max=255"`
	Status        uint   `json:"status" validate:"oneof=1 2"`
	DepartmentId  []uint `json:"departmentId" validate:"required"`
	Source        string `json:"source" validate:"min=0,max=50"`
	RoleIds       []uint `json:"roleIds" validate:"required"`
	SourceUserId  string `json:"sourceUserId"`  // Third-party user ID
	SourceUnionId string `json:"sourceUnionId"` // Unique third-party union ID
}

// WeComUserAddReq is the WeCom user creation request struct
type WeComUserAddReq struct {
	Username      string `json:"username" validate:"required,min=2,max=50"`
	Password      string `json:"password"`
	Nickname      string `json:"nickname" validate:"required,min=0,max=50"`
	GivenName     string `json:"givenName" validate:"min=0,max=50"`
	Mail          string `json:"mail" validate:"required,min=0,max=100"`
	JobNumber     string `json:"jobNumber" validate:"required,min=0,max=20"`
	PostalAddress string `json:"postalAddress" validate:"min=0,max=255"`
	Departments   string `json:"departments" validate:"min=0,max=512"`
	Position      string `json:"position" validate:"min=0,max=128"`
	Mobile        string `json:"mobile" validate:"required,checkMobile"`
	Avatar        string `json:"avatar"`
	Introduction  string `json:"introduction" validate:"min=0,max=255"`
	Status        uint   `json:"status" validate:"oneof=1 2"`
	DepartmentId  []uint `json:"departmentId" validate:"required"`
	Source        string `json:"source" validate:"min=0,max=50"`
	RoleIds       []uint `json:"roleIds" validate:"required"`
	SourceUserId  string `json:"sourceUserId"`  // Third-party user ID
	SourceUnionId string `json:"sourceUnionId"` // Unique third-party union ID
}

// UserUpdateReq is the resource update request struct
type UserUpdateReq struct {
	ID            uint   `json:"id" validate:"required"`
	Username      string `json:"username" validate:"required,min=2,max=50"`
	Nickname      string `json:"nickname" validate:"min=0,max=20"`
	GivenName     string `json:"givenName" validate:"min=0,max=50"`
	Mail          string `json:"mail" validate:"min=0,max=100"`
	JobNumber     string `json:"jobNumber" validate:"min=0,max=20"`
	PostalAddress string `json:"postalAddress" validate:"min=0,max=255"`
	Departments   string `json:"departments" validate:"min=0,max=512"`
	Position      string `json:"position" validate:"min=0,max=128"`
	Mobile        string `json:"mobile" validate:"checkMobile"`
	Avatar        string `json:"avatar"`
	Introduction  string `json:"introduction" validate:"min=0,max=255"`
	DepartmentId  []uint `json:"departmentId" validate:"required"`
	Source        string `json:"source" validate:"min=0,max=50"`
	RoleIds       []uint `json:"roleIds" validate:"required"`
}

// UserDeleteReq is the bulk resource deletion request struct
type UserDeleteReq struct {
	UserIds []uint `json:"userIds" validate:"required"`
}

// UserChangePwdReq is the password change request struct
type UserChangePwdReq struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

// UserResetPasswordReq is the password reset request struct
type UserResetPasswordReq struct {
	Username string `json:"username" validate:"required"`
}

// UserChangeUserStatusReq is the user status change request struct
type UserChangeUserStatusReq struct {
	ID     uint `json:"id" validate:"required"`
	Status uint `json:"status" validate:"oneof=1 2"`
}

// UserGetUserInfoReq is the user information request struct
type UserGetUserInfoReq struct {
}

// SyncDingUserReq syncs DingTalk user information
type SyncDingUserReq struct {
}

// SyncWeComUserReq syncs WeCom user information
type SyncWeComUserReq struct {
}

// SyncFeiShuUserReq syncs Feishu user information
type SyncFeiShuUserReq struct {
}

// SyncOpenLdapUserReq syncs LDAP user information
type SyncOpenLdapUserReq struct {
}
type SyncSqlUserReq struct {
	UserIds []uint `json:"userIds" validate:"required"`
}

// UserListReq is the user list request struct
type UserListReq struct {
	Username     string `json:"username" form:"username"`
	Mobile       string `json:"mobile" form:"mobile" `
	Nickname     string `json:"nickname" form:"nickname"`
	GivenName    string `json:"givenName" form:"givenName"`
	DepartmentId []uint `json:"departmentId" form:"departmentId"`
	Status       uint   `json:"status" form:"status" `
	SyncState    uint   `json:"syncState" form:"syncState" `
	PageNum      int    `json:"pageNum" form:"pageNum"`
	PageSize     int    `json:"pageSize" form:"pageSize"`
}

// RegisterAndLoginReq is the user login request struct
type RegisterAndLoginReq struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
