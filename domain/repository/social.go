package repository

import "github.com/akinolaemmanuel49/memo-api/domain/models"

type SocialRepository interface {
	Follow(followerID, subjectID string) (models.Follow, error)
	Unfollow(followerID, subjectID string) (models.Follow, error)
	CreateComment(comment *models.Comment) (models.Comment, error)
	GetComment(commentID string) (models.Comment, error)
	UpdateComment(commentID string, updatedComment models.Comment) (models.Comment, error)
	GetCommentsByMemoID(memoID string, page, pageSize int) ([]models.Comment, error)
	GetRepliesByParentID(parentID string, page, pageSize int) ([]models.Comment, error)
	//GetRepliesByParentID(parentID string, page, pageSize int) ([]models.Comment, error)
	//GetReplies(ID string, page, pageSize int) ([]models.Comment, error)
	//GetAllComments(memoID, parentID string) ([]models.Comment, error)
	//DeleteComment(ownerID string, commentID string) (models.Comment, error)
}
