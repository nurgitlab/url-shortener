package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"url-shortener/internal/lib/logger/sl"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open(
		"postgres", storagePath,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url (
		id SERIAL NOT NULL,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB Ping failed: %w", err)
	}

	stid, err := db.Prepare(`
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
`)
	if _, err := stid.Exec(); err != nil {
		fmt.Println("idx")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, err
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.postgresql.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(urlToSave, alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgresql.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", op, sl.Err(err))
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgresql.DeleteURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("HERE")
		return fmt.Errorf("%s: %w", op, err)
	}
	fmt.Println("HE1RE")

	stmt, err = s.db.Prepare("DELETE FROM url WHERE alias = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.Exec(alias); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(alias).Scan(&alias)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return fmt.Errorf("%s: %w", op, sl.Err(err))
}
