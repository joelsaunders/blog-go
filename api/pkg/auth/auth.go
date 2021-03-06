package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/joelsaunders/blog-go/api/pkg/repository"
)

type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(id int, email string, jwtKey []byte) (string, error) {
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(20 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		ID:    id,
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

func CheckCredentials(ctx context.Context, email, password string, store repository.UserStore) (int, error) {
	storeUser, err := store.GetByEmail(ctx, email)
	if err != nil {
		return 0, fmt.Errorf("could not retreive user to check credentials: %s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storeUser.Password), []byte(password))

	if err != nil {
		return 0, fmt.Errorf("authenication failed: %s", err)
	}
	return storeUser.ID, nil
}
