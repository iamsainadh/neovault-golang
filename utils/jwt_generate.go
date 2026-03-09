package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(jwtSecret string, userID string, role string, orgId string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iss":     "neo-vault",
	}

	if orgId != "" {
		claims["org_id"] = orgId
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}
