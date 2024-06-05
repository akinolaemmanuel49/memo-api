package response

import (
	"time"

	"github.com/akinolaemmanuel49/memo-api/domain/models"
)

type Memo struct {
	ID          string    `json:"id,omitempty"`
	Type        string    `json:"type"`
	Content     string    `json:"content"`
	Likes       int64     `json:"likes,omitempty"`
	Shares      int64     `json:"shares,omitempty"`
	ResourceURL string    `json:"resource_url,omitempty"`
	Description string    `json:"description,omitempty"`
	Deleted     bool      `json:"deleted,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	OwnerID     string    `json:"owner_id,omitempty"`
}

func MemoResponseFromModel(memo models.Memo) Memo {
	return Memo{
		ID:          memo.ID,
		Type:        memo.Type,
		Content:     memo.Content,
		Likes:       memo.Likes,
		Shares:      memo.Shares,
		ResourceURL: memo.ResourceURL,
		Description: memo.Description,
		Deleted:     memo.Deleted,
		CreatedAt:   memo.CreatedAt,
		UpdatedAt:   memo.UpdatedAt,
		OwnerID:     memo.OwnerID,
	}
}

func MultipleMemoResponseFromModel(memos []models.Memo) []Memo {
	var memoResponses []Memo
	for _, memo := range memos {
		memoResponse := MemoResponseFromModel(memo)
		memoResponses = append(memoResponses, memoResponse)
	}
	return memoResponses
}
