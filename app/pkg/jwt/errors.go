package jwt

import "errors"

var (
	ErrInvalidToken         = errors.New("jwt: token is invalid")
	ErrParseTokenClaims     = errors.New("jwt: failed to parse token claims")
	ErrMissingTokenClaims   = errors.New("jwt: failed to get token claims")
	ErrExpiredToken         = errors.New("jwt: token is expired")
	ErrInvalidClaims        = errors.New("jwt: user claims are invalid")
	ErrInvalidSigningMethod = errors.New("jwt: invalid signing method")
	ErrInvalidSigningKey    = errors.New("jwt: invalid signing key")
	ErrFailedDecode         = errors.New("jwt: failed to decode token")
	ErrInvalidIssuerClaim   = errors.New("jwt: invalid issuer claim")
	ErrInvalidIssuedAtClaim = errors.New("jwt: IssuedAt claim is not valid")
	ErrInvalidUserIdClaim   = errors.New("jwt: UserID claim is not valid")
)
