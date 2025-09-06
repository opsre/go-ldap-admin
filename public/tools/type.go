package tools

import "github.com/mozillazg/go-pinyin"

// H is a shortcut for map[string]any
type H map[string]any

func ConvertToPinYin(src string) (dst string) {
	args := pinyin.NewArgs()
	args.Fallback = func(r rune, args pinyin.Args) []string {
		return []string{string(r)}
	}

	for _, singleResult := range pinyin.Pinyin(src, args) {
		for _, result := range singleResult {
			dst = dst + result
		}
	}
	return
}
