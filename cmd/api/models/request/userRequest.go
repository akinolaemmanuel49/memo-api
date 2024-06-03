package request

import (
	"errors"

	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"github.com/akinolaemmanuel49/memo-api/internal/helpers"
)

type User struct {
	Username  *string `json:"username" validate:"omitempty,min=5,excludes= "`
	Email     *string `json:"email" validate:"omitempty,email"`
	FirstName *string `json:"firstName" validate:"omitempty,min=2,max=100,excludes= "`
	LastName  *string `json:"lastName" validate:"omitempty,min=2,max=100,excludes= "`
	Password  *string `json:"password" validate:"omitempty,min=6,excludes= "`
	Status    *string `json:"status" validate:"omitempty"`
	About     *string `json:"about" validate:"omitempty"`
}

const (
	UserFieldUsername = iota
	UserFieldEmail
	UserFieldFirstName
	UserFieldLastName
	UserFieldPassword
)

func (u User) ToModel() models.User {
	return models.User{
		Username:  helpers.SafeDereference(u.Username),
		Email:     helpers.SafeDereference(u.Email),
		FirstName: helpers.SafeDereference(u.FirstName),
		LastName:  helpers.SafeDereference(u.LastName),
		Password:  helpers.SafeDereference(u.Password),
		Status:    helpers.SafeDereference(u.Status),
		About:     helpers.SafeDereference(u.About),
	}
}

// ValidateRequired verifies that the required fields for the request are provided.
func (u User) ValidateRequired(required ...int) error {
	for _, field := range required {
		switch field {
		case UserFieldUsername:
			if u.Username == nil {
				return errors.New("username is required")
			}

		case UserFieldEmail:
			if u.Email == nil {
				return errors.New("email is required")
			}

		case UserFieldFirstName:
			if u.FirstName == nil {
				return errors.New("first name is required")
			}

		case UserFieldLastName:
			if u.LastName == nil {
				return errors.New("last name is required")
			}

		case UserFieldPassword:
			if u.Password == nil {
				return errors.New("password is required")
			}
		}
	}

	return nil
}
