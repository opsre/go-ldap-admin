package request

// BaseSendCodeReq sends a verification code
type BaseSendCodeReq struct {
	Mail string `json:"mail" validate:"required,min=0,max=100"`
}

// BaseChangePwdReq is the password change request struct
type BaseChangePwdReq struct {
	Mail string `json:"mail" validate:"required,min=0,max=100"`
	Code string `json:"code" validate:"required,len=6"`
}

// BaseDashboardReq is the system home page data request struct
type BaseDashboardReq struct {
}

// EncryptPasswdReq
type EncryptPasswdReq struct {
	Passwd string `json:"passwd" form:"passwd" validate:"required"`
}

// DecryptPasswdReq
type DecryptPasswdReq struct {
	Passwd string `json:"passwd" form:"passwd" validate:"required"`
}

// BaseConfigReq is the system configuration request struct
type BaseConfigReq struct {
}

// BaseVersionReq is the version information request struct
type BaseVersionReq struct {
}
