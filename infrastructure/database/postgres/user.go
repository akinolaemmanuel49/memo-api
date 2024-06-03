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

type user struct {
	Db *sql.DB
}

func NewUserInfrastructure(db *sql.DB) repository.UserRepository {
	return user{Db: db}
}

// users table constraints
const (
	duplicateUsername = "users_username_key"
	duplicateEmail    = "users_email_key"
)

// Create registers and returns an instance of a new user, it returns an error if a duplicate username or email is used.
// repository.ErrDuplicateDetails is returned if at least the username or the email already exists in the database.
func (u user) Create(user *models.User) (models.User, error) {
	query := `
	INSERT INTO public.users(username, first_name, last_name, email, password_hash)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	newUser := *user
	err := u.Db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	).Scan(&newUser.ID, &newUser.CreatedAt, &newUser.UpdatedAt)

	if err != nil {
		switch {
		case
			strings.Contains(err.Error(), duplicateUsername),
			strings.Contains(err.Error(), duplicateEmail):
			return models.User{}, repository.ErrDuplicateDetails

		default:
			return models.User{}, err
		}
	}

	return newUser, nil
}

// GetById retrieves an existing user via their ID.
// repository.ErrRecordNotFound is returned if no user matches the query.
func (u user) GetById(id string) (models.User, error) {
	query := `
	SELECT 
		id,
		username,
		first_name,
		last_name,
		email,
		password_hash,
		avatar,
		status,
		about,
		follower_count,
		following_count,
		deleted,
		is_activated,
		created_at,
		updated_at,
		_version
	FROM public.users
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	foundUser := models.User{}
	err := u.Db.QueryRowContext(ctx, query, id).
		Scan(
			&foundUser.ID,
			&foundUser.Username,
			&foundUser.FirstName,
			&foundUser.LastName,
			&foundUser.Email,
			&foundUser.Password,
			&foundUser.AvatarURL,
			&foundUser.Status,
			&foundUser.About,
			&foundUser.FollowerCount,
			&foundUser.FollowingCount,
			&foundUser.Deleted,
			&foundUser.IsActivated,
			&foundUser.CreatedAt,
			&foundUser.UpdatedAt,
			&foundUser.Version,
		)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, repository.ErrRecordNotFound

		case foundUser.Deleted:
			return foundUser, repository.ErrRecordDeleted

		default:
			return models.User{}, err
		}
	}

	return foundUser, nil
}

// GetByEmail retrieves an existing user via their email address.
// repository.ErrRecordNotFound is returned if no user matches the query.
func (u user) GetByEmail(email string) (models.User, error) {
	query := `
	SELECT 
		id,
		username,
		first_name,
		last_name,
		email,
		password_hash,
		avatar,
		status,
		about,
		follower_count,
		following_count,
		deleted,
		is_activated,
		created_at,
		updated_at,
		_version
	FROM public.users
	WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	foundUser := models.User{}
	err := u.Db.QueryRowContext(ctx, query, email).
		Scan(
			&foundUser.ID,
			&foundUser.Username,
			&foundUser.FirstName,
			&foundUser.LastName,
			&foundUser.Email,
			&foundUser.Password,
			&foundUser.AvatarURL,
			&foundUser.Status,
			&foundUser.About,
			&foundUser.FollowerCount,
			&foundUser.FollowingCount,
			&foundUser.Deleted,
			&foundUser.IsActivated,
			&foundUser.CreatedAt,
			&foundUser.UpdatedAt,
			&foundUser.Version,
		)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, repository.ErrRecordNotFound

		case foundUser.Deleted:
			return foundUser, repository.ErrRecordDeleted

		default:
			return models.User{}, err
		}
	}

	return foundUser, nil
}

// GetBySearchString retrieves users matching the search string via their first name, last name, username, or email.
// repository.ErrRecordNotFound is returned if no user matches the query.
func (u user) GetBySearchString(searchString string) ([]models.User, error) {
	query := `
	SELECT 
		id,
		username,
		first_name,
		last_name,
		email,
		password_hash,
		avatar,
		status,
		about,
		follower_count,
		following_count,
		deleted,
		is_activated,
		created_at,
		updated_at,
		_version
	FROM public.users
	WHERE
		first_name ILIKE $1 OR
		last_name ILIKE $1 OR
		username ILIKE $1 OR
		email ILIKE $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	// Use a wildcard search for the search string
	searchPattern := "%" + searchString + "%"

	rows, err := u.Db.QueryContext(ctx, query, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.AvatarURL,
			&user.Status,
			&user.About,
			&user.FollowerCount,
			&user.FollowingCount,
			&user.Deleted,
			&user.IsActivated,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version,
		)
		if err != nil {
			return nil, err
		}

		if user.Deleted {
			return nil, repository.ErrRecordDeleted
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, repository.ErrRecordNotFound
	}

	return users, nil
}

