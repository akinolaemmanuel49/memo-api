package repository

import "github.com/akinolaemmanuel49/memo-api/domain/models"

type UserRepository interface {
	Create(user *models.User) (models.User, error)
	GetById(id string) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetFollowersOfUser(id string, page, pageSize int) ([]models.User, error)
	GetUsersFollowedBy(id string, page, pageSize int) ([]models.User, error)
	Update(id string, updatedUser models.User) (models.User, error)
	Delete(id string, deletedUser models.User) (models.User, error)
}
