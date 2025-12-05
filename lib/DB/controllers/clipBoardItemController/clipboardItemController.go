package clipboarditemcontroller

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type Kind int


const (
    KindUnknown Kind = iota
    KindText
    KindImage
    KindFiles
)


type ClipboardItem struct {
    ID        string   
    Kind      Kind      
    Text      string    
    ImageData []byte    
    Files     []string 
    Source    string    // "local" or device ID in future
    CreatedAt time.Time
}

func New(db *sql.DB,item ClipboardItem, ctx context.Context) error {
	var insertClipboardItemSQL string = `
        INSERT INTO clipboard_items (
            id, kind, text, image_data, files, source, created_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?);
    `

	filesJson,err := json.Marshal(item.Files)

	if err != nil {
		return err
	}

	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}

	_,err = db.ExecContext(
		ctx,
		insertClipboardItemSQL,
		item.ID,
		item.Kind,
		item.Text,
		item.ImageData,
		filesJson,
		item.Source,
		item.CreatedAt,
	)

	if err != nil{
		return err
	}

	return nil
}

func Delete(db *sql.DB,itemId string, ctx context.Context) error {
	var deleteClipboardItemSQL string = `
        DELETE FROM clipboard_items WHERE id = ?
    `

	_, err := db.ExecContext(
		ctx,
		deleteClipboardItemSQL,
		itemId,
	)

	if err != nil{
		return err
	}

	return nil
}


type PageResult struct {
    Items      []ClipboardItem
    TotalCount int
    Page       int
    PageSize   int
}

func ListClipboardItems(ctx context.Context, db *sql.DB, page, pageSize int) (PageResult, error) {
	const listClipboardItemsSQL = `
		SELECT
			id,
			kind,
			text,
			image_data,
			files,
			source,
			created_at
		FROM clipboard_items
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?;
	`

	const countClipboardItemsSQL = `
		SELECT COUNT(*) FROM clipboard_items;
	`

    if page < 1 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 20
    }

    offset := (page - 1) * pageSize

    // 1) total count
    var total int
    if err := db.QueryRowContext(ctx, countClipboardItemsSQL).Scan(&total); err != nil {
        return PageResult{}, err
    }

    // 2) page data
    rows, err := db.QueryContext(ctx, listClipboardItemsSQL, pageSize, offset)
    if err != nil {
        return PageResult{}, err
    }
    defer rows.Close()

    items := []ClipboardItem{}

    for rows.Next() {
        var (
            id        string
            kindStr   string
            text      sql.NullString
            imageData []byte
            filesJSON sql.NullString
            source    string
            createdAtStr string
        )

        if err := rows.Scan(
            &id,
            &kindStr,
            &text,
            &imageData,
            &filesJSON,
            &source,
            &createdAtStr,
        ); err != nil {
            return PageResult{}, err
        }

        t, err := time.Parse(time.RFC3339, createdAtStr)
        if err != nil {
            // fallback: ignore parse error, use zero time
            t = time.Time{}
        }

        // decode files []string
        var files []string
        if filesJSON.Valid && filesJSON.String != "" {
            _ = json.Unmarshal([]byte(filesJSON.String), &files)
        }

        // convert kind string to  Kind type
        var kind Kind

        switch kindStr {
            case "0":
                kind = KindUnknown
            case "1":
			    kind = KindText
            case "2":
			    kind = KindImage
            case "3":
			    kind = KindFiles 
        }

        item := ClipboardItem{
            ID:        id,
            Kind:      kind,
            Text:      text.String, 
            ImageData: imageData,
            Files:     files,
            Source:    source,
            CreatedAt: t,
        }

        items = append(items, item)
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