// GetAll retrieves a list of all users.
func (u user) GetAll(page, pageSize int) ([]models.User, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	query := `
	SELECT 
		id,
		username,
		first_name,
		last_name,
		email,
		avatar,
		status,
		about,
		follower_count,
		following_count,
		deleted,
		created_at,
		updated_at
	FROM public.users
	LIMIT $1 OFFSET $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	rows, err := u.Db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var users []models.User
	for rows.Next() {
		var follower models.User
		err := rows.Scan(
			&follower.ID,
			&follower.Username,
			&follower.FirstName,
			&follower.LastName,
			&follower.Email,
			&follower.AvatarURL,
			&follower.Status,
			&follower.About,
			&follower.FollowerCount,
			&follower.FollowingCount,
			&follower.Deleted,
			&follower.CreatedAt,
			&follower.UpdatedAt,
		)

		if err != nil {
			switch {
			default:
				return nil, err
			}
		}
		users = append(users, follower)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetFollowersOfUser retrieves a list of users that follow the user with matching id.
func (u user) GetFollowersOfUser(id string, page, pageSize int) ([]models.User, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	query := `
	SELECT 
		u.id,
		u.username,
		u.first_name,
		u.last_name,
		u.email,
		u.avatar,
		u.status,
		u.about,
		u.follower_count,
		u.following_count,
		u.deleted,
		u.created_at,
		u.updated_at
	FROM public.users u
	JOIN public.follow f ON u.id = f.follower_id
	WHERE f.subject_id = $1
	LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	rows, err := u.Db.QueryContext(ctx, query, id, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var followers []models.User
	for rows.Next() {
		var follower models.User
		err := rows.Scan(
			&follower.ID,
			&follower.Username,
			&follower.FirstName,
			&follower.LastName,
			&follower.Email,
			&follower.AvatarURL,
			&follower.Status,
			&follower.About,
			&follower.FollowerCount,
			&follower.FollowingCount,
			&follower.Deleted,
			&follower.CreatedAt,
			&follower.UpdatedAt,
		)

		if err != nil {
			switch {
			default:
				return nil, err
			}
		}
		followers = append(followers, follower)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

// GetUsersFollowedBy retrieves a list of users that the user with matching id follows.
func (u user) GetUsersFollowedBy(id string, page, pageSize int) ([]models.User, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	query := `
    SELECT 
        u.id,
        u.username,
        u.first_name,
        u.last_name,
        u.email,
        u.avatar,
        u.status,
        u.about,
        u.follower_count,
        u.following_count,
        u.deleted,
        u.created_at,
        u.updated_at
    FROM public.users u
    JOIN public.follow f ON u.id = f.subject_id
    WHERE f.follower_id = $1
    LIMIT $2 OFFSET $3
    `

	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	rows, err := u.Db.QueryContext(ctx, query, id, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var usersFollowed []models.User
	for rows.Next() {
		var userFollowed models.User
		err := rows.Scan(
			&userFollowed.ID,
			&userFollowed.Username,
			&userFollowed.FirstName,
			&userFollowed.LastName,
			&userFollowed.Email,
			&userFollowed.AvatarURL,
			&userFollowed.Status,
			&userFollowed.About,
			&userFollowed.FollowerCount,
			&userFollowed.FollowingCount,
			&userFollowed.Deleted,
			&userFollowed.CreatedAt,
			&userFollowed.UpdatedAt,
		)

		if err != nil {
			switch {
			default:
				return nil, err
			}
		}
		usersFollowed = append(usersFollowed, userFollowed)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usersFollowed, nil
}

// Update updates the instance of the user with matching id.
// repository.ErrRecordNotFound is returned if no user matches the query.
// repository.ErrDuplicateDetails is returned if at least the username or the email already exists in the database.
// repository.ErrConcurrentUpdate is returned if multiple users attempt to update the same record at a time.
func (u user) Update(id string, updatedUser models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), helpers.UploadTimeoutDuration)
	defer cancel()

	// Query statements
	selectQuery := `SELECT _version from public.users WHERE id = $1 FOR NO KEY UPDATE;`
	updateQuery := `
	UPDATE public.users
		SET 
		    username = $1,
		    email = $2,
		    first_name = $3,
		    last_name = $4,
		    password_hash = $5,
		    status = $6,
		    about = $7,
		    avatar = $8,
		    updated_at = $9,
		    _version = _version + 1
		WHERE id = $10 AND _version = $11;`

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return models.User{}, err
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
			return models.User{}, repository.ErrRecordNotFound
		default:
			return models.User{}, err
		}
	}
	// Check versions
	if currentVersion != updatedUser.Version {
		return models.User{}, repository.ErrConcurrentUpdate
	}

	// Update user instance
	_, err = tx.ExecContext(ctx,
		updateQuery,
		updatedUser.Username,
		updatedUser.Email,
		updatedUser.FirstName,
		updatedUser.LastName,
		updatedUser.Password,
		updatedUser.Status,
		updatedUser.About,
		updatedUser.AvatarURL,
		time.Now().UTC(),
		id,
		updatedUser.Version)
	// Handle errors arising from update
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, repository.ErrRecordNotFound
		case strings.Contains(err.Error(), duplicateUsername),
			strings.Contains(err.Error(), duplicateEmail):
			return models.User{}, repository.ErrDuplicateDetails
		default:
			return models.User{}, err
		}
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return models.User{}, err
	}
	// Fetch updated user
	updatedUser, err = u.GetById(id)
	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}

// Delete sets the deleted field  for instance of user with matching id.
// repository.ErrRecordNotFound is returned if no user matches the query.
func (u user) Delete(id string, deletedUser models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), helpers.TimeoutDuration)
	defer cancel()

	// Query statements
	selectQuery := `SELECT _version from public.users WHERE id = $1 FOR NO KEY UPDATE;`
	deleteQuery := `
	UPDATE public.users
	SET
		deleted = TRUE,
		updated_at = $1,
		_version = _version + 1
	WHERE id = $2 AND _version = $3;`

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return models.User{}, err
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
			return models.User{}, repository.ErrRecordNotFound
		default:
			return models.User{}, err
		}
	}
	// Check versions
	if currentVersion != deletedUser.Version {
		return models.User{}, repository.ErrConcurrentUpdate
	}

	// Set deleted flag to TRUE
	_, err = tx.ExecContext(ctx,
		deleteQuery,
		time.Now().UTC(),
		id,
		deletedUser.Version)
	// Handle errors arising from delete
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, repository.ErrRecordNotFound
		default:
			return models.User{}, err
		}
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return models.User{}, err
	}

	return deletedUser, nil
}
