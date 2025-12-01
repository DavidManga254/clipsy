package clipboardbItemMigration

var ClipboardItemSchema string = `
	CREATE TABLE IF NOT EXISTS clipboard_items (
    id          TEXT PRIMARY KEY,          -- UUID or any string ID
    kind        TEXT NOT NULL,             -- or INTEGER if your Kind is int-backed
    text        TEXT,                      -- clipboard text (nullable)
    image_data  BLOB,                      -- raw image bytes (nullable)
    files       TEXT,                      -- JSON-encoded []string
    source      TEXT NOT NULL,             -- "local" or device id
    created_at  DATETIME NOT NULL          -- ISO8601 string (time.Time)
);
`