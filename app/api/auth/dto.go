package auth

type signUpRequestDto struct {
	email    string `json:"username" binding:"required"`
	password string `json:"password" binding:"required"`
}

type loginRequestDto struct {
	email    string `json:"username" binding:"required"`
	password string `json:"password" binding:"required"`
}
