// Package handler implementation
package handler

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/foolxdev/url-shortener/internal/model"
	"github.com/foolxdev/url-shortener/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLModel struct {
	DB *pgxpool.Pool
}

type Application struct {
	URLModel *URLModel
}

var errExpired = errors.New("short URL has expired")

func NewApplication(db *pgxpool.Pool) *Application {
	return &Application{
		URLModel: &URLModel{
			DB: db,
		},
	}
}

func (app *Application) Insert(
	u model.UrlCreate,
) (model.Url, error) {
	ctx := context.Background()

	tx, err := app.URLModel.DB.Begin(ctx)
	if err != nil {
		return model.Url{}, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", "error", err)
		}
	}()

	// Ask PostgreSQL for the next ID instead of depending on a particular
	// sequence name. Sequence names differ across existing deployments and
	// PostgreSQL may use an identity column rather than a serial sequence.
	var id int64
	err = tx.QueryRow(ctx, `SELECT nextval(pg_get_serial_sequence('urls', 'id'))`).Scan(&id)
	if err != nil {
		return model.Url{}, err
	}

	hash, err := service.ShortenURL(id)
	if err != nil {
		return model.Url{}, err
	}

	var createdAt time.Time
	stmt := `
		INSERT INTO urls (id, name, link, hash, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at;
	`

	err = tx.QueryRow(ctx, stmt, id, u.Name, u.Link, hash, u.ExpiresAt).Scan(&createdAt)
	if err != nil {
		return model.Url{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.Url{}, err
	}

	return model.Url{
		ID:        id,
		Name:      u.Name,
		Link:      u.Link,
		Hash:      hash,
		CreatedAt: createdAt,
		ExpiresAt: u.ExpiresAt,
	}, nil
}

func (app *Application) GetHash(
	hash string,
) (string, error) {
	stmt := `SELECT link, expires_at FROM urls WHERE hash = $1`

	var (
		link      string
		expiresAt *time.Time
	)

	err := app.URLModel.DB.QueryRow(
		context.Background(),
		stmt,
		hash,
	).Scan(&link, &expiresAt)
	if err != nil {
		return "", err
	}
	if expiresAt != nil && !expiresAt.After(time.Now()) {
		return "", errExpired
	}

	return link, nil
}
