package entities

import (
	"arm/erpsto-loader/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Moving struct {
	ID              string           `json:"id"`
	Index           string           `json:"index"`
	DocDate         string           `json:"docDate"`
	DocNumber       string           `json:"docNumber"`
	Status          bool             `json:"status"`
	StoFromId       string           `json:"stoFromId"`
	StoToId         string           `json:"stoToId"`
	WarehouseFromId string           `json:"warehouseFromId"`
	WarehouseToId   string           `json:"warehouseToId"`
	Sum             string           `json:"sum"`
	Comment         string           `json:"comment"`
	Positions       []MovingPosition `json:"positions"`
}

type MovingPosition struct {
	ID        string `json:"id"`
	ProductId string `json:"productId"`
	Amount    string `json:"amount"`
	Price     string `json:"price"`
	Sum       string `json:"sum"`
}

func SaveMoving(db *sql.DB, data []byte, entityHash string, usedProductIDs map[string]bool) error {
	var entity Moving
	if err := json.Unmarshal(data, &entity); err != nil {
		return fmt.Errorf("failed to unmarshal moving: %w", err)
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

	// Сохраняем данные - шапка
	_, err = tx.Exec(`
		INSERT OR REPLACE INTO moving (
			id, ext_index, doc_date, doc_number, status, sto_from_id, 
			sto_to_id, warehouse_from_id, warehouse_to_id, sum, comment, 
			updated_at, is_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entity.ID,
		entity.Index,
		entity.DocDate,
		entity.DocNumber,
		entity.Status,
		entity.StoFromId,
		entity.StoToId,
		entity.WarehouseFromId,
		entity.WarehouseToId,
		entity.Sum,
		entity.Comment,
		time.Now().Format(time.RFC3339),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to save moving: %w", err)
	}

	// Удаляем старые строки ТЧ.Товары
	_, err = tx.Exec(`
		DELETE FROM moving_position WHERE moving_id = ?`,
		entity.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete moving position: %w", err)
	}

	// Сохраняем данные - ТЧ.Товары
	for _, pos := range entity.Positions {

		_, err = tx.Exec(`
			INSERT OR REPLACE INTO moving_position (
				id, moving_id, product_id, amount, price, sum
			) VALUES (?, ?, ?, ?, ?, ?)`,
			pos.ID,
			entity.ID,
			pos.ProductId,
			pos.Amount,
			pos.Price,
			pos.Sum,
		)
		if err != nil {
			return fmt.Errorf("failed to save moving position: %w", err)
		}

		// заполняем карту используемых элементов Product
		usedProductIDs[pos.ProductId] = true
	}

	// Сохраняем хэш
	if err := storage.SaveEntityHash(tx, "moving", entity.ID, entityHash); err != nil {
		return fmt.Errorf("failed to save moving hash: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
