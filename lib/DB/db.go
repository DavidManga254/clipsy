package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	"github.com/DavidManga254/clipsy/lib/DB/migrations"
	"github.com/DavidManga254/clipsy/lib/logger"
)

var database *sql.DB

func InitDB(){
	var dbPath string = "./clipsy.db"

	dataSourceName := "file:" + dbPath + "?_pragma=busy_timeout(5000)"

	db, err := sql.Open("sqlite",dataSourceName)

	if err != nil {
        log.Fatalf("failed to open db: %v", err)
    }

	database = db


	if err := database.Ping(); err != nil {
        log.Fatalf("failed to ping db: %v", err)
    }

    if err := migrations.Migrate(database); err != nil {
        log.Fatalf("failed to migrate db: %v", err)
    }

    logger.Info("DB initialized at: %v", dbPath)
}