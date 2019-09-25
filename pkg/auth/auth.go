package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joelsaunders/bilbo-go/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(email string, jwtKey []byte) (string, error) {
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(10 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	return token.SignedString(jwtKey)
}

func HashPassword(password string) string {
	passwordBytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	hashedPassword := string(hash)
	return hashedPassword
}

func CheckCredentials(ctx context.Context, email, password string, store repository.UserStore) error {
	storeUser, err := store.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("could not retreive user to check credentials: %s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storeUser.Password), []byte(password))

	if err != nil {
		return fmt.Errorf("authenication failed: %s", err)
	}
	return nil
}
