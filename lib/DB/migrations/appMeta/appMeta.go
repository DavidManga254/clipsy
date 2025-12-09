package appmetamigration

var AppMetaSchema string = `
CREATE TABLE IF NOT EXISTS app_meta (
    key         TEXT PRIMARY KEY,     -- e.g. "device_id"
    value       TEXT NOT NULL,        -- the actual value
    created_at  DATETIME NOT NULL     -- when this key was first set
);
`
