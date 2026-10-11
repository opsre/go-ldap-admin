package i18n

import (
	"os"
	"testing"
)

func TestNormalizeLocale(t *testing.T) {
	SetSupportedLocales([]string{"en-US", "zh-CN", "ja-JP", "es-ES", "ko-KR"}, "en-US")

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
		{name: "unknown", in: "it-IT", want: "en-US"},
		{name: "empty", in: "", want: "en-US"},
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
	SetSupportedLocales([]string{"en-US", "zh-CN"}, "en-US")
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
	if err := Init("en-US", []string{"en-US", "zh-CN", "ja-JP", "es-ES", "ko-KR"}); err != nil {
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

func TestSyncErrorChineseMessagesPreserveOriginalText(t *testing.T) {
	if err := Init("en-US", []string{"en-US", "zh-CN", "ja-JP", "es-ES", "ko-KR"}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		key  string
		args Args
		want string
	}{
		{"sqlsync.user_id_not_found", Args{"id": 17}, "用户ID[17]不存在"},
		{"sqlsync.user_info_failed", Args{"error": "detail"}, "获取用户信息失败: detail"},
		{"sqlsync.user_ldap_sync_failed", Args{"username": "alice", "error": "detail"}, "向LDAP同步用户[alice]失败：detail"},
		{"sqlsync.user_department_info_failed", Args{"username": "alice", "error": "detail"}, "用户[alice]根据部门ID获取部门信息失败: detail"},
		{"sqlsync.user_add_group_failed", Args{"username": "alice", "group": "Engineering", "error": "detail"}, "向Ldap添加用户[alice]到分组[Engineering]失败：detail"},
		{"sqlsync.user_status_update_failed", Args{"username": "alice", "error": "detail"}, "用户[alice]同步完毕之后更新状态失败：detail"},
		{"sqlsync.group_id_not_found", Args{"id": 23}, "分组ID[23]不存在"},
		{"sqlsync.group_info_failed", Args{"error": "detail"}, "获取分组信息失败: detail"},
		{"sqlsync.group_ldap_sync_failed", Args{"group": "Engineering", "error": "detail"}, "向LDAP同步分组[Engineering]失败：detail"},
		{"sqlsync.group_user_sync_failed", Args{"group": "Engineering", "username": "alice", "error": "detail"}, "同步分组[Engineering]之后处理分组内的用户[alice]失败：detail"},
		{"sqlsync.group_status_update_failed", Args{"group": "Engineering", "error": "detail"}, "分组[Engineering]同步完毕之后更新状态失败：detail"},
		{"openldap_sync.user_mysql_create_failed", Args{"error": "detail"}, "向MySQL创建用户失败：detail"},
		{"openldap_sync.user_department_info_failed", Args{"error": "detail"}, "根据部门ID获取部门信息失败detail"},
		{"openldap_sync.user_mysql_group_relation_failed", Args{"error": "detail"}, "向MySQL添加用户到分组关系失败：detail"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := T("zh-CN", tt.key, tt.args); got != tt.want {
				t.Fatalf("T(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestInitLoadsLocalesWhenWorkingDirectoryHasNoLocaleFiles(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatal(err)
		}
	})

	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if err := Init("en-US", []string{"en-US", "zh-CN"}); err != nil {
		t.Fatal(err)
	}
	if got := T("zh-CN", "common.success", nil); got != "成功" {
		t.Fatalf("T common.success = %q, want 成功", got)
	}
}
