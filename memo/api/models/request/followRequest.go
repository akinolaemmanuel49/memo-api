package request

import (
	"errors"

	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"github.com/akinolaemmanuel49/memo-api/internal/helpers"
)

type Follow struct {
	SubjectID  *string `json:"subjectID" validate:"omitempty"`
	FollowerID *string `json:"followerID" validate:"omitempty"`
}

const (
	FollowFieldSubjectID = iota
	FollowFieldFollowerID
)

func (f Follow) ToModel() models.Follow {
	return models.Follow{
		SubjectID:  helpers.SafeDereference(f.SubjectID),
		FollowerID: helpers.SafeDereference(f.FollowerID),
	}
}

// ValidateRequired verifies that the required fields for the request are provided.
func (f Follow) ValidateRequired(required ...int) error {
	for _, field := range required {
		switch field {
		case FollowFieldSubjectID:
			if f.SubjectID == nil {
				return errors.New("subjectID is required")
			}

		case FollowFieldFollowerID:
			if f.FollowerID == nil {
				return errors.New("followerID is required")
			}
		}
	}
	return nil
}
