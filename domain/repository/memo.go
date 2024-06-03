package repository

import "github.com/akinolaemmanuel49/memo-api/domain/models"

type MemoRepository interface {
	CreateMemo(ownerID string, memo *models.Memo) (models.Memo, error)
	GetMemo(id string) (models.Memo, error)
	LikeMemo(likerID string, memoID string) (models.Like, error)
	UnlikeMemo(likerID string, memoID string) error
	GetAllMemos(page, pageSize int) ([]models.Memo, error)
	FindMemos(searchString string, page, pageSize int) ([]models.Memo, error)
	GetMemosByFollowing(ownerID string, page, pageSize int) ([]models.Memo, error)
	Update(id string, updatedMemo models.Memo) (models.Memo, error)
	Delete(id string, deletedMemo models.Memo) (models.Memo, error)
	ShareMemo(sharerID string, memoID string) (models.Share, error)
	UnshareMemo(sharerID string, memoID string) error
	GetMemosByOwnerID(ownerID string, page, pageSize int) ([]models.Memo, error)
	//ReportMemo(id string) error
}
