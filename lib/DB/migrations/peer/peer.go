package peersMigration

var PeersSchema string = `
CREATE TABLE IF NOT EXISTS peers (
    device_id    TEXT PRIMARY KEY,     -- DeviceID from discovery payload
    ip           TEXT NOT NULL,        -- last known IP (e.g. "192.168.0.12")
    port         INTEGER NOT NULL,     -- TCP port the peer listens on
    device_name  TEXT,                 -- Metadata.DeviceName (nullable)
    status       TEXT NOT NULL,        -- "discovered", "connecting", "connected", "dead"
    protocol_id  TEXT,                 -- optional: protocol identifier/version
    metadata     TEXT,                 -- JSON-encoded extra metadata (nullable)
    last_seen    DATETIME NOT NULL,    -- last discovery ping time (ISO8601)
    created_at   DATETIME NOT NULL     -- when this peer was first seen
);
`
