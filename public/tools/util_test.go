package tools

import (
	"errors"
	"fmt"
	"testing"
)

func TestGenPass(t *testing.T) {
	fmt.Printf("密码为：%s\n", NewGenPasswd("123456"))
	// err := ComparePasswd("$2a$10$Fy8p0nCixgWKzLfO3SgdhOzAF7YolSt6dHj6QidDGYlzLJDpniXB6", "123456")
	// if err != nil {
	// 	fmt.Printf("密码错误：%s\n", err)
	// }
}

func TestArrUintCmp(t *testing.T) {
	a := []uint{1, 2, 3, 4, 6, 9}
	b := []uint{1, 2, 3, 4, 5, 6, 7}
	c, d := ArrUintCmp(a, b)
	fmt.Printf("%v\n", c)
	fmt.Printf("%v\n", d)
}

func TestSliceToString(t *testing.T) {
	a := []uint{1}
	fmt.Printf("%s\n", SliceToString(a, ","))
}

func TestEncodePass(t *testing.T) {
	// to encode a password into ssha
	hashed := EncodePass([]byte("testpass"))
	fmt.Println(string(hashed))
	// to validate a password against saved hash.
	if Matches([]byte(hashed), []byte("testpass")) {
		fmt.Println("Its a match.")
	} else {
		fmt.Println("its not match")
	}
}

// 测试密码强度是否符合规则
func TestNewPasswdStrength(t *testing.T) {
	tests := []struct {
		name string
		pwd  string
		err  error
	}{
		// TODO: Add test cases.
		{name: "[False]长度不足8", pwd: "abc@!Q", err: errors.New("pasword invalid")},
		{name: "[False]长度大于16", pwd: "abc@!Qaaassss3729128aaaaaaaaa", err: errors.New("pasword invalid")},
		{name: "[False]包含空格", pwd: "abAa a32!aCa", err: errors.New("pasword invalid")},
		{name: "[False]全大写", pwd: "AAGGGAAGAAAAAA", err: errors.New("pasword invalid")},
		{name: "[False]全小写", pwd: "abcabcabc", err: errors.New("pasword invalid")},
		{name: "[False]全数字", pwd: "111111111", err: errors.New("pasword invalid")},
		{name: "[False]全特殊字符", pwd: "!@#$%^&*(()_+", err: errors.New("pasword invalid")},

		{name: "[False]大写+小写", pwd: "abcabcQQQQ", err: errors.New("pasword invalid")},
		{name: "[False]大写+数字", pwd: "1111QQQQ", err: errors.New("pasword invalid")},
		{name: "[False]大写+特殊字符", pwd: "QQQQA(())", err: errors.New("pasword invalid")},
		{name: "[False]小写+数字", pwd: "abcabc7890", err: errors.New("pasword invalid")},
		{name: "[False]小写+特殊字符", pwd: "abcabc.(())", err: errors.New("pasword invalid")},
		{name: "[False]数字+特殊字符", pwd: "12345678&()", err: errors.New("pasword invalid")},
		{name: "[False]大写+小写+数字", pwd: "AAcabc7890", err: errors.New("pasword invalid")},
		{name: "[False]大写+小写+特殊字符", pwd: "aBcabcQ(A)Q", err: errors.New("pasword invalid")},
		{name: "[False]小写+数字+特殊字符", pwd: "abc1239)))", err: errors.New("pasword invalid")},

		{name: "[False]unicode特殊字符#00", pwd: "æææææAb#c1", err: errors.New("pasword invalid")},
		{name: "[False]unicode特殊字符#01", pwd: "\b5Ὂg̀9! ℃ᾭG", err: errors.New("pasword invalid")},
		{name: "[False]unicode特殊字符#02", pwd: "中文是密码Ý", err: errors.New("pasword invalid")},
		{name: "[False]unicode特殊字符#03", pwd: "notEng3.14Ý", err: errors.New("pasword invalid")},
		{name: "[False]unicode特殊字符#04", pwd: "ĒNĜĹis Ĥ", err: errors.New("pasword invalid")},
		{name: "[False]unicode特殊字符#05", pwd: "😂😂😂😂😂", err: errors.New("pasword invalid")},

		{name: "[True]大写-小写-数字-特殊字符", pwd: "AbcFAc.163606", err: nil},
		{name: "[True]大写-小写-数字-特殊字符", pwd: "AbcP1bc!!#.#", err: nil},
		{name: "[True]大写-小写数字-特殊字符", pwd: "AUbEXZ#14159", err: nil},
		{name: "[True]大写-小写-数字-特殊字符", pwd: "iyTmqp@14159", err: nil},
		{name: "[True]大写-小写-数字-特殊字符", pwd: "G1thub:liuup", err: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswdStrength(tt.pwd)
			if (err != nil && tt.err == nil) || (err == nil && tt.err != nil) {
				t.Errorf("CheckPasswdStrength() error = %v", err)
			}
		})
	}
}
