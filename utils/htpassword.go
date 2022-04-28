package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Htpasswd(user, pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", user, string(hashedPass)), nil
}
