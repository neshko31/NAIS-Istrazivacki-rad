package models

import (
	"errors"
	"regexp"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserRegister struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	PasswordC string `json:"password_c" binding:"required,min=6"`
	Email     string `json:"email" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=tourist guide TOURIST GUIDE"`
}

func (u *UserRegister) Validate() error {
	if u.Password != u.PasswordC {
		return errors.New("Password mismatch")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("Invalid email format")
	}

	if u.Role == "" {
		return errors.New("Invalid role")
	}

	return nil
}
