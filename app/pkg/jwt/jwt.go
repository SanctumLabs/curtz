package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtGen interface {
	Encode(uid string, signingKey string, issuer string, expireDelta int) (string, error)
	EncodeRefreshToken(uid string, signingKey string, issuer string, expireDelta int) (string, error)
	Decode(tokenString string, issuer string, signingKey string) (string, time.Time, error)
}

type jwtGen struct{}

// Claims is our custom metadata, which will be hashed and sent as the second segment in our JWT
type Claims struct {
	jwt.StandardClaims
	UserId string `json:"id"`
}

func New() JwtGen {
	return new(jwtGen)
}

// Encode a claim into a JWT token
func (j jwtGen) Encode(uid string, signingKey string, issuer string, expireDelta int) (string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(expireDelta)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
			Subject:   uid,
		},
		UserId: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

// EncodeRefreshToken encodes claims into a refresh token
func (j jwtGen) EncodeRefreshToken(uid string, signingKey string, issuer string, expireDelta int) (string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(expireDelta)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
			Subject:   uid,
		},
		UserId: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

// Decode a jwt token and returns user id if valid
func (j jwtGen) Decode(tokenString string, issuer string, signingKey string) (string, time.Time, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(signingKey), nil
	})

	var uid string

	if err != nil {
		return uid, time.Time{}, ErrParseTokenClaims
	}

	if !token.Valid {
		return uid, time.Time{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return uid, time.Time{}, ErrMissingTokenClaims
	}

	if claims.UserId == "" {
		return uid, time.Time{}, ErrInvalidUserIdClaim
	}

	if claims.IssuedAt == 0 {
		return uid, time.Time{}, ErrInvalidIssuedAtClaim
	}

	if claims.Issuer == "" || claims.Issuer != issuer {
		return uid, time.Time{}, ErrInvalidIssuerClaim
	}

	return claims.UserId, time.Unix(claims.IssuedAt, 0), nil
}
