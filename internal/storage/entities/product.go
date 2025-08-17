package entities

import (
	"arm/erpsto-loader/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Product struct {
	ID          string  `json:"id"`
	ArtId       string  `json:"artId"`
	PIN         string  `json:"pin"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	PercentVat  float64 `json:"percentVat"`
	TtnDate     string  `json:"ttnDate"`
	TtnNumber   string  `json:"ttnNumber"`
	Cost        string  `json:"cost"`
}

// Cохранение сущности в базу данных с хешированием
func SaveProduct(db *sql.DB, data []byte, entityHash string, codePage string) error {
	var entity Product
	if err := json.Unmarshal(data, &entity); err != nil {
		return fmt.Errorf("failed to unmarshal product: %w", err)
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
		INSERT OR REPLACE INTO product (
			id, art_id, pin, name, description, unit, 
			percent_vat, ttn_date, ttn_number, cost, 
			updated_at, is_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entity.ID,
		entity.ArtId,
		entity.PIN,
		entity.Name,
		entity.Description,
		entity.Unit,
		entity.PercentVat,
		entity.TtnDate,
		entity.TtnNumber,
		entity.Cost,
		time.Now().Format(time.RFC3339),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to save product: %w", err)
	}

	// Сохраняем хэш
	if err := storage.SaveEntityHash(tx, "product", entity.ID, entityHash); err != nil {
		return fmt.Errorf("failed to save product hash: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
