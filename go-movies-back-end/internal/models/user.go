package models

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// List of valid roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User struct

type User struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Password  string  `json:"password,omitempty"`
	Approved  bool    `json:"approved"`
	Roles     []*Role `json:"roles"`
	// Timestamps
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Role struct {
	ID     int    `json:"-"`
	UserID int    `json:"-"`
	Role   string `json:"role"`
	// Timestamps
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type JwtUser struct {
	ID        int      `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
}

// Function to compare a plain-text password with a hashed password

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			//invalid password
			fmt.Printf("Mismatch-%s", err)
			return false, nil
		default:
			fmt.Println(err)
			return false, err
		}

	}
	return true, nil

}

// Function to hash a password using bcrypt
func (u *User) HashPassword() error {
	// Hash the password using bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}
	u.Password = string(hashedBytes)
	return nil
}

// Validation for User object
func (u *User) Validate() error {
	if u.FirstName == "" {
		return errors.New("first name cannot be empty")
	}
	if u.LastName == "" {
		return errors.New("last name cannot be empty")
	}
	if u.Email == "" {
		return errors.New("email cannot be empty")
	}
	if u.Password == "" {
		return errors.New("password cannot be empty")
	}
	// Add password length validation
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	// email validation
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return fmt.Errorf("invalid email format: %v", err)
	}
	// Make sure roles are valid
	for _, role := range u.Roles {
		if role.Role != RoleAdmin && role.Role != RoleUser {
			return errors.New("invalid role")
		}
	}

	return nil
}

// Function to verify if user has given role
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r.Role == role {
			return true
		}
	}
	return false
}

// Get list of roles from user object a list of string
func (u *User) GetRoles() []string {
	roles := make([]string, len(u.Roles))
	for i, role := range u.Roles {
		roles[i] = role.Role
	}
	return roles
}

// Function to create JWT User
func (u *User) ToJwtUser() *JwtUser {
	return &JwtUser{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Roles:     u.GetRoles(),
	}
}
