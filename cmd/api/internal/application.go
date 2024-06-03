package internal

import "github.com/akinolaemmanuel49/memo-api/domain/repository"

// Application is a container for data required at different points throughout the server.
type Application struct {
	Config       Config
	Repositories repository.Repositories
}
