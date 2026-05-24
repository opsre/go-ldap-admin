package middleware

import (
	"testing"

	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/gin-gonic/gin"
)

func TestLocalizeAuthFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.SetSupportedLocales([]string{"zh-CN", "en-US"}, "zh-CN")
	i18n.SetMessages(map[string]map[string]string{
		"zh-CN": {
			"auth.not_logged_in":      "用户未登录",
			"auth.user_not_found":     "用户不存在",
			"auth.user_disabled":      "当前用户已被禁用",
			"auth.password_incorrect": "密码错误",
			"auth.jwt_failed":         "JWT认证失败, 错误码: {code}, 错误信息: {message}",
		},
		"en-US": {
			"auth.not_logged_in":      "user is not logged in",
			"auth.user_not_found":     "user does not exist",
			"auth.user_disabled":      "the current user has been disabled",
			"auth.password_incorrect": "incorrect password",
			"auth.jwt_failed":         "JWT authentication failed, code: {code}, message: {message}",
		},
	})

	c, _ := gin.CreateTestContext(nil)
	i18n.SetLocale(c, "en-US")

	tests := map[string]string{
		"用户未登录": "user is not logged in",
		"用户不存在": "user does not exist",
		"用户被禁用": "the current user has been disabled",
		"密码错误":  "incorrect password",
	}
	for raw, want := range tests {
		if got := localizeAuthFailure(c, raw); got != want {
			t.Fatalf("localizeAuthFailure(%q) = %q, want %q", raw, got, want)
		}
	}
}
