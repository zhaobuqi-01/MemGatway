package pkg

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// GenerateSaltedMD5 generates an MD5 hash with a salt.
func GenerateSaltedMD5(input, salt string) string {
	hasher := md5.New()
	hasher.Write([]byte(input + salt))
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}

// GenerateSalt generates a cryptographically secure random salt of the given length.
func GenerateSalt(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	salt := base64.URLEncoding.EncodeToString(randomBytes)
	return salt, nil
}

func GenSaltPassword(salt, password string) string {
	s1 := sha256.New()
	s1.Write([]byte(password))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(salt + str1))
	return fmt.Sprintf("%x", s2.Sum(nil))
}
