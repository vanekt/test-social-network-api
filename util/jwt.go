package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

const JwtSecretKey = "81ni1ke0jz122HHz8m8ZQR9m3a3pdIXmTRiELTOSqe9xEcfbJoZ47SwlC94k8XwO"

func CreateToken(userId uint32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
	})

	tokenString, err := token.SignedString([]byte(JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserIdFromToken(tokenString string) (uint32, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecretKey), nil
	})

	if err != nil {
		return 0, errors.New("GetUserIdFromToken: Can't parse auth token: " + err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("GetUserIdFromToken: Something went wrong")
	}

	return uint32(claims["userId"].(float64)), nil
}
