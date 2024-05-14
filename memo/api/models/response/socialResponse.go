package response

import (
	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"time"
)

type Comment struct {
	ID          string    `json:"id,omitempty"`
	OwnerID     string    `json:"owner_id,omitempty"`
	MemoID      string    `json:"memo_id,omitempty"`
	ParentID    string    `json:"parent_id,omitempty"`
	CommentType string    `json:"comment_type"`
	Content     string    `json:"content"`
	Likes       int64     `json:"likes,omitempty"`
	Caption     string    `json:"caption,omitempty"`
	Transcript  string    `json:"transcript,omitempty"`
	Deleted     bool      `json:"deleted,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func CommentResponseFromModel(comment models.Comment) Comment {
	return Comment{
		ID:          comment.ID,
		OwnerID:     comment.OwnerID,
		MemoID:      comment.MemoID,
		ParentID:    comment.ParentID.String,
		CommentType: comment.CommentType,
		Content:     comment.Content,
		Likes:       comment.Likes,
		Caption:     comment.Caption.String,
		Transcript:  comment.Transcript.String,
		Deleted:     comment.Deleted,
		CreatedAt:   comment.CreatedAt,
		UpdatedAt:   comment.UpdatedAt,
	}
}

func MultipleCommentResponseFromModel(comments []models.Comment) []Comment {
	var commentResponses []Comment
	for _, comment := range comments {
		commentResponse := CommentResponseFromModel(comment)
		commentResponses = append(commentResponses, commentResponse)
	}
	return commentResponses
}
