package migrations

import (
	"database/sql"

	appmetamigration "github.com/DavidManga254/clipsy/lib/DB/migrations/appMeta"
	clipboardbItemMigration "github.com/DavidManga254/clipsy/lib/DB/migrations/clipboardItem"
	peersMigration "github.com/DavidManga254/clipsy/lib/DB/migrations/peer"
)

func Migrate(db *sql.DB) error {

	//app meta migration 
	_,err := db.Exec(appmetamigration.AppMetaSchema)
	if err != nil {
		return err
	}
	
	//clipboardItem table
	_,err = db.Exec(clipboardbItemMigration.ClipboardItemSchema)
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