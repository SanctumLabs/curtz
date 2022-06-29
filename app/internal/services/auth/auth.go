package auth

import (
	"time"

	"github.com/sanctumlabs/curtz/app/config"
	"github.com/sanctumlabs/curtz/app/pkg/jwt"
)

type AuthService struct {
	config config.AuthConfig
}

func NewService(config config.AuthConfig) *AuthService {
	return &AuthService{config}
}

func (svc *AuthService) Authenticate(tokenString string) (string, time.Time, error) {
	return jwt.Decode(tokenString, svc.config.Jwt.Issuer, svc.config.Jwt.Secret)
}

func (svc *AuthService) GenerateToken(userId string) (string, error) {
	return jwt.Encode(userId, svc.config.Jwt.Secret, svc.config.Jwt.Issuer, svc.config.ExpireDelta)
}
