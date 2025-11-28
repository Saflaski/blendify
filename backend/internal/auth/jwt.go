package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/glog"
)

func (h *AuthHandler) CreateJWT(secret []byte, userID string) (string, error) {
	expiration := time.Second * time.Duration(h.jwtExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiredAT": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *AuthHandler) ValidateWithJWT() func(next http.Handler) http.Handler {
	// ctx := context.Background()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			glog.Info("In Middleware ValidateWithJWT ")
			//Get the JWT token from the client
			//Parse the JWT token and verify with secret
			//Compare the user and parsed JWT token
			//

		})
	}
}
