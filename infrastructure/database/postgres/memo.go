package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"github.com/akinolaemmanuel49/memo-api/domain/repository"
	"github.com/akinolaemmanuel49/memo-api/internal/helpers"
)

type memo struct {
	Db *sql.DB
}

func NewMemoInfrastructure(db *sql.DB) repository.MemoRepository {
	return memo{Db: db}
}

// CreateMemo creates and returns an instance of a new text memo,
// it returns an error if ownerID is not set.
func (m memo) CreateMemo(ownerID string, memo *models.Memo) (models.Memo, error) {
	query := `
	INSERT INTO public.memos(content, owner_id, type, resource_url, description)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.UploadTimeoutDuration)
	defer cancel()

	newMemo := *memo
	err := m.Db.QueryRowContext(
		ctx,
		query,
		memo.Content,
		ownerID,
		memo.Type,
		memo.ResourceURL,
		memo.Description,
	).Scan(&newMemo.ID, &newMemo.CreatedAt, &newMemo.UpdatedAt)

	if err != nil {
		switch {
		default:
			return models.Memo{}, err
		}
	}

	return newMemo, nil
}

// GetMemo retrieves an existing memo via its ID.
// repository.ErrRecordNotFound is returned if no text memo matches the query.
func (m memo) GetMemo(id string) (models.Memo, error) {
	query := `
	SELECT
		id,
		content,
		type,
		likes,
		shares,
		resource_url,
		description,
		deleted,
		created_at,
		updated_at,
		owner_id,
		_version
	FROM public.memos
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	foundMemo := models.Memo{}
	err := m.Db.QueryRowContext(ctx, query, id).
		Scan(
			&foundMemo.ID,
			&foundMemo.Content,
			&foundMemo.Type,
			&foundMemo.Likes,
			&foundMemo.Shares,
			&foundMemo.ResourceURL,
			&foundMemo.Description,
			&foundMemo.Deleted,
			&foundMemo.CreatedAt,
			&foundMemo.UpdatedAt,
			&foundMemo.OwnerID,
			&foundMemo.Version,
		)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Memo{}, repository.ErrRecordNotFound

		default:
			return models.Memo{}, err
		}
	}

	if foundMemo.Deleted {
		return foundMemo, repository.ErrRecordDeleted
	}

	return foundMemo, nil
}

