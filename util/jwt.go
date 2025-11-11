package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/config"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, config *config.Config) (string, error) {
	fmt.Printf("=== JWT GenerateToken 调试 ===\n")
	fmt.Printf("cfg 是否为 nil: %t\n", config == nil)
	if config != nil {
		//fmt.Printf("cfg.JWT 是否为 nil: %t\n", config.JWT. == nil)
		fmt.Printf("完整的 JWT 配置: %+v\n", config.JWT)
		fmt.Printf("Secret: %s\n", config.JWT.Secret)
		fmt.Printf("ExpirationHours: %d\n", config.JWT.ExpirationHours)
		fmt.Printf("Issuer: %s\n", config.JWT.Issuer)
		fmt.Printf("Audience: %s\n", config.JWT.Audience)
	}
	expirationTime := time.Now().Add(time.Duration(config.JWT.ExpirationHours) * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.JWT.Issuer,
			Audience:  []string{config.JWT.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT.Secret))
}

func ValidateToken(tokenString string, config *config.Config) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
