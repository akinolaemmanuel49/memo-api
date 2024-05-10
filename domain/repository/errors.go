package repository

import "errors"

var (
	ErrDuplicateDetails   = errors.New("username or email already exists")
	ErrRecordNotFound     = errors.New("no matching record found")
	ErrRecordDeleted      = errors.New("record has been deleted")
	ErrUnapprovedFileType = errors.New("provided file type is not allowed")
	ErrDuplicateFollow    = errors.New("identical follow instance already exists")
	ErrCheckFollow        = errors.New("followerID and subjectID must not be the same")
	ErrMemoIDQueryMissing = errors.New("memoID is missing in the URL query parameter")
	ErrConcurrentUpdate   = errors.New("concurrent update detected")
)
