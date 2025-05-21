package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthService interface {
	GenerateToken(companyID string) (string, error)
}

type authServiceImpl struct {
	jwtSecret []byte
}

func NewAuthService(jwtSecret []byte) AuthService {
	return &authServiceImpl{
		jwtSecret: jwtSecret,
	}
}

func (s *authServiceImpl) GenerateToken(companyID string) (string, error) {
	claims := jwt.MapClaims{
		"company_id": companyID,
		"iss":     "courtly-api",
		"exp":        time.Now().Add(24 * time.Hour * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.jwtSecret)
}
