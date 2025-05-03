package main

import (
	"ivanjabrony/refstudy/cmd/app"
	"ivanjabrony/refstudy/cmd/config"
	"ivanjabrony/refstudy/cmd/initDB"
	"log"

	_ "ivanjabrony/refstudy/docs"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// @title           Refstudy API
// @version         1.0
// @description     Refstude managing API
func main() {
	cfg := config.New()

	db, err := initDB.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := initDB.RunMigrations(db, cfg.Database.Name, "file://migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	application := app.New(db)
	if err := application.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
