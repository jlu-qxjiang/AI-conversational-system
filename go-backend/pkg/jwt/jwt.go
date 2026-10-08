package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const CtxtUserIDKey = "userID"
const CtxtUsernameKey = "username"
const CtxtRoleKey = "role"
const (
	TokenExpireDuration = time.Hour * 24 * 365 // JWT 过期时间：365天
)

// 密钥（从环境变量读取，避免明文写入代码）
var mySecret = []byte(os.Getenv("JWT_SECRET"))

// CustomClaims 自定义声明
type CustomClaims struct {
	UserID   int64  `json:"userID"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string, role int) (string, error) {
	claims := CustomClaims{
		Username: username,
		UserID:   userID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "bluebell",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySecret)
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
