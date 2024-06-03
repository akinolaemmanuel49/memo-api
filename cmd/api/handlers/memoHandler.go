package handlers

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/akinolaemmanuel49/memo-api/memo/api/models/response"

	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"github.com/akinolaemmanuel49/memo-api/domain/repository"
	"github.com/akinolaemmanuel49/memo-api/memo/api/helpers"
	"github.com/akinolaemmanuel49/memo-api/memo/api/internal"
	"github.com/akinolaemmanuel49/memo-api/memo/api/models/request"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MemoHandler interface {
	CreateTextMemo(ctx *gin.Context)
	CreateImageMemo(ctx *gin.Context)
	CreateVideoMemo(ctx *gin.Context)
	CreateAudioMemo(ctx *gin.Context)
	GetMemo(ctx *gin.Context)
	DeleteMemo(ctx *gin.Context)
	LikeMemo(ctx *gin.Context)
	UnlikeMemo(ctx *gin.Context)
	ShareMemo(ctx *gin.Context)
	UnshareMemo(ctx *gin.Context)
	GetAllMemos(ctx *gin.Context)
	GetSubscribedMemos(ctx *gin.Context)
	GetMemosByOwnerID(ctx *gin.Context)
	GetOwnMemos(ctx *gin.Context)
}

type memoHandler struct {
	app internal.Application
}

func NewMemoHandler(app internal.Application) MemoHandler {
	return memoHandler{app: app}
}

// CreateTextMemo creates a new instance of a text based memo.
func (mh memoHandler) CreateTextMemo(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	requestBody := request.TextMemo{}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := requestBody.ValidateRequired(
		request.TextMemoFieldContent)

	if err != nil {
		helpers.HandleValidationError(ctx, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		helpers.HandleValidationError(ctx, err)
		return
	}

	// convert request to text memo model
	textMemo := requestBody.ToModel()
	textMemo.OwnerID = user.ID
	textMemo.MemoType = "text"

	// attempt to save text memo in repository
	newTextMemo, err := mh.app.Repositories.Memo.CreateMemo(user.ID, &textMemo)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return newly created text memo
	ctx.JSON(
		http.StatusCreated,
		response.MemoResponseFromModel(newTextMemo),
	)
}

// CreateImageMemo creates a new instance of an image based memo.
func (mh memoHandler) CreateImageMemo(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// Validate request data
	caption := ctx.PostForm("caption")

	imageMemo := models.Memo{
		OwnerID:  user.ID,
		MemoType: "image",
		Caption:  caption,
	}

	// attempt to save image memo in repository
	newImageMemo, err := mh.app.Repositories.Memo.CreateMemo(user.ID, &imageMemo)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// Process the file provided with the form as the new storage of the memo
	var memoURL string
	memoFile, _, err := ctx.Request.FormFile("memoFile")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			// When no file is provided with the form set memoURL to ""
			memoURL = ""
		default:
			helpers.HandleInternalServerError(ctx, err)
			return
		}
	} else {
		// When a file is provided with the form, attempt to upload it as the memo's new storage
		// and update the memoURL to be saved as part of the memo
		memoURL, err = mh.app.Repositories.File.UploadMemoMedia(newImageMemo.ID, memoFile, newImageMemo.MemoType)
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

	newImageMemo.Content = memoURL

	updatedMemo, err := mh.app.Repositories.Memo.Update(newImageMemo.ID, newImageMemo)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	data := response.MemoResponseFromModel(updatedMemo)

	// return newly create image memo
	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"data":    data,
			"message": `Image memo was successfully created.`,
		},
	)
}

// CreateVideoMemo creates a new instance of a video based memo.
func (mh memoHandler) CreateVideoMemo(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// Validate request data
	caption := ctx.PostForm("caption")

	videoMemo := models.Memo{
		OwnerID:  user.ID,
		MemoType: "video",
		Caption:  caption,
	}

	// attempt to save video memo in repository
	newVideoMemo, err := mh.app.Repositories.Memo.CreateMemo(user.ID, &videoMemo)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// Process the file provided with the form as the new storage of the memo
	var memoURL string
	memoFile, _, err := ctx.Request.FormFile("memoFile")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			// When no file is provided with the form set memoURL to ""
			memoURL = ""
		default:
			helpers.HandleInternalServerError(ctx, err)
			return
		}
	} else {
		// When a file is provided with the form, attempt to upload it as the memo's new storage
		// and update the memoURL to be saved as part of the memo
		memoURL, err = mh.app.Repositories.File.UploadMemoMedia(newVideoMemo.ID, memoFile, newVideoMemo.MemoType)
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

	newVideoMemo.Content = memoURL

	updatedMemo, err := mh.app.Repositories.Memo.Update(newVideoMemo.ID, newVideoMemo)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	data := response.MemoResponseFromModel(updatedMemo)

	// return newly create video memo
	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"data":    data,
			"message": `Video memo was successfully created.`,
		},
	)
}

// CreateAudioMemo creates a new instance of an audio based memo.
func (mh memoHandler) CreateAudioMemo(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// Validate request data
	caption := ctx.PostForm("caption")

	audioMemo := models.Memo{
		OwnerID:  user.ID,
		MemoType: "audio",
		Caption:  caption,
	}

	// attempt to save audio memo in repository
	newAudioMemo, err := mh.app.Repositories.Memo.CreateMemo(user.ID, &audioMemo)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// Process the file provided with the form as the new storage of the memo
	var memoURL string
	memoFile, _, err := ctx.Request.FormFile("memoFile")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			// When no file is provided with the form set memoURL to ""
			memoURL = ""
		default:
			helpers.HandleInternalServerError(ctx, err)
			return
		}
	} else {
		// When a file is provided with the form, attempt to upload it as the memo's new storage
		// and update the memoURL to be saved as part of the memo
		memoURL, err = mh.app.Repositories.File.UploadMemoMedia(newAudioMemo.ID, memoFile, newAudioMemo.MemoType)
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

	newAudioMemo.Content = memoURL

	updatedMemo, err := mh.app.Repositories.Memo.Update(newAudioMemo.ID, newAudioMemo)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	data := response.MemoResponseFromModel(updatedMemo)

	// return newly create audio memo
	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"data":    data,
			"message": `Audio memo was successfully created.`,
		},
	)
}

