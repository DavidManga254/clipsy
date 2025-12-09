package appmetacontroller

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

const deviceIDKey = "device_id"

// GetDeviceID returns the current device ID if present.
// If not set, it returns sql.ErrNoRows.
func GetDeviceID(ctx context.Context, db *sql.DB) (string, error) {
	const selectSQL = `
        SELECT value
        FROM app_meta
        WHERE key = ?;
    `

	var value string
	err := db.QueryRowContext(ctx, selectSQL, deviceIDKey).Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}


func SetDeviceID(ctx context.Context, db *sql.DB, id string) error {
	const upsertSQL = `
        INSERT INTO app_meta (key, value, created_at)
        VALUES (?, ?, ?)
        ON CONFLICT(key) DO UPDATE SET
            value = excluded.value,
            created_at = excluded.created_at;
    `

	now := time.Now().UTC()

	_, err := db.ExecContext(
		ctx,
		upsertSQL,
		deviceIDKey,
		id,
		now,
	)
	return err
}

// GetOrCreateDeviceID fetches the device ID if it exists,
// otherwise generates a new one, persists it, and returns it.
func GetOrCreateDeviceID(ctx context.Context, db *sql.DB) (string, error) {
	id, err := GetDeviceID(ctx, db)
	if err == nil && id != "" {
		return id, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	newID, err := generateRandomID(32) 
	if err != nil {
		return "", err
	}

	if err := SetDeviceID(ctx, db, newID); err != nil {
		return "", err
	}

	return newID, nil
}

func generateRandomID(nBytes int) (string, error) {
	buf := make([]byte, nBytes/2)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
