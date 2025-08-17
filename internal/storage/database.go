package storage

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// Подключение к SQLite
func NewSQLiteStorage(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

// Нужно ли обновлять сущность - проверяем хэша пакета в базе и хэш входного пакета
func CheckIfNeedsUpdate(db *sql.DB, objectType, id, newHash string) (bool, error) {
	var currentHash string
	err := db.QueryRow(`
		SELECT hash FROM entity_hashes WHERE object_type = ? AND entity_id = ?`,
		objectType,
		id,
	).Scan(&currentHash)

	if err == sql.ErrNoRows {
		return true, nil // Сущности нет в базе, нужно сохранить
	}
	if err != nil {
		return false, err
	}

	return currentHash != newHash, nil // Сравниваем хэши
}

// Записываем хэш новой записи в базу
func SaveEntityHash(tx *sql.Tx, objectType, id, hash string) error {
	_, err := tx.Exec(`
		INSERT OR REPLACE INTO entity_hashes 
		(object_type, entity_id, hash, updated_at)
		VALUES (?, ?, ?, ?)`,
		objectType,
		id,
		hash,
		time.Now(),
	)

	return err
}
