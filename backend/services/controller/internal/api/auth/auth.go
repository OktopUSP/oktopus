package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJwtKey() []byte {
	jwtKey, ok := os.LookupEnv("SECRET_API_KEY")
	if !ok || jwtKey == "" {
		return []byte("supersecretkey")
	}
	return []byte(jwtKey)
}

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(email string, username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(4 * time.Hour)
	claims := &JWTClaim{
		username,
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "Oktopus",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(getJwtKey())
	return
}

func ValidateToken(signedToken string) (email string, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return getJwtKey(), nil
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	email = claims.Email

	return
}
