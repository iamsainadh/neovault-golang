package model

import (
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	Mfa_enabled   string    `json:"mfa_enabled"`
	Is_active     bool      `json:"is_active"`
	Is_verified   bool      `json:"is_verified"`
	Deleted_at    time.Time `json:"deleted_at"`
	Last_login_at time.Time `json:"Last_login_at"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	Name     string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	UserId string
	Token  string
}

type LoginUserResponse struct {
	UserId   string
	Name     string
	Email    string
	Password string
	Token    string
}

type CreatedUser struct {
	ID    string
	Name  string
	Email string
}
type LoginUserInput struct {
	Email    string
	Password string
}
