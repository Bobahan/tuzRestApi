package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	database *sql.DB
}

func New(databasePath string) (*Storage, error) {
	const op = "storage.sqlite.New" // op - операция, подсказывает где произошла ошибка

	database, err := sql.Open("sqlite3", databasePath)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	stmt, err := database.Prepare(`
    CREATE TABLE IF NOT EXISTS url(
      id INTEGER PRIMARY KEY,
      alias TEXT NOT NULL UNIQUE,
      url TEXT NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
  `)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	storage := &Storage{
		database: database,
	}

	return storage, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.database.Prepare(`INSERT INTO url(url, alias) VALUES(?, ?)`)

	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)

	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("%s, failed to get last inserted id %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	row := s.database.QueryRow(`SELECT url FROM url WHERE alias = ?`, alias)

	var url string

	if err := row.Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}

		return "", fmt.Errorf("%s, %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.sqlite.DeleteUrl"

	stmt, err := s.database.Prepare(`DELETE FROM url WHERE alias = ?`)

	if err != nil {
		return fmt.Errorf("%s, failed to prepare query %w", op, err)
	}

	_, err = stmt.Exec(alias)

	if err != nil {
		return fmt.Errorf("%s, failed to exec query %w", op, err)
	}

	return nil
}
