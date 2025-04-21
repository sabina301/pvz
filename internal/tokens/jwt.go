package tokens

import (
	"github.com/golang-jwt/jwt/v5"
	"pvz/configs"
	"pvz/internal/models/auth"
	"time"
)

type Claims struct {
	UserId string    `json:"user_id"`
	Role   auth.Role `json:"auth"`
	jwt.RegisteredClaims
}

func GenerateJwt(userId string, role auth.Role) (string, error) {
	claims := Claims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(configs.AppConfiguration.Auth.Expiration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.AppConfiguration.Auth.JwtSecret))
}

func GenerateDummyJwt(role auth.Role) (string, error) {
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(configs.AppConfiguration.Auth.Expiration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.AppConfiguration.Auth.JwtSecret))
}

func ParseJwt(jwtToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.AppConfiguration.Auth.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*Claims), nil
}
