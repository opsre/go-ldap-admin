package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/gin-gonic/gin"
)

func TestErrKeepsRawErrorCompatibility(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.SetSupportedLocales([]string{"zh-CN", "en-US"}, "zh-CN")
	i18n.SetMessages(map[string]map[string]string{
		"zh-CN": {"common.success": "success"},
		"en-US": {"common.success": "success"},
	})

	r := gin.New()
	r.GET("/err", func(c *gin.Context) {
		Err(c, NewValidatorError(fmt.Errorf("用户名已存在")), nil)
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/err", nil))

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["msg"] != "用户名已存在" {
		t.Fatalf("msg = %#v", body["msg"])
	}
	if _, ok := body["msgKey"]; ok {
		t.Fatalf("raw error should not include msgKey: %#v", body)
	}
}

func TestErrLocalizesUnmappedChineseErrorForNonDefaultLocale(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.SetSupportedLocales([]string{"zh-CN", "en-US"}, "zh-CN")
	i18n.SetMessages(map[string]map[string]string{
		"zh-CN": {
			"error.validator": "参数校验失败",
		},
		"en-US": {
			"error.validator": "validation failed",
		},
	})

	r := gin.New()
	r.GET("/err", func(c *gin.Context) {
		i18n.SetLocale(c, "en-US")
		Err(c, NewValidatorError(fmt.Errorf("对应平台的动态字段关系已存在，请勿重复添加")), nil)
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/err", nil))

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["msg"] != "validation failed" {
		t.Fatalf("msg = %#v", body["msg"])
	}
	if body["msgKey"] != "error.validator" {
		t.Fatalf("msgKey = %#v", body["msgKey"])
	}
}

func TestErrLocalizesI18nError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.SetSupportedLocales([]string{"zh-CN", "en-US"}, "zh-CN")
	i18n.SetMessages(map[string]map[string]string{
		"zh-CN": {"user.username_exists": "用户名已存在"},
		"en-US": {"user.username_exists": "username already exists"},
	})

	r := gin.New()
	r.GET("/err", func(c *gin.Context) {
		i18n.SetLocale(c, "en-US")
		Err(c, NewValidatorI18nError("user.username_exists", nil), nil)
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/err", nil))

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["msg"] != "username already exists" {
		t.Fatalf("msg = %#v", body["msg"])
	}
	if body["msgKey"] != "user.username_exists" {
		t.Fatalf("msgKey = %#v", body["msgKey"])
	}
}

func TestErrLocalizesFieldRelationI18nError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.SetSupportedLocales([]string{"zh-CN", "en-US"}, "zh-CN")
	i18n.SetMessages(map[string]map[string]string{
		"zh-CN": {"field_relation.exists": "对应平台的动态字段关系已存在，请勿重复添加"},
		"en-US": {"field_relation.exists": "the dynamic field relation for this platform already exists"},
	})

	r := gin.New()
	r.GET("/err", func(c *gin.Context) {
		i18n.SetLocale(c, "en-US")
		Err(c, NewValidatorI18nError("field_relation.exists", nil), nil)
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/err", nil))

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["msg"] != "the dynamic field relation for this platform already exists" {
		t.Fatalf("msg = %#v", body["msg"])
	}
	if body["msgKey"] != "field_relation.exists" {
		t.Fatalf("msgKey = %#v", body["msgKey"])
	}
}

func TestErrLocalizesLegacyRawError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.SetSupportedLocales([]string{"zh-CN", "en-US"}, "zh-CN")
	i18n.SetMessages(map[string]map[string]string{
		"zh-CN": {
			"legacy.common.current_user_failed": "获取当前登陆用户信息失败",
			"legacy.api.create_failed":          "创建接口失败: {error}",
		},
		"en-US": {
			"legacy.common.current_user_failed": "failed to get current login user information",
			"legacy.api.create_failed":          "failed to create API: {error}",
		},
	})

	r := gin.New()
	r.GET("/exact", func(c *gin.Context) {
		i18n.SetLocale(c, "en-US")
		Err(c, NewMySqlError(fmt.Errorf("获取当前登陆用户信息失败")), nil)
	})
	r.GET("/prefix", func(c *gin.Context) {
		i18n.SetLocale(c, "en-US")
		Err(c, NewMySqlError(fmt.Errorf("创建接口失败: duplicate")), nil)
	})

	tests := []struct {
		path       string
		wantMsg    string
		wantMsgKey string
	}{
		{
			path:       "/exact",
			wantMsg:    "failed to get current login user information",
			wantMsgKey: "legacy.common.current_user_failed",
		},
		{
			path:       "/prefix",
			wantMsg:    "failed to create API: duplicate",
			wantMsgKey: "legacy.api.create_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body["msg"] != tt.wantMsg {
				t.Fatalf("msg = %#v", body["msg"])
			}
			if body["msgKey"] != tt.wantMsgKey {
				t.Fatalf("msgKey = %#v", body["msgKey"])
			}
		})
	}
}
