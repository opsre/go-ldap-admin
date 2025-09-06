package tools

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SystemErr    = 500
	MySqlErr     = 501
	LdapErr      = 505
	OperationErr = 506
	ValidatorErr = 412
)

type RspError struct {
	code int
	err  error
}

func (re *RspError) Error() string {
	return re.err.Error()
}

func (re *RspError) Code() int {
	return re.code
}

// NewRspError New
func NewRspError(code int, err error) *RspError {
	return &RspError{
		code: code,
		err:  err,
	}
}

// NewMySqlError mysql错误
func NewMySqlError(err error) *RspError {
	return NewRspError(MySqlErr, err)
}

// NewValidatorError 验证错误
func NewValidatorError(err error) *RspError {
	return NewRspError(ValidatorErr, err)
}

// NewLdapError ldap错误
func NewLdapError(err error) *RspError {
	return NewRspError(LdapErr, err)
}

// NewOperationError 操作错误
func NewOperationError(err error) *RspError {
	return NewRspError(OperationErr, err)
}

// ReloadErr 重新加载错误
func ReloadErr(err any) *RspError {
	rspErr, ok := err.(*RspError)
	if !ok {
		rspError, ok := err.(error)
		if !ok {
			return &RspError{
				code: SystemErr,
				err:  fmt.Errorf("unknow error"),
			}
		}
		return &RspError{
			code: SystemErr,
			err:  rspError,
		}
	}
	return rspErr
}

// Success http 成功
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

// Err http 错误
func Err(c *gin.Context, err *RspError, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": err.Code(),
		"msg":  err.Error(),
		"data": data,
	})
}

// 返回前端
func Response(c *gin.Context, httpStatus int, code int, data gin.H, message string) {
	c.JSON(httpStatus, gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	})
}
