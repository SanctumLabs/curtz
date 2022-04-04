package user

import "golang.org/x/crypto/bcrypt"

// Password contains user password information
type Password struct {
	Password string `db:"password"`
}

func NewPassword(value string) Password {
	return Password{
		Password: value,
	}
}

//HashPassword hashes the user password using bcrypt hash function
func (p *Password) HashPassword() error {

	pwd := []byte(p.Password)

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.Password = string(hash)

	return nil
}

//Compare compares the password hash against the passed in password string
func (p Password) Compare(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
