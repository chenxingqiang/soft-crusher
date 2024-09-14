// File: internal/auth/jwt.go

package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chenxingqiang/soft-crusher/pkg/errors"
	"github.com/chenxingqiang/soft-crusher/pkg/logging"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your-secret-key") // In production, use a secure method to manage this key

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			errors.RespondWithError(w, http.StatusUnauthorized, "Missing authorization token")
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		logging.Info("User authenticated", zap.String("username", claims.Username))
		next.ServeHTTP(w, r)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, found := users[creds.Username]
	if !found {
		errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := GenerateToken(creds.Username)
	if err != nil {
		errors.HandleError(w, err, "Failed to generate token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
