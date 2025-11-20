package utils

import (
	"errors"
	"time"

	"github.com/demo/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义 JWT Claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成 JWT Token
func GenerateJWT(userID int64, username string) (string, error) {
	cfg := config.GetConfig().JWT

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(cfg.Secret))
}

// ParseJWT 解析 JWT Token
func ParseJWT(tokenString string) (*Claims, error) {
	cfg := config.GetConfig().JWT

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshJWT 刷新 Token (如果距离过期时间小于 2 小时，则刷新)
func RefreshJWT(tokenString string) (string, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return "", err
	}

	// 如果距离过期时间大于 2 小时，则不刷新
	if time.Until(claims.ExpiresAt.Time) > 2*time.Hour {
		return tokenString, nil
	}

	// 生成新 Token
	return GenerateJWT(claims.UserID, claims.Username)
}
