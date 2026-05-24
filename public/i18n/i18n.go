package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

const (
	ContextLocaleKey = "locale"
	defaultLocale    = "zh-CN"
)

type Args map[string]any

var (
	mu               sync.RWMutex
	activeDefault    = defaultLocale
	supportedLocales = []string{"zh-CN", "en-US", "ja-JP", "es-ES", "ko-KR"}
	messages         = map[string]map[string]string{}
)

func Init(defaultLoc string, supported []string) error {
	SetSupportedLocales(supported, defaultLoc)

	loaded := make(map[string]map[string]string)
	for _, locale := range SupportedLocales() {
		data, err := readLocaleFile(locale)
		if err != nil {
			return fmt.Errorf("read locale %s: %w", locale, err)
		}
		flat, err := flattenYAML(data)
		if err != nil {
			return fmt.Errorf("parse locale %s: %w", locale, err)
		}
		loaded[locale] = flat
	}
	SetMessages(loaded)
	return nil
}

func readLocaleFile(locale string) ([]byte, error) {
	filename := filepath.Join("locales", locale+".yaml")
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for {
		candidate := filepath.Join(dir, filename)
		if data, err := os.ReadFile(candidate); err == nil {
			return data, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil, os.ErrNotExist
}

func SetSupportedLocales(locales []string, defaultLoc string) {
	mu.Lock()
	defer mu.Unlock()

	if len(locales) == 0 {
		locales = []string{defaultLocale}
	}
	normalized := make([]string, 0, len(locales))
	for _, locale := range locales {
		if norm := canonicalLocale(locale); norm != "" {
			normalized = append(normalized, norm)
		}
	}
	if len(normalized) == 0 {
		normalized = []string{defaultLocale}
	}
	supportedLocales = dedupe(normalized)

	defaultLoc = canonicalLocale(defaultLoc)
	if !contains(supportedLocales, defaultLoc) {
		defaultLoc = supportedLocales[0]
	}
	activeDefault = defaultLoc
}

func SupportedLocales() []string {
	mu.RLock()
	defer mu.RUnlock()

	out := make([]string, len(supportedLocales))
	copy(out, supportedLocales)
	return out
}

func SetMessages(next map[string]map[string]string) {
	mu.Lock()
	defer mu.Unlock()

	messages = next
}

func NormalizeLocale(raw string) string {
	mu.RLock()
	defer mu.RUnlock()

	return normalizeLocaleLocked(raw)
}

func LocaleFromContext(c *gin.Context) string {
	if c == nil {
		return NormalizeLocale("")
	}
	if value, ok := c.Get(ContextLocaleKey); ok {
		if locale, ok := value.(string); ok {
			return NormalizeLocale(locale)
		}
	}
	return NormalizeLocale("")
}

func SetLocale(c *gin.Context, locale string) {
	c.Set(ContextLocaleKey, NormalizeLocale(locale))
}

func TC(c *gin.Context, key string, args Args) string {
	return T(LocaleFromContext(c), key, args)
}

func T(locale string, key string, args Args) string {
	mu.RLock()
	defer mu.RUnlock()

	locale = normalizeLocaleLocked(locale)
	if msg := messages[locale][key]; msg != "" {
		return format(msg, args)
	}
	if msg := messages[activeDefault][key]; msg != "" {
		return format(msg, args)
	}
	if activeDefault != defaultLocale {
		if msg := messages[defaultLocale][key]; msg != "" {
			return format(msg, args)
		}
	}
	return key
}

func normalizeLocaleLocked(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return activeDefault
	}

	for _, part := range strings.Split(raw, ",") {
		token := strings.TrimSpace(strings.Split(part, ";")[0])
		if token == "" {
			continue
		}
		token = canonicalLocale(token)
		if contains(supportedLocales, token) {
			return token
		}
		lang := strings.Split(token, "-")[0]
		for _, supported := range supportedLocales {
			if strings.HasPrefix(supported, lang+"-") {
				return supported
			}
		}
	}
	return activeDefault
}

func canonicalLocale(raw string) string {
	raw = strings.TrimSpace(strings.ReplaceAll(raw, "_", "-"))
	if raw == "" {
		return ""
	}
	parts := strings.Split(raw, "-")
	if len(parts) == 1 {
		return strings.ToLower(parts[0])
	}
	return strings.ToLower(parts[0]) + "-" + strings.ToUpper(parts[1])
}

func flattenYAML(data []byte) (map[string]string, error) {
	var nested map[string]any
	if err := yaml.Unmarshal(data, &nested); err != nil {
		return nil, err
	}
	out := map[string]string{}
	flatten("", nested, out)
	return out, nil
}

func flatten(prefix string, value any, out map[string]string) {
	switch typed := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			next := key
			if prefix != "" {
				next = prefix + "." + key
			}
			flatten(next, typed[key], out)
		}
	case string:
		out[prefix] = typed
	default:
		out[prefix] = fmt.Sprint(typed)
	}
}

func format(message string, args Args) string {
	for key, value := range args {
		message = strings.ReplaceAll(message, "{"+key+"}", fmt.Sprint(value))
	}
	return message
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func dedupe(items []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}
