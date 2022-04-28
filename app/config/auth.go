package config

type AuthConfig struct {
	Jwt struct {
		// Secret key used to sign the JWT token
		Secret string
		// Expiration time of the JWT token
		ExpireDelta int
		Issuer      string
	}
}
