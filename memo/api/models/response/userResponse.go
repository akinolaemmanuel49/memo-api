package response

import (
	"time"

	"github.com/akinolaemmanuel49/memo-api/domain/models"
)

type User struct {
	ID             string    `json:"id,omitempty"`
	Username       string    `json:"username,omitempty"`
	Email          string    `json:"email,omitempty"`
	FirstName      string    `json:"firstName,omitempty"`
	LastName       string    `json:"lastName,omitempty"`
	AvatarURL      string    `json:"avatarURL,omitempty"`
	Status         string    `json:"status"`
	About          string    `json:"about"`
	Deleted        bool      `json:"deleted"`
	FollowerCount  int64     `json:"followerCount"`
	FollowingCount int64     `json:"followingCount"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}

func UserResponseFromModel(user models.User) User {
	return User{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AvatarURL:      user.AvatarURL,
		Status:         user.Status,
		About:          user.About,
		Deleted:        user.Deleted,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func MultipleUserResponseFromModel(users []models.User) []User {
	var userResponses []User
	for _, user := range users {
		userResponse := UserResponseFromModel(user)
		userResponses = append(userResponses, userResponse)
	}
	return userResponses
}