// GetMemo fetches an instance of a text based memo that matches a query.
func (mh memoHandler) GetMemo(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	memoID := ctx.Param("memoID")
	if memoID == "" {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, repository.ErrMemoIDQueryMissing)
	}

	// Fetch memo data
	memo, err := mh.app.Repositories.Memo.GetMemo(memoID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, err)
		case errors.Is(err, repository.ErrRecordDeleted):
			data := response.MemoResponseFromModel(memo)
			helpers.HandleLogicalDeleteError(ctx, data, err)
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return memo
	ctx.JSON(
		http.StatusOK,
		response.MemoResponseFromModel(memo),
	)
}

// LikeMemo adds a new like instance to the likes table.
func (mh memoHandler) LikeMemo(ctx *gin.Context) {
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	memoID := ctx.Param("memoID")
	if memoID == "" {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, repository.ErrMemoIDQueryMissing)
	}

	_, err := mh.app.Repositories.Memo.LikeMemo(user.ID, memoID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, err)
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return success response
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Like was successful.",
		},
	)
}

// UnlikeMemo deletes a like instance from the likes table.
func (mh memoHandler) UnlikeMemo(ctx *gin.Context) {
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	memoID := ctx.Param("memoID")
	if memoID == "" {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, repository.ErrMemoIDQueryMissing)
	}

	err := mh.app.Repositories.Memo.UnlikeMemo(user.ID, memoID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, err)
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return success response
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Unlike was successful.",
		},
	)
}

// ShareMemo adds a new share instance to the shares table.
func (mh memoHandler) ShareMemo(ctx *gin.Context) {
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	memoID := ctx.Param("memoID")
	if memoID == "" {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, repository.ErrMemoIDQueryMissing)
	}

	_, err := mh.app.Repositories.Memo.ShareMemo(user.ID, memoID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, err)
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return success response
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Share was successful.",
		},
	)
}

// UnshareMemo deletes a share instance from the shares table.
func (mh memoHandler) UnshareMemo(ctx *gin.Context) {
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	memoID := ctx.Param("memoID")
	if memoID == "" {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, repository.ErrMemoIDQueryMissing)
	}

	err := mh.app.Repositories.Memo.UnshareMemo(user.ID, memoID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, err)
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return success response
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Unshare was successful.",
		},
	)
}

func (mh memoHandler) GetAllMemos(ctx *gin.Context) {
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
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		}
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		switch {
		default:
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		}
		return
	}

	// retrieve list of memos by followed users from the database
	memos, err := mh.app.Repositories.Memo.GetAllMemos(page, pageSize)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return fetched memos
	ctx.JSON(
		http.StatusOK,
		response.MultipleMemoResponseFromModel(memos))
}

func (mh memoHandler) GetSubscribedMemos(ctx *gin.Context) {
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
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		}
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		switch {
		default:
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		}
		return
	}

	// retrieve list of memos by followed users from the database
	memos, err := mh.app.Repositories.Memo.GetMemosByFollowing(user.ID, page, pageSize)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return fetched memos
	ctx.JSON(
		http.StatusOK,
		response.MultipleMemoResponseFromModel(memos))
}

// DeleteMemo creates a new instance of an image based memo.
func (mh memoHandler) DeleteMemo(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	memoID := ctx.Param("memoID")
	if memoID == "" {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, repository.ErrMemoIDQueryMissing)
	}

	memo, err := mh.app.Repositories.Memo.GetMemo(memoID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, err)
		case errors.Is(err, repository.ErrRecordDeleted):
			data := response.MemoResponseFromModel(memo)
			helpers.HandleLogicalDeleteError(ctx, data, err)
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}
	if memo.MemoType != "text" {
		err = mh.app.Repositories.File.DeleteMemoMedia(memoID)
		if err != nil {
			helpers.HandleInternalServerError(ctx, err)
			return
		}
	}

	_, err = mh.app.Repositories.Memo.Delete(memoID, memo)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	// return newly create text memo
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Memo was successfully deleted",
		},
	)
}

// GetMemosByOwnerID fetches all memos owned by a user with matching ID.
func (mh memoHandler) GetMemosByOwnerID(ctx *gin.Context) {
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// retrieve query params for pagination
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("pageSize")
	ownerID := ctx.Param("ownerID")

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
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		}
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		switch {
		default:
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		}
		return
	}

	// retrieve list of memos by followed users from the database
	memos, err := mh.app.Repositories.Memo.GetMemosByOwnerID(ownerID, page, pageSize)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return fetched memos
	ctx.JSON(
		http.StatusOK,
		response.MultipleMemoResponseFromModel(memos))
}

// GetOwnMemos fetches all memos owned by authenticated user.
func (mh memoHandler) GetOwnMemos(ctx *gin.Context) {
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
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		}
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		switch {
		default:
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		}
		return
	}

	// retrieve list of memos by followed users from the database
	memos, err := mh.app.Repositories.Memo.GetMemosByOwnerID(user.ID, page, pageSize)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return fetched memos
	ctx.JSON(
		http.StatusOK,
		response.MultipleMemoResponseFromModel(memos))
}
