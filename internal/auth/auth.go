package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CourtlyClaims struct {
    Sub string `json:"sub"`
    Role string `json:"role"`
    jwt.RegisteredClaims
}

type AuthService interface {
	GenerateToken(companyID string) (string, error)
	GetSecretKey() []byte
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
	claims := CourtlyClaims{
        Sub: companyID,
        Role: "authenticated",
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    "courtly-api",
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)),
            Audience: jwt.ClaimStrings{"authenticated"},
        },
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authServiceImpl) GetSecretKey() []byte {
	return s.jwtSecret
}
