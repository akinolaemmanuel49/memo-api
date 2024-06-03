package handlers

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/akinolaemmanuel49/memo-api/cmd/api/helpers"
	"github.com/akinolaemmanuel49/memo-api/cmd/api/internal"
	"github.com/akinolaemmanuel49/memo-api/cmd/api/models/response"
	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"github.com/akinolaemmanuel49/memo-api/domain/repository"
)

type UserHandler interface {
	Get(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetFollowers(ctx *gin.Context)
	GetFollowing(ctx *gin.Context)
	DeleteAvatar(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type userHandler struct {
	app internal.Application
}

func NewUserHandler(app internal.Application) UserHandler {
	return userHandler{app: app}
}

// Get returns the profile of an authenticated user.
func (uh userHandler) Get(ctx *gin.Context) {
	// fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)
	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// return fetched user
	ctx.JSON(
		http.StatusOK,
		response.UserResponseFromModel(user),
	)

}

// GetAll retrieves a list of users.
func (uh userHandler) GetAll(ctx *gin.Context) {
	// fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)
	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// retrieve query params for pagination
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("pageSize")

	if pageStr == "" {
		pageStr = helpers.DefaultPage
	}
	if pageSizeStr == "" {
		pageSizeStr = helpers.DefaultPageSize
	}

	// convert query strings to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		switch {
		default:
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid value for page parameter"))
		}
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		switch {
		default:
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid value for pageSize parameter"))
		}
		return
	}

	// retrieve list of users from database
	followers, err := uh.app.Repositories.Users.GetAll(page, pageSize)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return fetched users
	ctx.JSON(
		http.StatusOK,
		response.MultipleUserResponseFromModel(followers),
	)
}

// Update updates and returns the profile of an authenticated user.
// Updatable details include: username, first name, last name, storage, status, about.
func (uh userHandler) Update(ctx *gin.Context) {
	// Fetch authenticated user from context
	user := helpers.ContextGetUser(ctx)

	// Return an authentication error if users is not in context
	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// Validate request data
	username := ctx.PostForm("username")
	email := ctx.PostForm("email")
	firstName := ctx.PostForm("firstName")
	lastName := ctx.PostForm("lastName")
	password := ctx.PostForm("password")
	status := ctx.PostForm("status")
	about := ctx.PostForm("about")

	if password != "" {
		if err := helpers.HashPassword(&password); err != nil {
			helpers.HandleInternalServerError(ctx, err)
			return
		}
	}

	// Process the file provided with the form as the new storage of the user
	var avatarURL string
	avatarFile, _, err := ctx.Request.FormFile("avatarFile")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			// When no file is provided with the form,
			// simply update avatarURL, retaining the user's current storage
			avatarURL = user.AvatarURL
		default:
			helpers.HandleInternalServerError(ctx, err)
			return
		}
	} else {
		// When a file is provided with the form, attempt to upload it as the user's new storage
		// and update the avatarURL to be saved as part of the user profile
		avatarURL, err = uh.app.Repositories.File.UploadAvatar(user.ID, avatarFile)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrUnapprovedFileType):
				helpers.HandleValidationError(ctx, err)
			default:
				helpers.HandleInternalServerError(ctx, err)
			}
			return
		}
	}

	updatedUser := user

	// Update user profile
	if username != "" {
		updatedUser.Username = username
	}
	if email != "" {
		updatedUser.Email = email
	}
	if firstName != "" {
		updatedUser.FirstName = firstName
	}
	if lastName != "" {
		updatedUser.LastName = lastName
	}
	if password != "" {
		updatedUser.Password = password
	}
	if status != "" {
		updatedUser.Status = status
	}
	if about != "" {
		updatedUser.About = about
	}

	updatedUser.AvatarURL = avatarURL

	if _, err := uh.app.Repositories.Users.Update(user.ID, updatedUser); err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, err)
		case errors.Is(err, repository.ErrDuplicateDetails):
			helpers.HandleErrorResponse(ctx, http.StatusConflict, errors.New("username already exists"))
		case errors.Is(err, repository.ErrRecordDeleted):
			data := response.UserResponseFromModel(user)
			helpers.HandleLogicalDeleteError(ctx, data, err)
		case errors.Is(err, repository.ErrConcurrentUpdate):
			helpers.HandleErrorResponse(ctx, http.StatusConflict, err)
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "User was successfully updated",
		},
	)
}

// GetFollowers retrieves a list of users following the authenticated user.
func (uh userHandler) GetFollowers(ctx *gin.Context) {
	// fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)
	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// retrieve query params for pagination
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("pageSize")

	if pageStr == "" {
		pageStr = helpers.DefaultPage
	}
	if pageSizeStr == "" {
		pageSizeStr = helpers.DefaultPageSize
	}

	// convert query strings to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		switch {
		default:
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid value for page parameter"))
		}
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		switch {
		default:
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid value for pageSize parameter"))
		}
		return
	}

	// retrieve list of followers from database
	followers, err := uh.app.Repositories.Users.GetFollowersOfUser(user.ID, page, pageSize)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return fetched users
	ctx.JSON(
		http.StatusOK,
		response.MultipleUserResponseFromModel(followers),
	)
}

// GetFollowing retrieves a list of users being followed by the authenticated user.
func (uh userHandler) GetFollowing(ctx *gin.Context) {
	// fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)
	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// retrieve query params for pagination
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("pageSize")

	if pageStr == "" {
		pageStr = helpers.DefaultPage
	}
	if pageSizeStr == "" {
		pageSizeStr = helpers.DefaultPageSize
	}

	// convert query strings to integers
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		switch {
		default:
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid value for page parameter"))
		}
		return
	}
	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil {
		switch {
		default:
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid value for pageSize parameter"))
		}
		return
	}

	// retrieve list of users being followed by authenticated user from database
	followers, err := uh.app.Repositories.Users.GetUsersFollowedBy(user.ID, page, pageSize)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return fetched users
	ctx.JSON(
		http.StatusOK,
		response.MultipleUserResponseFromModel(followers),
	)
}

// Delete performs a soft delete of a user instance.
func (uh userHandler) Delete(ctx *gin.Context) {
	user := helpers.ContextGetUser(ctx)
	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	if _, err := uh.app.Repositories.Users.Delete(user.ID, user); err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, err)
		case errors.Is(err, repository.ErrRecordDeleted):
			data := response.UserResponseFromModel(user)
			helpers.HandleLogicalDeleteError(ctx, data, err)
		case errors.Is(err, repository.ErrConcurrentUpdate):
			helpers.HandleErrorResponse(ctx, http.StatusConflict, err)
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "User was successfully deleted",
		},
	)
}

// DeleteAvatar deletes the storage of a user.
func (uh userHandler) DeleteAvatar(ctx *gin.Context) {
	user := helpers.ContextGetUser(ctx)
	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	err := uh.app.Repositories.File.DeleteAvatar(user.ID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	// reset avatarURL of user in repository
	user.AvatarURL = ""
	if _, err := uh.app.Repositories.Users.Update(user.ID, user); err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, err)
		case errors.Is(err, repository.ErrConcurrentUpdate):
			helpers.HandleErrorResponse(ctx, http.StatusConflict, err)
		case errors.Is(err, repository.ErrRecordDeleted):
			data := response.UserResponseFromModel(user)
			helpers.HandleLogicalDeleteError(ctx, data, err)
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Avatar was successfully deleted",
		},
	)
}
