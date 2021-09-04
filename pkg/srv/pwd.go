// Package srv password utilities
package srv

import "golang.org/x/crypto/bcrypt"

func GenPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	} else {
		return string(hash), nil
	}
}

func CheckPwd(hash string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	if err != nil {
		return false
	} else {
		return true
	}
}
