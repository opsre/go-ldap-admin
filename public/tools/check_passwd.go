package tools

import (
	"fmt"
	"strings"
	"unicode"
)

/*
新密码规则：
1. 8位(包含)-16位(包含)
2. 包含至少一个拉丁大写字母
3. 包含至少一个拉丁小写字母
4. 包含至少一个数字
5.包含至少一个特殊字符  !@#%^&*()-_+={}][|;:<>,.?
6. 密码不能包含空格，特殊unicode字符
*/

// 校验新密码强度是否符合规则
func CheckPasswdStrength(newPasswd string) (err error) {
	// 校验密码长度
	if len(newPasswd) < 8 || len(newPasswd) > 16 {
		return fmt.Errorf("密码长度不符合规则")
	}

	// flags[0] hasUpper
	// flags[1] hasLower
	// flags[2] hasNumber
	// flags[3] hasSpecial
	var flags = make([]bool, 4)

	// 以unicode字符形式遍历
	for _, u := range newPasswd {
		// 密码不能包含空格
		if unicode.IsSpace(u) {
			return fmt.Errorf("密码不能包含空格")
		}

		// 如果既不是大小写字母，也不是数字，也不是标点符号
		// 为什么不用unicode.IsLetter()是为了排除如'æ'拉丁字母的干扰
		if !isLetter(u) && !unicode.IsNumber(u) && !isSpecial(u) {
			return fmt.Errorf("密码不能包含Unicode特殊字符")
		}

		// 判断是否为大写字母
		if unicode.IsUpper(u) {
			flags[0] = true
		}

		// 判断是否为小写字母
		if unicode.IsLower(u) {
			flags[1] = true
		}

		// 判断是否为数字
		if unicode.IsNumber(u) {
			flags[2] = true
		}

		// 判断特殊字符
		if isSpecial(u) {
			flags[3] = true
		}
	}

	// 判断四个条件是否满足
	for _, flag := range flags {
		if !flag {
			return fmt.Errorf("密码不符合规则")
		}
	}

	return nil
}

// 判断特殊字符
func isSpecial(u rune) bool {
	return strings.Contains("!@#%^&*()-_+={}][|;:<>,.?", string(u))
}

// 判断是否为拉丁字母
func isLetter(u rune) bool {
	return (u >= 'a' && u <= 'z') || (u >= 'A' && u <= 'Z')
}