// GetAllMemos fetches all memo instances from all users.
func (m memo) GetAllMemos(page, pageSize int) ([]models.Memo, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	// Query for all posts made by all users.
	query := `
	SELECT
		id,
		content,
		type,
		likes,
		shares,
		resource_url,
		description,
		deleted,
		created_at,
		updated_at,
		owner_id
	FROM public.memos
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	rows, err := m.Db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	memos := make([]models.Memo, 0)
	for rows.Next() {
		var memo models.Memo
		err := rows.Scan(
			&memo.ID,
			&memo.Content,
			&memo.Type,
			&memo.Likes,
			&memo.Shares,
			&memo.ResourceURL,
			&memo.Description,
			&memo.Deleted,
			&memo.CreatedAt,
			&memo.UpdatedAt,
			&memo.OwnerID,
		)
		if err != nil {
			return nil, err
		}
		memos = append(memos, memo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memos, nil
}

func (m memo) FindMemos(searchString string, page, pageSize int) ([]models.Memo, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	// SQL query with conditional search based on memo_type
	query := `
	SELECT
		id,
		content,
		type,
		likes,
		shares,
		resource_url,
		description,
		deleted,
		created_at,
		updated_at,
		owner_id
	FROM public.memos
	WHERE content ILIKE $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	rows, err := m.Db.QueryContext(ctx, query, "%"+searchString+"%", pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	memos := make([]models.Memo, 0)
	for rows.Next() {
		var memo models.Memo
		err := rows.Scan(
			&memo.ID,
			&memo.Content,
			&memo.Type,
			&memo.Likes,
			&memo.Shares,
			&memo.ResourceURL,
			&memo.Description,
			&memo.Deleted,
			&memo.CreatedAt,
			&memo.UpdatedAt,
			&memo.OwnerID,
		)
		if err != nil {
			return nil, err
		}
		memos = append(memos, memo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memos, nil
}

// GetMemosByFollowing fetches all memo instances from followed users.
func (m memo) GetMemosByFollowing(userID string, page, pageSize int) ([]models.Memo, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	// Query for all posts made by those followed users.
	query := `
	SELECT
		id,
		content,
		type,
		likes,
		shares,
		resource_url,
		description,
		deleted,
		created_at,
		updated_at,
		owner_id
	FROM public.memos
	WHERE owner_id IN (
		SELECT subject_id::uuid
		FROM public.follow
		WHERE follower_id = $1)
		OR owner_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	rows, err := m.Db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	memos := make([]models.Memo, 0)
	for rows.Next() {
		var memo models.Memo
		err := rows.Scan(
			&memo.ID,
			&memo.Content,
			&memo.Type,
			&memo.Likes,
			&memo.Shares,
			&memo.ResourceURL,
			&memo.Description,
			&memo.Deleted,
			&memo.CreatedAt,
			&memo.UpdatedAt,
			&memo.OwnerID,
		)
		if err != nil {
			return nil, err
		}
		memos = append(memos, memo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memos, nil
}

// LikeMemo creates a new instance in the likes table and increments the number of likes on the memos table.
func (m memo) LikeMemo(likerID string, memoID string) (models.Like, error) {
	query := `INSERT INTO public.likes(liked_by, memo_id) VALUES($1, $2) RETURNING id, created_at, updated_at`
	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	newLike := models.Like{
		LikedBy: likerID,
		MemoID:  memoID,
	}

	err := m.Db.QueryRowContext(
		ctx,
		query,
		likerID,
		memoID).Scan(&newLike.ID, &newLike.CreatedAt, &newLike.UpdatedAt)
	if err != nil {
		switch {
		default:
			return models.Like{}, err
		}
	}

	// increment memo like count
	err = m.updateLikeCount(memoID, 1)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Like{}, repository.ErrRecordNotFound
		default:
			return models.Like{}, err
		}
	}
	return newLike, nil
}

// UnlikeMemo deletes an instance with matching parameters from the likes table.
func (m memo) UnlikeMemo(likerID string, memoID string) error {
	query := `DELETE FROM public.likes WHERE liked_by = $1 AND memo_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	_, err := m.Db.ExecContext(ctx, query, likerID, memoID)
	if err != nil {
		switch {
		default:
			return err
		}
	}

	// decrement memo like count
	err = m.updateLikeCount(memoID, -1)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return repository.ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

// updateLikeCount increments/decrements the like value for an instance with matching memo id in the memos table.
func (m memo) updateLikeCount(memoID string, increment int) error {
	query := `UPDATE public.memos SET likes = likes + $1 WHERE id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	_, err := m.Db.ExecContext(ctx, query, increment, memoID)
	if err != nil {
		return err
	}
	return nil
}

// ShareMemo creates a new instance in the shares table and increments the number of shares on the memos table.
func (m memo) ShareMemo(sharerID string, memoID string) (models.Share, error) {
	query := `INSERT INTO public.shares(shared_by, memo_id) VALUES($1, $2) RETURNING id, created_at, updated_at`
	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	newShare := models.Share{
		SharedBy: sharerID,
		MemoID:   memoID,
	}

	err := m.Db.QueryRowContext(
		ctx,
		query,
		sharerID,
		memoID).Scan(&newShare.ID, &newShare.CreatedAt, &newShare.UpdatedAt)
	if err != nil {
		switch {
		default:
			return models.Share{}, err
		}
	}

	// increment memo share count
	err = m.updateShareCount(memoID, 1)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Share{}, repository.ErrRecordNotFound
		default:
			return models.Share{}, err
		}
	}
	return newShare, nil
}

// UnshareMemo deletes an instance with matching parameters from the shares table.
func (m memo) UnshareMemo(sharerID string, memoID string) error {
	query := `DELETE FROM public.shares WHERE shared_by = $1 AND memo_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	_, err := m.Db.ExecContext(ctx, query, sharerID, memoID)
	if err != nil {
		switch {
		default:
			return err
		}
	}

	// decrement memo share count
	err = m.updateShareCount(memoID, -1)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return repository.ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

// updateShareCount increments/decrements the share value for an instance with matching memo id in the memos table.
func (m memo) updateShareCount(memoID string, value int) error {
	query := `UPDATE public.memos SET shares = shares + $1 WHERE id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	_, err := m.Db.ExecContext(ctx, query, value, memoID)
	if err != nil {
		return err
	}
	return nil
}

func (m memo) Update(id string, updatedMemo models.Memo) (models.Memo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	// Query statements
	selectQuery := `SELECT _version from public.memos WHERE id = $1 FOR NO KEY UPDATE;`
	updateQuery := `
	UPDATE public.memos
		SET
		    resource_url = $1,
		    updated_at = $2,
		    _version = _version + 1
		WHERE id = $3 AND _version=$4;`

	tx, err := m.Db.BeginTx(ctx, nil)
	if err != nil {
		return models.Memo{}, err
	}

	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	var currentVersion int
	// Fetch current version value from user instance
	err = tx.QueryRowContext(ctx, selectQuery, id).Scan(&currentVersion)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Memo{}, repository.ErrRecordNotFound
		default:
			return models.Memo{}, err
		}
	}
	// Check versions
	if currentVersion != updatedMemo.Version {
		return models.Memo{}, repository.ErrConcurrentUpdate
	}

	// Update memo instance
	_, err = tx.ExecContext(ctx,
		updateQuery,
		updatedMemo.ResourceURL,
		time.Now().UTC(),
		id,
		updatedMemo.Version)
	// Handle errors arising from update
	if err != nil {
		switch {
		default:
			return models.Memo{}, err
		}
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return models.Memo{}, err
	}
	// Fetch updated memo
	updatedMemo, err = m.GetMemo(id)
	if err != nil {
		return models.Memo{}, err
	}
	return updatedMemo, nil
}

func (m memo) Delete(id string, deletedMemo models.Memo) (models.Memo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	// Query statements
	selectQuery := `SELECT _version from public.memos WHERE id = $1 FOR NO KEY UPDATE;`
	deleteQuery := `
	UPDATE public.memos
		SET
		    deleted = TRUE,
		    updated_at = $1,
		    _version = _version + 1
		WHERE id = $2 AND _version=$3;`

	tx, err := m.Db.BeginTx(ctx, nil)
	if err != nil {
		return models.Memo{}, err
	}

	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	var currentVersion int
	// Fetch current version value from user instance
	err = tx.QueryRowContext(ctx, selectQuery, id).Scan(&currentVersion)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Memo{}, repository.ErrRecordNotFound
		default:
			return models.Memo{}, err
		}
	}
	// Check versions
	if currentVersion != deletedMemo.Version {
		return models.Memo{}, repository.ErrConcurrentUpdate
	}

	// Set deleted flag to TRUE
	_, err = tx.ExecContext(ctx,
		deleteQuery,
		time.Now().UTC(),
		id,
		deletedMemo.Version)
	// Handle errors arising from delete
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Memo{}, repository.ErrRecordNotFound
		default:
			return models.Memo{}, err
		}
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return models.Memo{}, err
	}

	return deletedMemo, nil
}

func (m memo) GetMemosByOwnerID(ownerID string, page, pageSize int) ([]models.Memo, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	// Query for all posts made by user with matching id.
	query := `
	SELECT
		id,
		content,
		type,
		likes,
		shares,
		resource_url,
		description,
		deleted,
		created_at,
		updated_at,
		owner_id
	FROM public.memos
	WHERE owner_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	rows, err := m.Db.QueryContext(ctx, query, ownerID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	memos := make([]models.Memo, 0)
	for rows.Next() {
		var memo models.Memo
		err := rows.Scan(
			&memo.ID,
			&memo.Content,
			&memo.Type,
			&memo.Likes,
			&memo.Shares,
			&memo.ResourceURL,
			&memo.Description,
			&memo.Deleted,
			&memo.CreatedAt,
			&memo.UpdatedAt,
			&memo.OwnerID,
		)
		if err != nil {
			return nil, err
		}
		memos = append(memos, memo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memos, nil
}
