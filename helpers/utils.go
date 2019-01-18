package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/elfgzp/plumber/config"
	"time"
)

func GeneratePasswordHash(pwd string) string {
	return Md5(pwd)
}

func Md5(origin string) string {
	hasher := md5.New()
	hasher.Write([]byte(origin))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Second * time.Duration(config.JWTConf.ExpSec)).Unix(),
	})

	return token.SignedString([]byte(config.JWTConf.Secret))
}

func CheckToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.JWTConf.Secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["email"].(string), nil
	} else {
		return "", err
	}

}
