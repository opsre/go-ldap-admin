package tools

import (
	"crypto/rand"
	"math/big"

	"github.com/eryajf/go-ldap-admin/config"
)

// Password encryption uses an adaptive, irreversible hash algorithm
// func GenPasswd(passwd string) string {
// 	hashPasswd, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
// 	return string(hashPasswd)
// }

// Compare two string hashes to determine whether they came from the same plain text
// hashPasswd is the hash to compare
// passwd is the plain-text password
// func ComparePasswd(hashPasswd string, passwd string) error {
// 	// if err := bcrypt.CompareHashAndPassword([]byte(hashPasswd), []byte(passwd)); err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

// Encrypt the password
func NewGenPasswd(passwd string) string {
	pass, _ := RSAEncrypt([]byte(passwd), config.Conf.System.RSAPublicBytes)
	return string(pass)
}

// Decrypt the password
func NewParPasswd(passwd string) string {
	pass, _ := RSADecrypt([]byte(passwd), config.Conf.System.RSAPrivateBytes)
	return string(pass)
}

const (
	passwordLength = 8
	letters        = "abcdefghijklmnopqrstu@vwxyzABCDEFGHIJKL#MNOP*QRSTUVWXYZ0123456789"
	lettersLength  = len(letters)
)

// Generate a random password
func GenerateRandomPassword() string {
	password := make([]byte, passwordLength)

	for i := range password {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(lettersLength)))
		password[i] = letters[index.Int64()]
	}

	return string(password)
}
