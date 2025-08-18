package entities

import (
	"arm/erpsto-loader/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Dismantling struct {
	ID          string                `json:"id"`
	Index       string                `json:"index"`
	DocDate     string                `json:"docDate"`
	DocNumber   string                `json:"docNumber"`
	Status      bool                  `json:"status"`
	StoId       string                `json:"stoId"`
	WarehouseId string                `json:"warehouseId"`
	ProductId   string                `json:"productId"`
	Amount      string                `json:"amount"`
	Sum         string                `json:"sum"`
	Comment     string                `json:"comment"`
	Positions   []DismantlingPosition `json:"positions"`
}

type DismantlingPosition struct {
	ID        string `json:"id"`
	ProductId string `json:"productId"`
	Amount    string `json:"amount"`
	Price     string `json:"price"`
	Sum       string `json:"sum"`
}

func SaveDismantling(db *sql.DB, data []byte, entityHash string, usedProductIDs map[string]bool) error {
	var entity Dismantling
	if err := json.Unmarshal(data, &entity); err != nil {
		return fmt.Errorf("failed to unmarshal dismantling: %w", err)
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
		INSERT OR REPLACE INTO dismantling (
			id, ext_index, doc_date, doc_number, status, sto_id, 
			warehouse_id, product_id, amount, sum, comment, 
			updated_at, is_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entity.ID,
		entity.Index,
		entity.DocDate,
		entity.DocNumber,
		entity.Status,
		entity.StoId,
		entity.WarehouseId,
		entity.ProductId,
		entity.Amount,
		entity.Sum,
		entity.Comment,
		time.Now().Format(time.RFC3339),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to save dismantling: %w", err)
	}

	// заполняем карту используемых элементов Product
	usedProductIDs[entity.ProductId] = true

	// Удаляем старые строки ТЧ.Товары
	_, err = tx.Exec(`
		DELETE FROM dismantling_position WHERE dismantling_id = ?`,
		entity.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete dismantling position: %w", err)
	}

	// Сохраняем данные - ТЧ.Товары
	for _, pos := range entity.Positions {

		_, err = tx.Exec(`
			INSERT OR REPLACE INTO dismantling_position (
				id, dismantling_id, product_id, amount, price, sum
			) VALUES (?, ?, ?, ?, ?, ?)`,
			pos.ID,
			entity.ID,
			pos.ProductId,
			pos.Amount,
			pos.Price,
			pos.Sum,
		)
		if err != nil {
			return fmt.Errorf("failed to save dismantling position: %w", err)
		}

		// заполняем карту используемых элементов Product
		usedProductIDs[pos.ProductId] = true
	}

	// Сохраняем хэш
	if err := storage.SaveEntityHash(tx, "dismantling", entity.ID, entityHash); err != nil {
		return fmt.Errorf("failed to save dismantling hash: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
