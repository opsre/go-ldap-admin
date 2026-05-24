package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eryajf/go-ldap-admin/public/i18n"
	"github.com/gin-gonic/gin"
)

func TestLocaleMiddlewareUsesXLocale(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.SetSupportedLocales([]string{"zh-CN", "en-US"}, "zh-CN")

	r := gin.New()
	r.Use(LocaleMiddleware())
	r.GET("/locale", func(c *gin.Context) {
		c.String(http.StatusOK, i18n.LocaleFromContext(c))
	})

	req := httptest.NewRequest(http.MethodGet, "/locale", nil)
	req.Header.Set("X-Locale", "en-US")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "en-US" {
		t.Fatalf("locale = %q", rec.Body.String())
	}
}

func TestLocaleMiddlewareFallsBackToAcceptLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.SetSupportedLocales([]string{"zh-CN", "en-US"}, "zh-CN")

	r := gin.New()
	r.Use(LocaleMiddleware())
	r.GET("/locale", func(c *gin.Context) {
		c.String(http.StatusOK, i18n.LocaleFromContext(c))
	})

	req := httptest.NewRequest(http.MethodGet, "/locale", nil)
	req.Header.Set("Accept-Language", "en-GB,en;q=0.8")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "en-US" {
		t.Fatalf("locale = %q", rec.Body.String())
	}
}
