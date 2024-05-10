package repository

import "io"

type FileRepository interface {
	UploadAvatar(userID string, avatarFile io.Reader) (avatarURL string, err error)
	DeleteAvatar(userID string) error
	UploadMemoMedia(memoID string, memoFile io.Reader, typeMedia string) (memoURL string, err error)
	DeleteMemoMedia(memoID string) error
	UploadCommentMedia(commentID string, commentFile io.Reader, typeMedia string) (commentURL string, err error)
	DeleteCommentMedia(commentID string) error
}
