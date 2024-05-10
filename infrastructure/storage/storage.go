package storage

import (
	"github.com/akinolaemmanuel49/memo-api/domain/repository"
)

type file struct {
	CloudName string
	APIKey    string
	APISecret string
}

func NewFileInfrastructure(cloudName, apiKey, apiSecret string) repository.FileRepository {
	return file{
		CloudName: cloudName,
		APIKey:    apiKey,
		APISecret: apiSecret,
	}
}

const (
	TypeImage = "image"
)

var Invalidate = true
