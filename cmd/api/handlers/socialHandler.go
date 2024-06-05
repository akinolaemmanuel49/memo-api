package handlers

import (
	"database/sql"
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

type SocialHandler interface {
	Follow(ctx *gin.Context)
	Unfollow(ctx *gin.Context)
	IsFollower(ctx *gin.Context)
	CreateTextComment(ctx *gin.Context)
	CreateTextReply(ctx *gin.Context)
	GetComments(ctx *gin.Context)
	GetReplies(ctx *gin.Context)
}

type socialHandler struct {
	app internal.Application
}

func NewSocialHandler(app internal.Application) SocialHandler {
	return socialHandler{app: app}
}

// Follow creates a new relationship between an authenticated user and another user.
func (sh socialHandler) Follow(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	subjectID := ctx.Param("subjectID")
	if subjectID == user.ID {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, repository.ErrCheckFollow)
		return
	}

	// attempt to follow a user
	newFollow, err := sh.app.Repositories.Social.Follow(user.ID, subjectID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateFollow):
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		case errors.Is(err, repository.ErrCheckFollow):
			helpers.HandleErrorResponse(ctx, http.StatusBadRequest, err)
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	followedUser, err := sh.app.Repositories.Users.GetById(newFollow.SubjectID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, err)

		default:
			fmt.Println(err)
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	responseMessage := gin.H{
		"status":  "success",
		"message": "User successfully followed.",
		"data": gin.H{
			"userID":     followedUser.ID,
			"username":   followedUser.Username,
			"followedAt": newFollow.CreatedAt,
		},
	}

	ctx.JSON(
		http.StatusOK,
		responseMessage,
	)
}

// Unfollow deletes an existing relationship between an authenticated user and another user.
func (sh socialHandler) Unfollow(ctx *gin.Context) {
	// fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)
	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// retrieve the user id for the user to follow from query params
	subjectID := ctx.Param("subjectID")
	if subjectID == user.ID {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, repository.ErrCheckFollow)
		return
	}

	// unfollow a user
	newUnfollow, err := sh.app.Repositories.Social.Unfollow(user.ID, subjectID)
	if err != nil {
		switch {
		default:
			fmt.Println(err)
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	unfollowedUser, err := sh.app.Repositories.Users.GetById(newUnfollow.SubjectID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, err)

		default:
			fmt.Println(err)
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	responseMessage := gin.H{
		"status":  "success",
		"message": "User successfully unfollowed.",
		"data": gin.H{
			"userID":       unfollowedUser.ID,
			"username":     unfollowedUser.Username,
			"unfollowedAt": newUnfollow.CreatedAt,
		},
	}

	ctx.JSON(
		http.StatusOK,
		responseMessage,
	)
}

// IsFollower checks if the authenticated user is following another user.
func (sh socialHandler) IsFollower(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	subjectID := ctx.Param("subjectID")
	if subjectID == user.ID {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, repository.ErrCheckFollow)
		return
	}

	// Check if the user is following the subject
	isFollower, err := sh.app.Repositories.Social.IsFollower(user.ID, subjectID)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	responseMessage := gin.H{
		"status":     "success",
		"isFollower": isFollower,
	}

	ctx.JSON(
		http.StatusOK,
		responseMessage,
	)
}

// CreateTextComment creates a new instance of a text based comment.
func (sh socialHandler) CreateTextComment(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	memoID := ctx.Param("memoID")

	comment := ctx.PostForm("comment")

	textComment := models.Comment{
		OwnerID:     user.ID,
		MemoID:      memoID,
		CommentType: "text",
		Content:     comment,
	}

	newTextComment, err := sh.app.Repositories.Social.CreateComment(&textComment)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return newly created text comment
	ctx.JSON(
		http.StatusCreated,
		response.CommentResponseFromModel(newTextComment))
}

// CreateImageComment creates a new instance of an image based comment.
func (sh socialHandler) CreateImageComment(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	memoID := ctx.Param("memoID")
	// Validate request data
	caption := ctx.PostForm("caption")

	sqlCaption := sql.NullString{
		String: caption,
		Valid:  true,
	}

	imageComment := models.Comment{
		OwnerID:     user.ID,
		MemoID:      memoID,
		CommentType: "image",
		Caption:     sqlCaption,
	}

	// attempt to save image memo in repository
	newImageComment, err := sh.app.Repositories.Social.CreateComment(&imageComment)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// Process the file provided with the form as the new storage of the comment
	var commentURL string
	commentFile, _, err := ctx.Request.FormFile("commentFile")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			// When no file is provided with the form set memoURL to ""
			commentURL = ""
		default:
			helpers.HandleInternalServerError(ctx, err)
			return
		}
	} else {
		// When a file is provided with the form, attempt to upload it as the comment's new storage
		// and update the commentURL to be saved as part of the comment
		commentURL, err = sh.app.Repositories.File.UploadCommentMedia(newImageComment.ID, commentFile, newImageComment.CommentType)
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

	newImageComment.Content = commentURL

	updatedComment, err := sh.app.Repositories.Social.UpdateComment(newImageComment.ID, newImageComment)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	data := response.CommentResponseFromModel(updatedComment)

	// return newly create image memo
	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"data":    data,
			"message": `Image comment was successfully created.`,
		},
	)
}

