package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/akinolaemmanuel49/memo-api/domain/repository"
	"github.com/akinolaemmanuel49/memo-api/memo/api/helpers"
	"github.com/akinolaemmanuel49/memo-api/memo/api/internal"
	"github.com/akinolaemmanuel49/memo-api/memo/api/models/request"
	"github.com/akinolaemmanuel49/memo-api/memo/api/models/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	SignUp(ctx *gin.Context)
	Token(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type authHandler struct {
	app internal.Application
}

func NewAuthHandler(app internal.Application) AuthHandler {
	return authHandler{app: app}
}

// SignUp registers a user and adds their details to the repository.
func (a authHandler) SignUp(ctx *gin.Context) {
	// validate request structure and contents
	requestBody := request.User{}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := requestBody.ValidateRequired(
		request.UserFieldUsername,
		request.UserFieldEmail,
		request.UserFieldFirstName,
		request.UserFieldLastName,
		request.UserFieldPassword)

	if err != nil {
		helpers.HandleValidationError(ctx, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		helpers.HandleValidationError(ctx, err)
		return
	}

	// convert request to user model and hash password
	user := requestBody.ToModel()
	if err := helpers.HashPassword(&user.Password); err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	// attempt to save user in repository
	newUser, err := a.app.Repositories.Users.Create(&user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateDetails):
			helpers.HandleErrorResponse(ctx, http.StatusConflict, errors.New("username or email already exists"))

		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// generate access and refresh tokens for user
	accessToken, refreshToken, err := helpers.GenerateTokens(a.app.Config.JWTSecret, newUser.ID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	// return newly created user with tokens
	ctx.JSON(
		http.StatusCreated,
		response.AuthResponse{
			Tokens: response.Tokens{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				ExpiresAt:    time.Now().Add(helpers.AccessTokenDuration).Unix(),
			},
			User: response.UserResponseFromModel(newUser),
		},
	)

	/** sample response
	{
		"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI1NzYyNzYsInN1YiI6ImFhMThmYWJmLTUxOTYtNDYyMC1iY2E3LTNiZDBkNGMwN2M4NSJ9.3JfRh-t91MIBPPa9DE7lGYRPyPu2VYZ8qNK55iCxNNE",
		"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwOTQ2NzYsInN1YiI6ImFhMThmYWJmLTUxOTYtNDYyMC1iY2E3LTNiZDBkNGMwN2M4NSJ9.lReNLhFW2zenoNKI7zsDNON-ICOGxxuV52OIh0I5lcM",
		"expiresAt": 1712576276,
		"profile": {
			"id": "aa18fabf-5196-4620-bca7-3bd0d4c07c85",
			"username": "janedoe",
			"email": "janedoe@mail.com",
			"firstName": "Jane",
			"lastName": "Doe",
			"status": "",
			"about": "",
			"deleted": false,
			"createdAt": "2024-04-07T12:37:56.58533+01:00",
			"updatedAt": "2024-04-07T12:37:56.58533+01:00"
		}
	}
	*/
}

// Token authenticates a single user.
func (a authHandler) Token(ctx *gin.Context) {
	// validate request structure and contents
	requestBody := request.User{}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := requestBody.ValidateRequired(
		request.UserFieldEmail,
		request.UserFieldPassword)
	if err != nil {
		helpers.HandleValidationError(ctx, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		helpers.HandleValidationError(ctx, err)
		return
	}

	// fetch user data
	user, err := a.app.Repositories.Users.GetByEmail(*requestBody.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("invalid credentials"))

		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// verify password
	valid, err := helpers.VerifyPassword(user.Password, *requestBody.Password)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}
	if !valid {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	// generate access and refresh tokens for user
	accessToken, refreshToken, err := helpers.GenerateTokens(a.app.Config.JWTSecret, user.ID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	// return logged-in user with tokens
	ctx.JSON(
		http.StatusOK,
		response.AuthResponse{
			Tokens: response.Tokens{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				ExpiresAt:    time.Now().Add(helpers.AccessTokenDuration).Unix(),
			},
			User: response.UserResponseFromModel(user),
		},
	)

	/** sample response
	{
		"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI1NzgyMTUsInN1YiI6ImFhMThmYWJmLTUxOTYtNDYyMC1iY2E3LTNiZDBkNGMwN2M4NSJ9.85-ea92idXNZzUNB1tFQmkVSJIixJo5jWLYvoTL6etc",
		"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwOTY2MTUsInN1YiI6ImFhMThmYWJmLTUxOTYtNDYyMC1iY2E3LTNiZDBkNGMwN2M4NSJ9.iUUjRc6tKEmdOKwRgppOF4lq0Q_DJAQ5MKaAgcw_KcM",
		"expiresAt": 1712578215,
		"profile": {
			"id": "aa18fabf-5196-4620-bca7-3bd0d4c07c85",
			"username": "janedoe",
			"email": "janedoe@mail.com",
			"firstName": "Jane",
			"lastName": "Doe",
			"status": "",
			"about": "",
			"deleted": false,
			"createdAt": "2024-04-07T12:37:56.58533+01:00",
			"updatedAt": "2024-04-07T12:37:56.58533+01:00"
		}
	}
	*/
}

// RefreshToken returns newly generated access and refresh tokens for a user.
func (a authHandler) RefreshToken(ctx *gin.Context) {
	// validate request
	requestBody := struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}{}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		helpers.HandleValidationError(ctx, err)
		return
	}

	// validate refresh token
	claims, err := helpers.ValidateToken(a.app.Config.JWTSecret, requestBody.RefreshToken)
	if err != nil {
		err := errors.New("invalid or expired refresh token")
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, err)
		return
	}

	// generate new access and refresh tokens
	accessToken, refreshToken, err := helpers.GenerateTokens(a.app.Config.JWTSecret, claims.Subject)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	// return new tokens
	ctx.JSON(
		http.StatusOK,
		response.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().Add(helpers.AccessTokenDuration).Unix(),
		},
	)

	/** sample response
	{
		"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI1Nzg4MjMsInN1YiI6ImFhMThmYWJmLTUxOTYtNDYyMC1iY2E3LTNiZDBkNGMwN2M4NSJ9.o_YfU4d05tozie6hXmvQlrFFpbmc1ItaCJa1rHbBwIs",
		"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwOTcyMjMsInN1YiI6ImFhMThmYWJmLTUxOTYtNDYyMC1iY2E3LTNiZDBkNGMwN2M4NSJ9.4O7K1JewYp3EjW8daLsBbbnmNc562oEyTA6FK8zKJSg",
		"expiresAt": 1712578823
	}
	*/
}
