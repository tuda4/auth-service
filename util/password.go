package util

import "golang.org/x/crypto/bcrypt"

func CreateHashPassword(password string) (hashPassword string, err error) {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	hashPassword = string(passwordByte)
	return
}

func CheckHashPassword(password string, hashPassword string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
