package tools

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

// code from https://gist.github.com/newm4n/ce9ac4308ae1beb4250efddad28e3f85

// Encode encodes the []byte of raw password
func EncodePass(rawPassPhrase []byte) string {
	hash := makeSSHAHash(rawPassPhrase, makeSalt())
	b64 := base64.StdEncoding.EncodeToString(hash)
	return string([]byte(fmt.Sprintf("{SSHA}%s", b64)))
}

// makeSalt make a 4 byte array containing random bytes.
func makeSalt() []byte {
	sbytes := make([]byte, 4)
	if _, err := rand.Read(sbytes); err != nil {
		panic(err) // crypto/rand should never fail
	}
	return sbytes
}

// makeSSHAHash make hasing using SHA-1 with salt. This is not the final output though. You need to append {SSHA} string with base64 of this hash.
func makeSSHAHash(passphrase, salt []byte) []byte {
	sha := sha1.New()
	sha.Write(passphrase)
	sha.Write(salt)

	h := sha.Sum(nil)
	return append(h, salt...)
}

// Matches matches the encoded password and the raw password
func Matches(encodedPassPhrase, rawPassPhrase []byte) bool {
	//strip the {SSHA}
	eppS := string(encodedPassPhrase)[6:]
	hash, err := base64.StdEncoding.DecodeString(eppS)
	if err != nil {
		return false
	}
	salt := hash[len(hash)-4:]

	sha := sha1.New()
	sha.Write(rawPassPhrase)
	sha.Write(salt)
	sum := sha.Sum(nil)

	return bytes.Equal(sum, hash[:len(hash)-4])
}
