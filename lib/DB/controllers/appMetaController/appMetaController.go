package appmetacontroller

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

const deviceIDKey = "device_id"

type DeviceDetails struct {
	deviceId string
	tcpPort int
}

// GetDeviceID returns the current device ID if present.
// If not set, it returns sql.ErrNoRows.
func GetDeviceDetails(ctx context.Context, db *sql.DB) (DeviceDetails, error) {
	const selectSQL = `
        SELECT value,tcpPort
        FROM app_meta
        WHERE key = ?;
    `

	var value string
	var tcpPort int

	err := db.QueryRowContext(ctx, selectSQL, deviceIDKey).Scan(&value,&tcpPort)
	if err != nil {
		return DeviceDetails{}, err
	}

	deviceDetails := DeviceDetails {
		deviceId: value,
		tcpPort: tcpPort,
	}



	return deviceDetails, nil
}

func SetDeviceTcpDetails(ctx context.Context,db *sql.DB,tcpPort int) error {
	const updateSql = `
		UPDATE app_meta
		SET 
			tcpPort = ?
		WHERE key = ?
	`

	_, err := db.ExecContext(
		ctx,
		updateSql,
		tcpPort,
		deviceIDKey,
	)

	return err
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
	deviceDetails, err := GetDeviceDetails(ctx, db)
	if err == nil && deviceDetails.deviceId != "" {
		return deviceDetails.deviceId, nil
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
