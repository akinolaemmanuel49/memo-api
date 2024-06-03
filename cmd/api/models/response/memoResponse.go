package response

import (
	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"time"
)

type Memo struct {
	ID         string    `json:"id,omitempty"`
	MemoType   string    `json:"memo_type"`
	Content    string    `json:"content"`
	Likes      int64     `json:"likes,omitempty"`
	Shares     int64     `json:"shares,omitempty"`
	Caption    string    `json:"caption,omitempty"`
	Transcript string    `json:"transcript,omitempty"`
	Deleted    bool      `json:"deleted,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	OwnerID    string    `json:"owner_id,omitempty"`
}

func MemoResponseFromModel(memo models.Memo) Memo {
	return Memo{
		ID:         memo.ID,
		MemoType:   memo.MemoType,
		Content:    memo.Content,
		Likes:      memo.Likes,
		Shares:     memo.Shares,
		Caption:    memo.Caption,
		Transcript: memo.Transcript,
		Deleted:    memo.Deleted,
		CreatedAt:  memo.CreatedAt,
		UpdatedAt:  memo.UpdatedAt,
		OwnerID:    memo.OwnerID,
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
