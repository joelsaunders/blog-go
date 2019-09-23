package auth

import (
	"log"

	"github.com/joelsaunders/bilbo-go/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateNewUserObj(email string, password string) *models.NewUser {
	passwordBytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	hashedPassword := string(hash)
	return &models.NewUser{Email: email, Password: hashedPassword}
}
