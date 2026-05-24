package i18n

import "testing"

func TestNormalizeLocale(t *testing.T) {
	SetSupportedLocales([]string{"zh-CN", "en-US", "ja-JP", "es-ES", "ko-KR"}, "zh-CN")

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "exact", in: "en-US", want: "en-US"},
		{name: "lower underscore", in: "en_us", want: "en-US"},
		{name: "language family", in: "en-GB", want: "en-US"},
		{name: "language only", in: "ja", want: "ja-JP"},
		{name: "accept language", in: "es-MX,es;q=0.9,en;q=0.8", want: "es-ES"},
		{name: "unknown", in: "it-IT", want: "zh-CN"},
		{name: "empty", in: "", want: "zh-CN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeLocale(tt.in); got != tt.want {
				t.Fatalf("NormalizeLocale(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestTranslateWithArgsAndFallback(t *testing.T) {
	SetSupportedLocales([]string{"zh-CN", "en-US"}, "zh-CN")
	SetMessages(map[string]map[string]string{
		"zh-CN": {
			"common.success":     "成功",
			"menu.create_failed": "创建菜单失败: {error}",
		},
		"en-US": {
			"common.success":     "success",
			"menu.create_failed": "failed to create menu: {error}",
		},
	})

	if got := T("en-US", "common.success", nil); got != "success" {
		t.Fatalf("T success = %q", got)
	}
	if got := T("en-US", "menu.create_failed", Args{"error": "duplicate"}); got != "failed to create menu: duplicate" {
		t.Fatalf("T args = %q", got)
	}
	if got := T("en-US", "missing.key", nil); got != "missing.key" {
		t.Fatalf("T missing = %q", got)
	}
	if got := T("en-GB", "common.success", nil); got != "success" {
		t.Fatalf("T fallback locale = %q", got)
	}
}

func TestEmailTemplatesLoadForAllLocales(t *testing.T) {
	if err := Init("zh-CN", []string{"zh-CN", "en-US", "ja-JP", "es-ES", "ko-KR"}); err != nil {
		t.Fatal(err)
	}

	keys := []string{
		"email.password_reset_subject",
		"email.password_reset_body",
		"email.verification_code_subject",
		"email.verification_code_body",
		"email.user_creation_subject",
		"email.user_creation_body",
		"email.admin_password_reset_subject",
		"email.admin_password_reset_body",
	}
	for _, locale := range SupportedLocales() {
		for _, key := range keys {
			if got := T(locale, key, nil); got == key || got == "" {
				t.Fatalf("%s %s = %q", locale, key, got)
			}
		}
	}
}
