package entities

import (
	"arm/erpsto-loader/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Client struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	ShortName   string `json:"shortName"`
	UNP         string `json:"unp"`
	FullName    string `json:"fullName"`
	Address     string `json:"address"`
	BankName    string `json:"bankName"`
	BIK         string `json:"bik"`
	BankAccount string `json:"bankAccount"`
}

// Cохранение сущности в базу данных с хешированием
func SaveClient(db *sql.DB, data []byte, entityHash string, codePage string) error {
	var entity Client
	if err := json.Unmarshal(data, &entity); err != nil {
		return fmt.Errorf("failed to unmarshal client: %w", err)
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
		INSERT OR REPLACE INTO client (
			id, type, short_name, unp, full_name, address, 
			bank_name, bik, bank_account, updated_at, is_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entity.ID,
		entity.Type,
		entity.ShortName,
		entity.UNP,
		entity.FullName,
		entity.Address,
		entity.BankName,
		entity.BIK,
		entity.BankAccount,
		time.Now().Format(time.RFC3339),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to save client: %w", err)
	}

	// Сохраняем хэш
	if err := storage.SaveEntityHash(tx, "client", entity.ID, entityHash); err != nil {
		return fmt.Errorf("failed to save client hash: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
