package helper

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		return "", err
	}
	return string(bytes), err
}

func CompareHashPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Info().Any("err", err).Msg("Error hash")
		return false
	}
	return true
}
