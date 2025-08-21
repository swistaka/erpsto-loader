package entities

import (
	"arm/erpsto-loader/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Refund struct {
	ID         string                `json:"id"`
	Index      string                `json:"index"`
	DocDate    string                `json:"docDate"`
	DocNumber  string                `json:"docNumber"`
	Status     bool                  `json:"status"`
	StoId      string                `json:"stoId"`
	ClientId   string                `json:"clientId"`
	Sum        string                `json:"sum"`
	SumVat     string                `json:"sumVat"`
	Comment    string                `json:"comment"`
	Positions  []RefundPosition `json:"positions"`
}

type RefundPosition struct {
	ID              string `json:"id"`
	ProductId       string `json:"productId"`
	Amount          string `json:"amount"`
	Price           string `json:"price"`
	PriceWithoutVat string `json:"priceWithoutVat"`
	SumWithoutVat   string `json:"sumWithoutVat"`
	PercentVat      string `json:"percentVat"`
	SumVat          string `json:"sumVat"`
	Sum             string `json:"sum"`
}

func SaveRefund(db *sql.DB, data []byte, entityHash string, usedProductIDs map[string]bool) error {
	var entity Refund
	if err := json.Unmarshal(data, &entity); err != nil {
		return fmt.Errorf("failed to unmarshal refund: %w", err)
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
		INSERT OR REPLACE INTO refund (
			id, ext_index, doc_date, doc_number, status, sto_id, 
			client_id, sum, sum_vat, comment, 
			updated_at, is_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entity.ID,
		entity.Index,
		entity.DocDate,
		entity.DocNumber,
		entity.Status,
		entity.StoId,
		entity.ClientId,
		entity.Sum,
		entity.SumVat,
		entity.Comment,
		time.Now().Format(time.RFC3339),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to save refund: %w", err)
	}

	// Удаляем старые строки ТЧ.Товары
	_, err = tx.Exec(`
		DELETE FROM refund_position WHERE refund_id = ?`,
		entity.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete refund position: %w", err)
	}

	// Сохраняем данные - ТЧ.Товары
	for _, pos := range entity.Positions {

		_, err = tx.Exec(`
			INSERT OR REPLACE INTO refund_position (
				id, refund_id, product_id, amount, price, price_without_vat, 
				sum_without_vat, percent_vat, sum_vat, sum
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			pos.ID,
			entity.ID,
			pos.ProductId,
			pos.Amount,
			pos.Price,
			pos.PriceWithoutVat,
			pos.SumWithoutVat,
			pos.PercentVat,
			pos.SumVat,
			pos.Sum,
		)
		if err != nil {
			return fmt.Errorf("failed to save refund position: %w", err)
		}

		// заполняем карту используемых элементов Product
		usedProductIDs[pos.ProductId] = true
	}

	// Сохраняем хэш
	if err := storage.SaveEntityHash(tx, "refund", entity.ID, entityHash); err != nil {
		return fmt.Errorf("failed to save refund hash: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
