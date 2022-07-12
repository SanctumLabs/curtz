package config

type Jwt struct {
	// Secret key used to sign the JWT token
	Secret string
	// Expiration time of the JWT token
	ExpireDelta int
	// Issure of the JWT token
	Issuer string
}

type AuthConfig struct {
	Jwt
}
