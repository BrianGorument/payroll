package utils

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/spf13/viper"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
func getJwtSecret() []byte {
	jwtSecret := viper.GetString("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in environment variables")
	}
	return []byte(jwtSecret)
}
func CreateJWTToken(userID int ,userName string , role string) (string, error) {

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userID,
		"userName":  userName,
		"role":      role,
		"exp"  :  time.Now().Add(time.Hour * 1).Unix(),
	})
	
	jwtSecret := getJwtSecret()
	
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	jwtSecret := getJwtSecret()
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}


func ConvertToUint(value interface{}) (uint, error) {
	floatVal, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("failed to convert %v to float64", value)
	}
	return uint(floatVal), nil
}

func ConvertToInt(value interface{}) (int, error) {
	floatVal, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("failed to convert %v to float64", value)
	}
	return int(floatVal), nil
}


