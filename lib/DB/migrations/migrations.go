package migrations

import (
	"database/sql"

	clipboardbItemMigration "github.com/DavidManga254/clipsy/lib/DB/migrations/clipboardItem"
)

func Migrate(db *sql.DB) error {
	
	//clipboardItem table
	_,err := db.Exec(clipboardbItemMigration.ClipboardItemSchema)
	if err != nil {
		return err
	}

	return err
}