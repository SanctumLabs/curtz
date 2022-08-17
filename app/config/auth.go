package config

type Jwt struct {
	// Secret key used to sign the JWT token
	Secret string
	// ExpireDelta is the expiration time of the JWT Access token
	ExpireDelta int
	// RefreshExpirationDelta is the expiry time of the JWT Refresh token
	RefreshExpireDelta int
	// Issuer of the JWT token
	Issuer string
}

type AuthConfig struct {
	Jwt
}