// CreateAudioComment creates a new instance of an audio based comment.
func (sh socialHandler) CreateAudioComment(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	memoID := ctx.Param("memoID")
	// Validate request data
	caption := ctx.PostForm("caption")
	transcript := ctx.PostForm("transcript")

	sqlCaption := sql.NullString{
		String: caption,
		Valid:  true,
	}

	sqlTranscript := sql.NullString{
		String: transcript,
		Valid:  true,
	}

	audioComment := models.Comment{
		OwnerID:     user.ID,
		MemoID:      memoID,
		CommentType: "audio",
		Caption:     sqlCaption,
		Transcript:  sqlTranscript,
	}

	// attempt to save image memo in repository
	newAudioComment, err := sh.app.Repositories.Social.CreateComment(&audioComment)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// Process the file provided with the form as the new storage of the comment
	var commentURL string
	commentFile, _, err := ctx.Request.FormFile("commentFile")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			// When no file is provided with the form set memoURL to ""
			commentURL = ""
		default:
			helpers.HandleInternalServerError(ctx, err)
			return
		}
	} else {
		// When a file is provided with the form, attempt to upload it as the comment's new storage
		// and update the commentURL to be saved as part of the comment
		commentURL, err = sh.app.Repositories.File.UploadCommentMedia(newAudioComment.ID, commentFile, newAudioComment.CommentType)
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

	newAudioComment.Content = commentURL

	updatedComment, err := sh.app.Repositories.Social.UpdateComment(newAudioComment.ID, newAudioComment)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	data := response.CommentResponseFromModel(updatedComment)

	// return newly create audio memo
	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"data":    data,
			"message": `Audio comment was successfully created.`,
		},
	)
}

// CreateVideoComment creates a new instance of a video based comment.
func (sh socialHandler) CreateVideoComment(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	memoID := ctx.Param("memoID")
	// Validate request data
	caption := ctx.PostForm("caption")
	transcript := ctx.PostForm("transcript")

	sqlCaption := sql.NullString{
		String: caption,
		Valid:  true,
	}

	sqlTranscript := sql.NullString{
		String: transcript,
		Valid:  true,
	}

	videoComment := models.Comment{
		OwnerID:     user.ID,
		MemoID:      memoID,
		CommentType: "video",
		Caption:     sqlCaption,
		Transcript:  sqlTranscript,
	}

	// attempt to save image memo in repository
	newVideoComment, err := sh.app.Repositories.Social.CreateComment(&videoComment)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// Process the file provided with the form as the new storage of the comment
	var commentURL string
	commentFile, _, err := ctx.Request.FormFile("commentFile")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			// When no file is provided with the form set memoURL to ""
			commentURL = ""
		default:
			helpers.HandleInternalServerError(ctx, err)
			return
		}
	} else {
		// When a file is provided with the form, attempt to upload it as the comment's new storage
		// and update the commentURL to be saved as part of the comment
		commentURL, err = sh.app.Repositories.File.UploadCommentMedia(newVideoComment.ID, commentFile, newVideoComment.CommentType)
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

	newVideoComment.Content = commentURL

	updatedComment, err := sh.app.Repositories.Social.UpdateComment(newVideoComment.ID, newVideoComment)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	data := response.CommentResponseFromModel(updatedComment)

	// return newly create video memo
	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"data":    data,
			"message": `Video comment was successfully created.`,
		},
	)
}

// CreateTextReply creates a new instance of a text based reply.
func (sh socialHandler) CreateTextReply(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	parentID := ctx.Param("parentID")
	fmt.Print(parentID)

	if parentID == "" {
		helpers.HandleErrorResponse(ctx, http.StatusBadRequest, errors.New("parentID parameter is required"))
		return
	}

	parentComment, err := sh.app.Repositories.Social.GetComment(parentID)
	fmt.Println(parentComment)

	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	memoID := parentComment.MemoID
	sqlParentID := sql.NullString{
		String: parentID,
		Valid:  true,
	}

	reply := ctx.PostForm("reply")

	textReply := models.Comment{
		OwnerID:     user.ID,
		MemoID:      memoID,
		CommentType: "text",
		Content:     reply,
		ParentID:    sqlParentID,
	}

	newTextReply, err := sh.app.Repositories.Social.CreateComment(&textReply)
	if err != nil {
		switch {
		default:
			helpers.HandleInternalServerError(ctx, err)
		}
		return
	}

	// return newly created text comment
	ctx.JSON(
		http.StatusCreated,
		response.CommentResponseFromModel(newTextReply))
}

// GetComments retrieves all comments for a memo.
func (sh socialHandler) GetComments(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// retrieve queries and params from URL
	memoID := ctx.Param("memoID")

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

	// Fetch comments by memoID
	comments, err := sh.app.Repositories.Social.GetCommentsByMemoID(memoID, page, pageSize)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	returned := response.MultipleCommentResponseFromModel(comments)

	// Return comments
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   returned,
	})
}

func (sh socialHandler) GetReplies(ctx *gin.Context) {
	// Fetch authenticated user from context and return authentication error if no user exists
	user := helpers.ContextGetUser(ctx)

	if reflect.DeepEqual(user, models.User{}) {
		helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthenticated"))
		return
	}

	// retrieve queries and params from URL
	commentID := ctx.Param("commentID")

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

	// Fetch replies by commentID
	replies, err := sh.app.Repositories.Social.GetRepliesByParentID(commentID, page, pageSize)
	if err != nil {
		helpers.HandleInternalServerError(ctx, err)
		return
	}

	returned := response.MultipleCommentResponseFromModel(replies)

	// Return comments
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   returned,
	})
}
