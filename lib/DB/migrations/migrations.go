package migrations

import (
	"database/sql"

	clipboardbItemMigration "github.com/DavidManga254/clipsy/lib/DB/migrations/clipboardItem"
	peersMigration "github.com/DavidManga254/clipsy/lib/DB/migrations/peer"
)

func Migrate(db *sql.DB) error {
	
	//clipboardItem table
	_,err := db.Exec(clipboardbItemMigration.ClipboardItemSchema)
	if err != nil {
		return err
	}

	//peers 
	_,err = db.Exec(peersMigration.PeersSchema)
	if err != nil {
		return err
	}

	return err
}