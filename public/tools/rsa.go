package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// RSA encryption
func RSAEncrypt(data, publicBytes []byte) ([]byte, error) {
	var res []byte
	// Parse the public key
	block, _ := pem.Decode(publicBytes)

	if block == nil {
		return res, fmt.Errorf("failed to encrypt: invalid public key")
	}

	// Parse the decoded data with X.509
	// x509.MarshalPKCS1PublicKey(block) cannot be used after parsing, so use ParsePKIXPublicKey instead
	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return res, fmt.Errorf("failed to encrypt: invalid public key: %v", err)
	}
	// Encrypt data with the public key
	pubKey := keyInit.(*rsa.PublicKey)
	res, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return res, fmt.Errorf("failed to encrypt: invalid public key: %v", err)
	}
	// Encode the encrypted data as Base64
	return []byte(EncodeStr2Base64(string(res))), nil
}

// Decrypt the data
func RSADecrypt(base64Data, privateBytes []byte) ([]byte, error) {
	var res []byte
	// Parse the Base64 data
	data := []byte(DecodeStrFromBase64(string(base64Data)))
	// Parse the private key
	block, _ := pem.Decode(privateBytes)
	if block == nil {
		return res, fmt.Errorf("failed to decrypt: invalid private key")
	}
	// Restore the data
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return res, fmt.Errorf("failed to decrypt: invalid private key: %v", err)
	}
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return res, fmt.Errorf("failed to decrypt: invalid private key: %v", err)
	}
	return res, nil
}

// Encrypt a Base64 string
func EncodeStr2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Decrypt a Base64 string
func DecodeStrFromBase64(str string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes)
}
