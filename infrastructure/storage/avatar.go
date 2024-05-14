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

// UploadAvatar uploads the avatarFile to Cloudinary as the user's profile storage,
// using the userID as the file name.
// The HTTPS URL of the uploaded image is returned if no error is encountered.
// repository.ErrUnapprovedFileType is return if a non-image file is uploaded.
func (f file) UploadAvatar(userID string, avatarFile io.Reader) (avatarURL string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), helpers.UploadTimeoutDuration)
	defer cancel()

	// create Cloudinary instance and upload file
	cld, err := cloudinary.NewFromParams(f.CloudName, f.APIKey, f.APISecret)
	if err != nil {
		return "", err
	}

	res, err := cld.Upload.Upload(
		ctx,
		avatarFile,
		uploader.UploadParams{
			PublicID:       userID,
			ResourceType:   TypeImage,
			AllowedFormats: []string{"jpeg", "jpg", "png"},
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

// DeleteAvatar destroys the file serving as the user's storage.
func (f file) DeleteAvatar(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// create Cloudinary instance to delete file
	cld, err := cloudinary.NewFromParams(f.CloudName, f.APIKey, f.APISecret)
	if err != nil {
		return err
	}

	// attempt to delete storage with matching user ID
	res, err := cld.Upload.Destroy(ctx,
		uploader.DestroyParams{
			PublicID:   userID,
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
