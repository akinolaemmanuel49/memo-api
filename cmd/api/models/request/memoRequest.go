package request

import (
	"errors"

	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"github.com/akinolaemmanuel49/memo-api/internal/helpers"
)

type TextMemo struct {
	Content *string `json:"content" validate:"omitempty"`
}

const (
	TextMemoFieldContent = iota
)

func (tm TextMemo) ToModel() models.Memo {
	return models.Memo{
		Content: helpers.SafeDereference(tm.Content),
	}
}

// ValidateRequired verifies that the required fields for the request are provided.
func (tm TextMemo) ValidateRequired(required ...int) error {
	for _, field := range required {
		switch field {
		case TextMemoFieldContent:
			if tm.Content == nil {
				return errors.New("content is required")
			}
		}
	}
	return nil
}

type ImageMemo struct {
	Caption *string `json:"caption" validate:"required"`
}

const (
	ImageMemoFieldCaption = iota
)

func (im ImageMemo) ToModel() models.Memo {
	return models.Memo{
		Caption: helpers.SafeDereference(im.Caption),
	}
}

// ValidateRequired verifies that the required fields for the request are provided.
func (im ImageMemo) ValidateRequired(required ...int) error {
	for _, field := range required {
		switch field {
		case ImageMemoFieldCaption:
			if im.Caption == nil {
				return errors.New("caption is required")
			}
		}
	}
	return nil
}

type AudioMemo struct {
	Caption    *string `json:"caption" validate:"required"`
	Transcript *string `json:"transcript" validate:"omitempty"`
}

const (
	AudioMemoFieldCaption = iota
)

func (am AudioMemo) ToModel() models.Memo {
	return models.Memo{
		Caption: helpers.SafeDereference(am.Caption),
	}
}

// ValidateRequired verifies that the required fields for the request are provided.
func (am AudioMemo) ValidateRequired(required ...int) error {
	for _, field := range required {
		switch field {
		case AudioMemoFieldCaption:
			if am.Caption == nil {
				return errors.New("caption is required")
			}
		}
	}
	return nil
}

type VideoMemo struct {
	Caption    *string `json:"caption" validate:"required"`
	Transcript *string `json:"transcript" validate:"omitempty"`
}

const (
	VideoMemoFieldCaption = iota
)

func (vm VideoMemo) ToModel() models.Memo {
	return models.Memo{
		Caption: helpers.SafeDereference(vm.Caption),
	}
}

// ValidateRequired verifies that the required fields for the request are provided.
func (vm VideoMemo) ValidateRequired(required ...int) error {
	for _, field := range required {
		switch field {
		case VideoMemoFieldCaption:
			if vm.Caption == nil {
				return errors.New("caption is required")
			}
		}
	}
	return nil
}
