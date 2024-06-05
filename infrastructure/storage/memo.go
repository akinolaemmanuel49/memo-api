package storage

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/akinolaemmanuel49/memo-api/domain/repository"
	"github.com/akinolaemmanuel49/memo-api/internal/helpers"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadMemoMedia uploads the memoFile to Cloudinary as the memo's storage,
// using the memoID as the file name.
// The HTTPS URL of the uploaded media is returned if no error is encountered.
// repository.ErrUnapprovedFileType is return if an unsupported file is uploaded.
func (f file) UploadMemoMedia(memoID string, resourceFile io.Reader, typeMedia string) (resourceURL string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), helpers.UploadTimeoutDuration)
	defer cancel()

	// create Cloudinary instance and upload file
	cld, err := cloudinary.NewFromParams(f.CloudName, f.APIKey, f.APISecret)
	if err != nil {
		return "", err
	}

	var allowedFormats []string

	// Add allowed formats based on typeMedia
	switch typeMedia {
	case "image":
		allowedFormats = []string{"jpeg", "jpg", "png", "gif", "bmp"}
	case "video":
		allowedFormats = []string{"mp4", "mov", "avi", "mkv", "wmv"}
	case "audio":
		allowedFormats = []string{"mp3", "wav", "ogg", "aac", "flac"}
	default:
		return "", errors.New("unsupported typeMedia")
	}

	res, err := cld.Upload.Upload(
		ctx,
		resourceFile,
		uploader.UploadParams{
			PublicID:       memoID,
			ResourceType:   typeMedia,
			AllowedFormats: allowedFormats,
			Tags:           []string{"storage"},
			Invalidate:     &Invalidate,
		},
	)
	if err != nil {
		return "", err
	}

	// resErrMessage represents a possible error returned from the Cloudinary server,
	// not one encountered before making the upload request
	if resErrMessage := res.Error.Message; resErrMessage != "" {
		switch {
		case strings.Contains(resErrMessage, "file format") && strings.Contains(resErrMessage, "not allowed"):
			// when a file not satisfying the uploader.UploadParams.AllowedFormats list is found
			// the response error is populated with a message in the format:
			// "Image file format <unapproved_type> not allowed"
			// this case-block takes advantage of that behaviour to return an error when such a situation arises
			return "", repository.ErrUnapprovedFileType
		default:
			// error messages for other situations are yet to be encountered so all other cases have
			// their messages propagated as errors back to the caller
			return "", errors.New(resErrMessage)
		}
	}

	// return URL of uploaded image
	return res.SecureURL, nil
}

// DeleteMemoMedia deletes the memo's media file from Cloudinary storage.
func (f file) DeleteMemoMedia(memoID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create Cloudinary instance to delete file
	cld, err := cloudinary.NewFromParams(f.CloudName, f.APIKey, f.APISecret)
	if err != nil {
		return err
	}

	// Attempt to delete storage with matching memo ID and typeMedia
	res, err := cld.Upload.Destroy(ctx,
		uploader.DestroyParams{
			PublicID:   memoID,
			Invalidate: &Invalidate,
		},
	)
	if err != nil {
		return err
	}
	if resErrMessage := res.Error.Message; resErrMessage != "" {
		return errors.New(resErrMessage)
	}

	return nil
}
