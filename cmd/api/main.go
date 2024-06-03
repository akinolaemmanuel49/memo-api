package main

import (
	"database/sql"
	"log"

	"github.com/akinolaemmanuel49/memo-api/memo/api/internal"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// load environment variables from dotenv file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading dotenv file: %s\n", err.Error())
	}

	// setup server configurations
	config := internal.Config{}
	config.Parse()
	if err := config.Validate(); err != nil {
		log.Fatalf("error validating configurations: %s\n\n", err)
	}

	// set Gin to release mode on production
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// open database connection
	db, err := openDB(config)
	if err != nil {
		log.Fatalf("database connection error: %s\n", err.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)
	log.Println("database connection established")

	// start server
	if err := serveApp(config, db); err != nil {
		log.Fatalln(err)
	}
}
