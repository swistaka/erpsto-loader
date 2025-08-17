package entities

import (
	"arm/erpsto-loader/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Car struct {
	ID                 string `json:"id"`
	ClientId           string `json:"clientId"`
	Brand              string `json:"brand"`
	Model              string `json:"model"`
	Color              string `json:"color"`
	Year               string `json:"year"`
	LicensePlateNumber string `json:"licensePlateNumber"`
	EngineCode         string `json:"engineCode"`
	VIN                string `json:"vin"`
	Mileage            string `json:"mileage"`
	Comment            string `json:"comment"`
}

// Cохранение сущности в базу данных с хешированием
func SaveCar(db *sql.DB, data []byte, entityHash string, codePage string) error {
	var entity Car
	if err := json.Unmarshal(data, &entity); err != nil {
		return fmt.Errorf("failed to unmarshal car: %w", err)
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
		INSERT OR REPLACE INTO car (
			id, client_id, brand, model, color, year, 
			license_plate_number, engine_code, vin, 
			mileage, comment, updated_at, is_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entity.ID,
		entity.ClientId,
		entity.Brand,
		entity.Model,
		entity.Color,
		entity.Year,
		entity.LicensePlateNumber,
		entity.EngineCode,
		entity.VIN,
		entity.Mileage,
		entity.Comment,
		time.Now().Format(time.RFC3339),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to save car: %w", err)
	}

	// Сохраняем хэш
	if err := storage.SaveEntityHash(tx, "car", entity.ID, entityHash); err != nil {
		return fmt.Errorf("failed to save car hash: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
