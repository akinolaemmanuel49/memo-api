package repository

import "io"

type FileRepository interface {
	UploadAvatar(userID string, avatarFile io.Reader) (avatarURL string, err error)
	DeleteAvatar(userID string) error
	UploadMemoMedia(memoID string, resourceFile io.Reader, typeMedia string) (resourceURL string, err error)
	DeleteMemoMedia(memoID string) error
	UploadCommentMedia(commentID string, resourceFile io.Reader, typeMedia string) (resourceURL string, err error)
	DeleteCommentMedia(commentID string) error
}
