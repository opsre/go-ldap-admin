package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Return to the frontend
func Response(c *gin.Context, httpStatus int, code int, data gin.H, message string) {
	c.JSON(httpStatus, gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	})
}

// Return to the frontend - success
func Success(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusOK, 200, data, message)
}

// Return to the frontend - failure
func Fail(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusBadRequest, 400, data, message)
}
