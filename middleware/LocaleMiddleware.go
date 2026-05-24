package middleware

import (
	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/gin-gonic/gin"
)

func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.GetHeader("X-Locale")
		if locale == "" {
			locale = c.GetHeader("Accept-Language")
		}
		i18n.SetLocale(c, locale)
		c.Next()
	}
}
