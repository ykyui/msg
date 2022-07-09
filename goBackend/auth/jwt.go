package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte("hmacSampleSecret")

func GenerateToken(userId string) (string, error) {
	// Create the JWT claims, which includes the username and expiry time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, struct {
		jwt.StandardClaims
	}{
		StandardClaims: jwt.StandardClaims{
			Id:        userId,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
