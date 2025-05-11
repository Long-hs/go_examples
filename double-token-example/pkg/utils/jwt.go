package utils

import (
	"double-token-example/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义JWT声明
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(jti string, userID int64, username string, expireSeconds int64, typ string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expireSeconds))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["typ"] = typ
	return token.SignedString([]byte(config.Cfg.JWT.Secret))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		typ := token.Header["typ"].(string)
		return claims, typ, nil
	}

	return nil, "", jwt.ErrInvalidKey
}

// GetUserIDFromToken 从token中获取用户ID
func GetUserIDFromToken(tokenString string) (int64, string, error) {
	claims, typ, err := ParseToken(tokenString)
	if err != nil {
		return 0, "", err
	}
	return claims.UserID, typ, nil
}
