package auth

import "golang.org/x/crypto/bcrypt"

func GererateHashedPassword(p string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p), 11)
	return string(hashed), err
}

func ComparePassword(value1, value2 string) error {
	return bcrypt.CompareHashAndPassword([]byte(value1), []byte(value2))
}
