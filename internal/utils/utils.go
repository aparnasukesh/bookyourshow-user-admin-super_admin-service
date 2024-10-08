package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
)

func GenCaptchaCode() (string, error) {

	codes := make([]byte, 6)
	if _, err := rand.Read(codes); err != nil {
		return "", err
	}

	for i := 0; i < 6; i++ {
		codes[i] = uint8(48 + (codes[i] % 10))
	}

	return string(codes), nil
}

func HashPassword(password string) string {
	data := []byte(password)
	password = fmt.Sprintf("%x", md5.Sum(data))
	return password
}

func VerifyPassword(requestPassword, dbPassword string) bool {
	requestPassword = fmt.Sprintf("%x", md5.Sum([]byte(requestPassword)))
	return requestPassword == dbPassword
}
