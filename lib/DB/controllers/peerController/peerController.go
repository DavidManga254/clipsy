package peerscontroller

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type PeerStatus string

const (
	PeerStatusDiscovered PeerStatus = "discovered"
	PeerStatusConnecting PeerStatus = "connecting"
	PeerStatusConnected  PeerStatus = "connected"
	PeerStatusDead       PeerStatus = "dead"
)

type Peer struct {
	DeviceID   string
	IP         string
	Port       int
	DeviceName string
	Status     PeerStatus
	ProtocolID string
	Metadata   map[string]string
	LastSeen   time.Time
	CreatedAt  time.Time
	Is_trusted int
}

func New(db *sql.DB, p Peer, ctx context.Context) error {
	const insertPeerSQL = `
        INSERT INTO peers (
            device_id,
            ip,
            port,
            device_name,
            status,
            protocol_id,
            metadata,
            last_seen,
            created_at,
			is_trusted
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,?);
    `

	metadataJSON, err := json.Marshal(p.Metadata)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.LastSeen.IsZero() {
		p.LastSeen = now
	}

	_, err = db.ExecContext(
		ctx,
		insertPeerSQL,
		p.DeviceID,
		p.IP,
		p.Port,
		p.DeviceName,
		string(p.Status),
		p.ProtocolID,
		string(metadataJSON),
		p.LastSeen,
		p.CreatedAt,
		p.Is_trusted,
	)
	if err != nil {
		return err
	}

	return nil
}

func Delete(db *sql.DB, deviceID string, ctx context.Context) error {
	const deletePeerSQL = `
        DELETE FROM peers WHERE device_id = ?;
    `

	_, err := db.ExecContext(
		ctx,
		deletePeerSQL,
		deviceID,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetPeerById(db *sql.DB, deviceID string, ctx context.Context) (Peer,error) {
	var getPeerSQL string = `SELECT device_id, ip, port, device_name, status, protocol_id, metadata, last_seen, created_at FROM peers WHERE device_id = ? LIMIT 1;`

	var (
		id           string
		ip           string
		port         int
		deviceName   sql.NullString
		statusStr    string
		protocolID   sql.NullString
		metadataJSON sql.NullString
		lastSeenStr  string
		createdAtStr string
		is_trusted int
	)

	err := db.QueryRowContext(ctx, getPeerSQL, deviceID).Scan(
		&id,
		&ip,
		&port,
		&deviceName,
		&statusStr,
		&protocolID,
		&metadataJSON,
		&lastSeenStr,
		&createdAtStr,
		&is_trusted,
	)
	if err != nil {
		return Peer{}, err
	}

	lastSeen, err := time.Parse(time.RFC3339, lastSeenStr)
	if err != nil {
		lastSeen = time.Time{}
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		createdAt = time.Time{}
	}

	metadata := map[string]string{}
	if metadataJSON.Valid && metadataJSON.String != "" {
		_ = json.Unmarshal([]byte(metadataJSON.String), &metadata)
	}

	p := Peer{
		DeviceID:   id,
		IP:         ip,
		Port:       port,
		DeviceName: deviceName.String,
		Status:     PeerStatus(statusStr),
		ProtocolID: protocolID.String,
		Metadata:   metadata,
		LastSeen:   lastSeen,
		CreatedAt:  createdAt,
		Is_trusted: is_trusted,
	}

	return p, nil

}

type PageResult struct {
	Items      []Peer
	TotalCount int
	Page       int
	PageSize   int
}

func ListPeers(ctx context.Context, db *sql.DB, page, pageSize int) (PageResult, error) {
	const listPeersSQL = `
		SELECT
			device_id,
			ip,
			port,
			device_name,
			status,
			protocol_id,
			metadata,
			last_seen,
			created_at,
			is_trusted
		FROM peers
		ORDER BY last_seen DESC
		LIMIT ? OFFSET ?;
	`

	const countPeersSQL = `
		SELECT COUNT(*) FROM peers;
	`

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var total int
	if err := db.QueryRowContext(ctx, countPeersSQL).Scan(&total); err != nil {
		return PageResult{}, err
	}

	rows, err := db.QueryContext(ctx, listPeersSQL, pageSize, offset)
	if err != nil {
		return PageResult{}, err
	}
	defer rows.Close()

	items := []Peer{}

	for rows.Next() {
		var (
			deviceID      string
			ip            string
			port          int
			deviceName    sql.NullString
			statusStr     string
			protocolID    sql.NullString
			metadataJSON  sql.NullString
			lastSeenStr   string
			createdAtStr  string
			is_trusted int
		)

		if err := rows.Scan(
			&deviceID,
			&ip,
			&port,
			&deviceName,
			&statusStr,
			&protocolID,
			&metadataJSON,
			&lastSeenStr,
			&createdAtStr,
			&is_trusted,
		); err != nil {
			return PageResult{}, err
		}

		lastSeen, err := time.Parse(time.RFC3339, lastSeenStr)
		if err != nil {
			lastSeen = time.Time{}
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			createdAt = time.Time{}
		}

		metadata := map[string]string{}
		if metadataJSON.Valid && metadataJSON.String != "" {
			_ = json.Unmarshal([]byte(metadataJSON.String), &metadata)
		}

		p := Peer{
			DeviceID:   deviceID,
			IP:         ip,
			Port:       port,
			DeviceName: deviceName.String,
			Status:     PeerStatus(statusStr),
			ProtocolID: protocolID.String,
			Metadata:   metadata,
			LastSeen:   lastSeen,
			CreatedAt:  createdAt,
			Is_trusted: is_trusted,
		}

		items = append(items, p)
	}

	if err := rows.Err(); err != nil {
		return PageResult{}, err
	}

	return PageResult{
		Items:      items,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
