package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/memo-api/domain/repository"
	"github.com/akinolaemmanuel49/memo-api/infrastructure/database/postgres"
	"github.com/akinolaemmanuel49/memo-api/infrastructure/storage"
	"github.com/akinolaemmanuel49/memo-api/memo/api/helpers"
	"github.com/akinolaemmanuel49/memo-api/memo/api/internal"
	"github.com/akinolaemmanuel49/memo-api/memo/api/routes"
)

// serveApp starts the server and handles its shutdown
func serveApp(config internal.Config, db *sql.DB) error {
	app := internal.Application{
		Config: config,
		Repositories: repository.Repositories{
			Users:  postgres.NewUserInfrastructure(db),
			Social: postgres.NewSocialInfrastructure(db),
			Memo:   postgres.NewMemoInfrastructure(db),
			File: storage.NewFileInfrastructure(
				config.Cloudinary.CloudName,
				config.Cloudinary.APIKey,
				config.Cloudinary.APISecret,
			),
		},
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      routes.Router(app),
		IdleTimeout:  helpers.IdleTimeout,
		ReadTimeout:  helpers.ReadTimeout,
		WriteTimeout: helpers.WriteTimeout,
	}

	// start server
	log.Printf("starting %s server on %s\n", app.Config.Env, srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	log.Println("server stopped")
	return nil
}
