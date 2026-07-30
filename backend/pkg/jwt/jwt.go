package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 载荷
type Claims struct {
	UserID            string `json:"user_id"`
	Username          string `json:"username"`
	Nickname          string `json:"nickname"`
	Role              string `json:"role"`
	Department        string `json:"department"`
	CanToggleAdminView bool  `json:"can_toggle_admin_view"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(secret string, expiry int, userID, username, nickname, role, department string, canToggleAdminView bool) (string, error) {
	if expiry <= 0 {
		expiry = 86400
	}
	claims := Claims{
		UserID:            userID,
		Username:          username,
		Nickname:          nickname,
		Role:              role,
		Department:        department,
		CanToggleAdminView: canToggleAdminView,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "osm-mngt",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析 JWT token
func ParseToken(secret string, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
