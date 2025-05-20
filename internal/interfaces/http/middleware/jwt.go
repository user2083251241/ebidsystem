package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/user2083251241/ebidsystem/internal/app/config"
)

// JWTClaims 自定义JWT声明
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, role string) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour) // 令牌有效期为3天
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Get("JWT_SECRET")))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid() {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKeyUsage
}
