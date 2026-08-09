package identity

// CreateUserRequest represents the request payload for creating a new user
type CreateUserRequest struct {
	Username     string         `json:"username"`
	FullName     UserFullName   `json:"full_name"`
	Email        Email          `json:"email"`
	PasswordHash string         `json:"password_hash"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

// UpdateUserRequest represents the request payload for updating an existing user
type UpdateUserRequest struct {
	ID        string  `json:"id"`
	Username  *string `json:"username,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
}

type UpdateUserVerificationRequest struct {
	ID        string  `json:"id"`
	Username  *string `json:"username,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
}
