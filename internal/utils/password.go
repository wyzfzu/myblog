package utils

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(pwd string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPassword), nil
}

func EqualPassword(hashPwd string, loginPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(loginPwd)); err != nil {
		return false
	}

	return true
}
