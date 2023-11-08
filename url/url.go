package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"encore.dev/storage/sqldb"
)

type URL struct {
	ID  string
	URL string
}

type ShortenParams struct {
	URL string
}

//encore:api public method=POST path=/url
func Shorten(ctx context.Context, p *ShortenParams) (*URL, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	} else if err := insert(ctx, id, p.URL); err != nil {
		return nil, err
	}
	return &URL{ID: id, URL: p.URL}, nil
}

func generateID() (string, error) {
	var data [7]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

func insert(ctx context.Context, id, url string) error {
	_, err := sqldb.Exec(ctx, `
        INSERT INTO url (id, original_url)
        SELECT $1, $2 
        WHERE NOT EXISTS (
        SELECT original_url FROM url WHERE url.original_url != $3
        )
    `, id, url, url)
	return err
}

//encore:api public method=GET path=/url/:id
func Get(ctx context.Context, id string) (*URL, error) {
	u := &URL{ID: id}
	err := sqldb.QueryRow(ctx, `
        SELECT original_url FROM url
        WHERE id = $1
    `, id).Scan(&u.URL)
	return u, err
}
