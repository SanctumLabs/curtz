package auth

import (
	"errors"
	"time"

	"github.com/sanctumlabs/curtz/app/config"
	"github.com/sanctumlabs/curtz/app/pkg/jwt"
)

// AuthService represents underlying authentication implementation
type AuthService struct {
	config config.AuthConfig
	jwt    jwt.JwtGen
}

func NewService(config config.AuthConfig, jwtGen jwt.JwtGen) *AuthService {
	return &AuthService{config, jwtGen}
}

func (svc *AuthService) Authenticate(tokenString string) (string, time.Time, error) {
	return svc.jwt.Decode(tokenString, svc.config.Jwt.Issuer, svc.config.Jwt.Secret)
}

func (svc *AuthService) GenerateToken(userId string) (string, error) {
	if token, err := svc.jwt.Encode(userId, svc.config.Jwt.Secret, svc.config.Jwt.Issuer, svc.config.ExpireDelta); err != nil {
		return "", errors.New("failed to create access token")
	} else {
		return token, nil
	}
}

func (svc *AuthService) GenerateRefreshToken(userId string) (string, error) {
	token, err := svc.jwt.EncodeRefreshToken(userId, svc.config.Jwt.Secret, svc.config.Jwt.Issuer, svc.config.RefreshExpireDelta)
	if err != nil {
		return "", errors.New("failed to create refresh token")
	}
	return token, nil
}
