package entity

import "golang.org/x/crypto/bcrypt"

// Password contains user password information
type Password struct {
	Value string
}

func NewPassword(value string) Password {
	return Password{
		Value: value,
	}
}

//HashPassword hashes the user password using bcrypt hash function
func (p *Password) HashPassword() error {

	pwd := []byte(p.Value)

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.Value = string(hash)

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
