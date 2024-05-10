package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"github.com/akinolaemmanuel49/memo-api/domain/repository"
	"github.com/akinolaemmanuel49/memo-api/internal/helpers"
)

type social struct {
	Db *sql.DB
}

func NewSocialInfrastructure(db *sql.DB) repository.SocialRepository {
	return social{Db: db}
}

const (
	duplicateFollowerSubjectPair = "unique_follower_subject_pair"
	checkFollowerSubjectPair     = "check_different_ids"
)

// Follow creates a new instance for a follow relationship between two users.
func (s social) Follow(followerID, subjectID string) (models.Follow, error) {
	query := `
	INSERT INTO public.follow(follower_id, subject_id)
	VALUES($1, $2)
	RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	newFollow := models.Follow{
		FollowerID: followerID,
		SubjectID:  subjectID,
	}

	err := s.Db.QueryRowContext(
		ctx,
		query,
		followerID,
		subjectID,
	).Scan(&newFollow.ID, &newFollow.CreatedAt, &newFollow.UpdatedAt)

	if err != nil {
		switch {
		case
			strings.Contains(err.Error(), duplicateFollowerSubjectPair):
			return models.Follow{}, repository.ErrDuplicateFollow
		case
			strings.Contains(err.Error(), checkFollowerSubjectPair):
			return models.Follow{}, repository.ErrCheckFollow
		default:
			return models.Follow{}, err
		}
	}
	err = s.updateFollowerCount(subjectID, 1)
	if err != nil {
		return models.Follow{}, err
	}
	err = s.updateFollowingCount(followerID, 1)
	if err != nil {
		return models.Follow{}, err
	}
	return newFollow, nil
}

// Unfollow deletes an instance for a follow relationship between two users.
func (s social) Unfollow(followerID, subjectID string) (models.Follow, error) {
	query := `
        DELETE FROM public.follow
        WHERE follower_id = $1 AND subject_id = $2
    `

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	_, err := s.Db.ExecContext(ctx, query, followerID, subjectID)
	if err != nil {
		return models.Follow{}, err // Return the error directly
	}

	// If the deletion was successful, construct and return the Follow struct
	newUnfollow := models.Follow{
		FollowerID: followerID,
		SubjectID:  subjectID,
	}

	// Update follower and following counts
	err = s.updateFollowerCount(subjectID, -1)
	if err != nil {
		return models.Follow{}, err
	}
	err = s.updateFollowingCount(followerID, -1)
	if err != nil {
		return models.Follow{}, err
	}

	return newUnfollow, nil
}

// UpdateFollowerCount updates the FollowerCount for a given user by a specified increment.
func (s social) updateFollowerCount(userID string, increment int) error {
	query := `
		UPDATE public.users
		SET follower_count = follower_count + $1
		WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	_, err := s.Db.ExecContext(ctx, query, increment, userID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateFollowingCount updates the FollowingCount for a given user by a specified increment.
func (s social) updateFollowingCount(userID string, increment int) error {
	query := `
		UPDATE public.users
		SET following_count = following_count + $1
		WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	_, err := s.Db.ExecContext(ctx, query, increment, userID)
	if err != nil {
		return err
	}

	return nil
}

// CreateComment creates a new instance of a comment in the comments table
func (s social) CreateComment(comment *models.Comment) (models.Comment, error) {
	query := `
	INSERT INTO public.comments(owner_id, memo_id, comment_type, comment_content, caption, transcript, parent_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return models.Comment{}, err
	}

	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	newComment := *comment
	err = tx.QueryRowContext(
		ctx,
		query,
		comment.OwnerID,
		comment.MemoID,
		comment.CommentType,
		comment.Content,
		comment.Caption,
		comment.Transcript,
		comment.ParentID,
	).Scan(&newComment.ID, &newComment.CreatedAt, &newComment.UpdatedAt)

	if err != nil {
		switch {
		default:
			return models.Comment{}, err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return models.Comment{}, err
	}

	return newComment, nil
}

func (s social) UpdateComment(id string, updatedComment models.Comment) (models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	// Query statements
	selectQuery := `SELECT _version from public.comments WHERE id = $1 FOR NO KEY UPDATE;`
	updateQuery := `
	UPDATE public.comments
		SET
		    comment_content = $1,
		    updated_at = $2,
		    _version = _version + 1
		WHERE id = $3 AND _version=$4;`

	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return models.Comment{}, err
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
			return models.Comment{}, repository.ErrRecordNotFound
		default:
			return models.Comment{}, err
		}
	}
	// Check versions
	if currentVersion != updatedComment.Version {
		return models.Comment{}, repository.ErrConcurrentUpdate
	}

	// Update memo instance
	_, err = tx.ExecContext(ctx,
		updateQuery,
		updatedComment.Content,
		time.Now().UTC(),
		id,
		updatedComment.Version)
	// Handle errors arising from update
	if err != nil {
		switch {
		default:
			return models.Comment{}, err
		}
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return models.Comment{}, err
	}
	// Fetch updated memo
	updatedComment, err = s.GetComment(id)
	if err != nil {
		return models.Comment{}, err
	}
	return updatedComment, nil
}

// GetComment retrieves an existing comment via its ID.
// repository.ErrRecordNotFound is returned if no text memo matches the query.
func (s social) GetComment(id string) (models.Comment, error) {
	query := `
	SELECT
		id,
		memo_id,
		parent_id,
		comment_content,
		comment_type,
		likes,
		caption,
		transcript,
		deleted,
		created_at,
		updated_at,
		owner_id,
		_version
	FROM public.comments
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	foundComment := models.Comment{}
	err := s.Db.QueryRowContext(ctx, query, id).
		Scan(
			&foundComment.ID,
			&foundComment.MemoID,
			&foundComment.ParentID,
			&foundComment.Content,
			&foundComment.CommentType,
			&foundComment.Likes,
			&foundComment.Caption,
			&foundComment.Transcript,
			&foundComment.Deleted,
			&foundComment.CreatedAt,
			&foundComment.UpdatedAt,
			&foundComment.OwnerID,
			&foundComment.Version,
		)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Comment{}, repository.ErrRecordNotFound

		default:
			return models.Comment{}, err
		}
	}

	if foundComment.Deleted {
		return foundComment, repository.ErrRecordDeleted
	}

	return foundComment, nil
}

func (s social) GetCommentsByMemoID(memoID string, page, pageSize int) ([]models.Comment, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	query := `
	SELECT 
		id,
       memo_id,
       parent_id,
       comment_content,
       comment_type,
       likes,
       caption,
       transcript,
       deleted,
       created_at,
       updated_at,
       owner_id
	FROM public.comments
	WHERE memo_id = $1 AND parent_id IS NULL
	ORDER BY created_at
	OFFSET $2 LIMIT $3`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	rows, err := s.Db.QueryContext(ctx, query, memoID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// Handle error
		}
	}(rows)

	comments := make([]models.Comment, 0)

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.MemoID,
			&comment.ParentID,
			&comment.Content,
			&comment.CommentType,
			&comment.Likes,
			&comment.Caption,
			&comment.Transcript,
			&comment.Deleted,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.OwnerID,
		)
		if err != nil {
			return nil, err
		}

		if !(comment.ParentID.Valid) {
			// If parent ID is nil, it's a top-level comment
			comments = append(comments, comment)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s social) GetRepliesByParentID(parentID string, page, pageSize int) ([]models.Comment, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	query := `
	SELECT 
		id,
       memo_id,
       parent_id,
       comment_content,
       comment_type,
       likes,
       caption,
       transcript,
       deleted,
       created_at,
       updated_at,
       owner_id
	FROM public.comments
	WHERE parent_id = $1
	ORDER BY created_at
	OFFSET $2 LIMIT $3`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	rows, err := s.Db.QueryContext(ctx, query, parentID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// Handle error
		}
	}(rows)

	replies := make([]models.Comment, 0)

	for rows.Next() {
		var reply models.Comment
		err := rows.Scan(
			&reply.ID,
			&reply.MemoID,
			&reply.ParentID,
			&reply.Content,
			&reply.CommentType,
			&reply.Likes,
			&reply.Caption,
			&reply.Transcript,
			&reply.Deleted,
			&reply.CreatedAt,
			&reply.UpdatedAt,
			&reply.OwnerID,
		)
		if err != nil {
			return nil, err
		}

		if reply.ParentID.Valid {
			// If parent ID is not nil, it's not a top-level comment
			replies = append(replies, reply)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return replies, nil
}

// Function to retrieve replies recursively
//func retrieveReplies(comment *models.Comment, commentMap map[string][]models.Comment) {
//	replies, ok := commentMap[comment.ID]
//	if ok {
//		comment.Replies = replies
//		for i := range replies {
//			retrieveReplies(&replies[i], commentMap)
//		}
//	}
//}
