package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func Validate(tokenString string) (uint, float64, error) {

	hmacSampleSecret := []byte(GoDotEnvVariable("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})

	if err != nil {
		return 0, 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		return uint(claims["user"].(float64)), claims["role"].(float64), nil

	} else {
		return 0, 0, err
	}
}
