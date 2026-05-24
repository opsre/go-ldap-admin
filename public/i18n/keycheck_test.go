package i18n

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestLocaleFilesHaveSameKeys(t *testing.T) {
	if err := Init("zh-CN", []string{"zh-CN", "en-US", "ja-JP", "es-ES", "ko-KR"}); err != nil {
		t.Fatal(err)
	}

	mu.RLock()
	defer mu.RUnlock()

	base := messages["zh-CN"]
	for locale, localeMessages := range messages {
		if locale == "zh-CN" {
			continue
		}
		for key := range base {
			if _, ok := localeMessages[key]; !ok {
				t.Fatalf("%s missing key %s", locale, key)
			}
		}
		for key := range localeMessages {
			if _, ok := base[key]; !ok {
				t.Fatalf("%s has extra key %s", locale, key)
			}
		}
	}
}

func TestMigratedBackendModulesDoNotUseChineseRawResponseErrors(t *testing.T) {
	files := []string{
		"logic/api_logic.go",
		"logic/operation_log_logic.go",
		"logic/field_relation_logic.go",
		"logic/role_logic.go",
		"logic/group_logic.go",
		"logic/user_logic.go",
		"logic/menu_logic.go",
		"logic/base_logic.go",
		"logic/dingtalk_logic.go",
		"logic/feishu_logic.go",
		"logic/wecom_logic.go",
		"logic/openldap_logic.go",
		"logic/a_logic.go",
		"service/ildap/group_ildap.go",
		"service/isql/api_isql.go",
		"service/isql/role_isql.go",
	}
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`New(Validator|MySql|Ldap|Operation)Error\(fmt\.Errorf\("[^"]*[\x{4e00}-\x{9fff}]`),
		regexp.MustCompile(`return "[^"]*[\x{4e00}-\x{9fff}]`),
	}
	root := findRepoRoot(t)
	for _, file := range files {
		source, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			t.Fatal(err)
		}
		for _, pattern := range patterns {
			if match := pattern.Find(source); match != nil {
				t.Fatalf("%s contains raw Chinese response error: %s", file, string(match))
			}
		}
	}
}

func TestServiceLayerErrorsThatReachResponsesDoNotUseChinese(t *testing.T) {
	files := []string{
		"service/isql/api_isql.go",
		"service/isql/role_isql.go",
		"service/ildap/group_ildap.go",
	}
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`errors\.New\("[^"]*[\x{4e00}-\x{9fff}]`),
		regexp.MustCompile(`fmt\.Errorf\("[^"]*[\x{4e00}-\x{9fff}]`),
	}
	root := findRepoRoot(t)
	for _, file := range files {
		source, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			t.Fatal(err)
		}
		for _, pattern := range patterns {
			if match := pattern.Find(source); match != nil {
				t.Fatalf("%s contains raw Chinese service error: %s", file, string(match))
			}
		}
	}
}

func TestUserServiceChineseSentinelsAreLimitedToAuthMapping(t *testing.T) {
	root := findRepoRoot(t)
	source, err := os.ReadFile(filepath.Join(root, "service/isql/user_isql.go"))
	if err != nil {
		t.Fatal(err)
	}
	allowed := map[string]bool{
		`errors.New("用户未登录")`: true,
		`errors.New("用户不存在")`: true,
		`errors.New("用户被禁用")`: true,
		`errors.New("密码错误")`:  true,
	}
	pattern := regexp.MustCompile(`errors\.New\("[^"]*[\x{4e00}-\x{9fff}][^"]*"\)`)
	for _, match := range pattern.FindAllString(string(source), -1) {
		if !allowed[match] {
			t.Fatalf("service/isql/user_isql.go contains unmapped Chinese sentinel: %s", match)
		}
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}
