package entities

import (
	"arm/erpsto-loader/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Contract struct {
	ID       string `json:"id"`
	ClientId string `json:"clientId"`
	Number   string `json:"number"`
	Date     string `json:"date"`
}

// Cохранение сущности в базу данных с хешированием
func SaveContract(db *sql.DB, data []byte, entityHash string) error {
	var entity Contract
	if err := json.Unmarshal(data, &entity); err != nil {
		return fmt.Errorf("failed to unmarshal contract: %w", err)
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
		INSERT OR REPLACE INTO contract (
			id, client_id, name, number, date, updated_at, is_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		entity.ID,
		entity.ClientId,
		fmt.Sprintf("%s от %s г.", entity.Number, entity.Date),
		entity.Number,
		entity.Date,
		time.Now().Format(time.RFC3339),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to save contract: %w", err)
	}

	// Сохраняем хэш
	if err := storage.SaveEntityHash(tx, "contract", entity.ID, entityHash); err != nil {
		return fmt.Errorf("failed to save contract hash: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
