package token

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
)

func StudentIdForToken(tokenString string) string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Println(err)
		return ""
	}
	claims := token.Claims.(jwt.MapClaims)
	studentId := claims["studentId"].(string)
	return studentId
}
