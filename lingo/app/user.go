package app

import (

	"golang.org/x/crypto/bcrypt"
)

func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
