package DB

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	"github.com/DavidManga254/clipsy/lib/db/migrations"
	"github.com/DavidManga254/clipsy/lib/logger"
)

var Database *sql.DB

func InitDB(){
	var dbPath string = "./clipsy.db"

	dataSourceName := "file:" + dbPath + "?_pragma=busy_timeout(5000)"

	db, err := sql.Open("sqlite",dataSourceName)

	if err != nil {
        log.Fatalf("failed to open db: %v", err)
    }

	Database = db


	if err := Database.Ping(); err != nil {
        log.Fatalf("failed to ping db: %v", err)
    }

    if err := migrations.Migrate(Database); err != nil {
        log.Fatalf("failed to migrate db: %v", err)
    }

    logger.Info("DB initialized at: %v", dbPath)
}