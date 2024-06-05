package handlers

import (
	"errors"
	"fmt"
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

type MemoHandler interface {
	CreateMemo(ctx *gin.Context)
	GetMemo(ctx *gin.Context)
	DeleteMemo(ctx *gin.Context)
	LikeMemo(ctx *gin.Context)
	UnlikeMemo(ctx *gin.Context)
	ShareMemo(ctx *gin.Context)
	UnshareMemo(ctx *gin.Context)
	GetAllMemos(ctx *gin.Context)
	FindMemos(ctx *gin.Context)
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

// CreateMemo creates a new instance of a memo.
func (mh memoHandler) CreateMemo(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// Parse URL params to get memo type
	memoType := ctx.Param("memoType")

	// Validate request data
	content := ctx.PostForm("content")
	description := ctx.PostForm("description")

	memo := models.Memo{
		OwnerID:     user.ID,
		Type:        memoType,
		Content:     content,
		Description: description,
	}

	// Attempt to save memo in repository
	newMemo, err := mh.app.Repositories.Memo.CreateMemo(user.ID, &memo)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// Process the file provided with the form as the new storage of the memo
	var resourceURL string
	resourceFile, _, err := ctx.Request.FormFile("resourceFile")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			// When no file is provided with the form set resourceURL to ""
			resourceURL = ""
		default:
			helpers.HandleInternalServerError(ctx, err)
			return
		}
	} else {
		// When a file is provided with the form, attempt to upload it as the memo's new storage
		// and update the memoURL to be saved as part of the memo
		resourceURL, err = mh.app.Repositories.File.UploadMemoMedia(newMemo.ID, resourceFile, newMemo.Type)
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

	newMemo.ResourceURL = resourceURL
	fmt.Println("BIG BAD WOLF 1")
	fmt.Println(newMemo.ResourceURL)

	updatedMemo, err := mh.app.Repositories.Memo.Update(newMemo.ID, newMemo)
	fmt.Println("BIG BAD WOLF 2")
	fmt.Println(updatedMemo.ResourceURL)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	data := response.MemoResponseFromModel(updatedMemo)

	// return newly created memo
	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"data":    data,
			"message": `Memo was successfully created.`,
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

func (mh memoHandler) FindMemos(ctx *gin.Context) {
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// retrieve query params for search
	searchString := ctx.Query("searchString")

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
	memos, err := mh.app.Repositories.Memo.FindMemos(searchString, page, pageSize)
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
	if memo.Type != "text" {
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
