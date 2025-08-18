package entities

import (
	"arm/erpsto-loader/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type STO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Index   string `json:"index"`
	Address string `json:"address"`
}

// Cохранение сущности в базу данных с хешированием
func SaveSTO(db *sql.DB, data []byte, entityHash string) error {
	var entity STO
	if err := json.Unmarshal(data, &entity); err != nil {
		return fmt.Errorf("failed to unmarshal sto: %w", err)
	}

	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Сохраняем данные
	_, err = tx.Exec(`
		INSERT OR REPLACE INTO sto (
			id, name, ext_index, address, updated_at, is_updated
		) VALUES (?, ?, ?, ?, ?, ?)`,
		entity.ID,
		entity.Name,
		entity.Index,
		entity.Address,
		time.Now().Format(time.RFC3339),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to save sto: %w", err)
	}

	// Сохраняем хэш
	if err := storage.SaveEntityHash(tx, "sto", entity.ID, entityHash); err != nil {
		return fmt.Errorf("failed to save sto hash: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
