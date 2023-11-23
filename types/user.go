package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

const (
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 7
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"isAdmin" json:"isAdmin"`
}

type NewUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func NewUserFromParams(params NewUserParams) (*User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encryptedPassword),
	}, nil
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p *UpdateUserParams) BSON() bson.M {

	b := bson.M{}
	if len(p.FirstName) > 0 {
		b["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		b["lastName"] = p.LastName
	}
	return b
}
func (params UpdateUserParams) Validate() []error {
	var errors []error

	if len(params.FirstName) < minFirstNameLength {
		errors = append(errors, fmt.Errorf("first name length should be at least %d characters", minFirstNameLength))
	}

	if len(params.LastName) < minLastNameLength {
		errors = append(errors, fmt.Errorf("last name length should be at least %d characters", minLastNameLength))
	}
	return errors
}
func (params NewUserParams) Validate() []error {
	var errors []error

	if len(params.FirstName) < minFirstNameLength {
		errors = append(errors, fmt.Errorf("first name length should be at least %d characters", minFirstNameLength))
	}

	if len(params.LastName) < minLastNameLength {
		errors = append(errors, fmt.Errorf("last name length should be at least %d characters", minLastNameLength))
	}

	if len(params.Password) < minPasswordLength {
		errors = append(errors, fmt.Errorf("password length should be at least %d characters", minPasswordLength))
	}

	if !isEmailValid(params.Email) {
		errors = append(errors, fmt.Errorf("email is invalid"))

	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func IsValidPassword(encryptedPassword, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}
