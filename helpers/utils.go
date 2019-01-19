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

func CheckToken(tokenString string) jwt.Claims {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("not authorization")
		}
		return []byte(config.JWTConf.Secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		now := time.Now().Unix()
		if int64(claims["exp"].(float64)) > now {
			return claims
		} else {
			return nil
		}
	} else {
		return nil
	}
}
