package main

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrDoesNotExist = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS kudos(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT NOT NULL UNIQUE,
        counter INTEGER NOT NULL
    );
    `

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(url string) (*Kudos, error) {
	res, err := r.db.Exec("INSERT INTO kudos(url, counter) values(?,0)", url)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	kudos := Kudos{ID: id, Url: url, Counter: 0}

	return &kudos, nil
}

func (r *SQLiteRepository) GetByUrl(url string) (*Kudos, error) {
	row := r.db.QueryRow("SELECT * FROM kudos WHERE url = ?", url)

	var kudos Kudos
	if err := row.Scan(&kudos.ID, &kudos.Url, &kudos.Counter); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDoesNotExist
		}
		return nil, err
	}
	return &kudos, nil
}

func (r *SQLiteRepository) IncreaseCounterByUrl(url string) (*Kudos, error) {
	kudos, err := r.GetByUrl(url)

	if err != nil && err != ErrDoesNotExist {
		return nil, err
	}

	if err == ErrDoesNotExist {
		kudos, err = r.Create(url)
	}

	if err != nil {
		return nil, err
	}

	res, err := r.db.Exec("UPDATE kudos SET counter = counter + 1 WHERE url = ?", url)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	kudos.Counter++

	return kudos, nil
}
