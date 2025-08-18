package entities

import (
	"arm/erpsto-loader/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Inventory struct {
	ID          string              `json:"id"`
	Index       string              `json:"index"`
	DocDate     string              `json:"docDate"`
	DocNumber   string              `json:"docNumber"`
	Status      bool                `json:"status"`
	StoId       string              `json:"stoId"`
	WarehouseId string              `json:"warehouseId"`
	Sum         string              `json:"sum"`
	Comment     string              `json:"comment"`
	Positions   []InventoryPosition `json:"positions"`
}

type InventoryPosition struct {
	ID         string `json:"id"`
	ProductId  string `json:"productId"`
	Amount     string `json:"amount"`
	AmountPlan string `json:"amountPlan"`
	Price      string `json:"price"`
	Sum        string `json:"sum"`
}

func SaveInventory(db *sql.DB, data []byte, entityHash string, usedProductIDs map[string]bool) error {
	var entity Inventory
	if err := json.Unmarshal(data, &entity); err != nil {
		return fmt.Errorf("failed to unmarshal inventory: %w", err)
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
		INSERT OR REPLACE INTO inventory (
			id, ext_index, doc_date, doc_number, status, sto_id, 
			warehouse_id, sum, comment, updated_at, is_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entity.ID,
		entity.Index,
		entity.DocDate,
		entity.DocNumber,
		entity.Status,
		entity.StoId,
		entity.WarehouseId,
		entity.Sum,
		entity.Comment,
		time.Now().Format(time.RFC3339),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to save inventory: %w", err)
	}

	// Удаляем старые строки ТЧ.Товары
	_, err = tx.Exec(`
		DELETE FROM inventory_position WHERE inventory_id = ?`,
		entity.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete inventory position: %w", err)
	}

	// Сохраняем данные - ТЧ.Товары
	for _, pos := range entity.Positions {

		_, err = tx.Exec(`
			INSERT OR REPLACE INTO inventory_position (
				id, inventory_id, product_id, amount, amount_plan, price, sum
			) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			pos.ID,
			entity.ID,
			pos.ProductId,
			pos.Amount,
			pos.AmountPlan,
			pos.Price,
			pos.Sum,
		)
		if err != nil {
			return fmt.Errorf("failed to save inventory position: %w", err)
		}

		// заполняем карту используемых элементов Product
		usedProductIDs[pos.ProductId] = true
	}

	// Сохраняем хэш
	if err := storage.SaveEntityHash(tx, "inventory", entity.ID, entityHash); err != nil {
		return fmt.Errorf("failed to save inventory hash: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
