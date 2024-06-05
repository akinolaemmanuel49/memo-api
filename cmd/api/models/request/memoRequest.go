package request

import (
	"errors"

	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"github.com/akinolaemmanuel49/memo-api/internal/helpers"
)

type Memo struct {
	Content     *string `json:"content" validate:"required"`
	Description *string `json:"transcript" validate:"omitempty"`
}

const (
	MemoFieldContent = iota
)

func (m Memo) ToModel() models.Memo {
	return models.Memo{
		Content: helpers.SafeDereference(m.Content),
	}
}

// ValidateRequired verifies that the required fields for the request are provided.
func (m Memo) ValidateRequired(required ...int) error {
	for _, field := range required {
		switch field {
		case MemoFieldContent:
			if m.Content == nil {
				return errors.New("content is required")
			}
		}
	}
	return nil
}
