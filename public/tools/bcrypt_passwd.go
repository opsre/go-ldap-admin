package tools

import (
	"crypto/rand"
	"math/big"

	"github.com/eryajf/go-ldap-admin/config"
)

// 密码加密 使用自适应hash算法, 不可逆
// func GenPasswd(passwd string) string {
// 	hashPasswd, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
// 	return string(hashPasswd)
// }

// 通过比较两个字符串hash判断是否出自同一个明文
// hashPasswd 需要对比的密文
// passwd 明文
// func ComparePasswd(hashPasswd string, passwd string) error {
// 	// if err := bcrypt.CompareHashAndPassword([]byte(hashPasswd), []byte(passwd)); err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

// 密码加密
func NewGenPasswd(passwd string) string {
	pass, _ := RSAEncrypt([]byte(passwd), config.Conf.System.RSAPublicBytes)
	return string(pass)
}

// 密码解密
func NewParPasswd(passwd string) string {
	pass, _ := RSADecrypt([]byte(passwd), config.Conf.System.RSAPrivateBytes)
	return string(pass)
}

const (
	passwordLength = 8
	letters        = "abcdefghijklmnopqrstu@vwxyzABCDEFGHIJKL#MNOP*QRSTUVWXYZ0123456789"
	lettersLength  = len(letters)
)

// 生成随机密码
func GenerateRandomPassword() string {
	password := make([]byte, passwordLength)

	for i := range password {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(lettersLength)))
		password[i] = letters[index.Int64()]
	}

	return string(password)
}
