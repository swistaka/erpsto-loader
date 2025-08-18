package entities

import (
	"arm/erpsto-loader/internal/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Request struct {
	ID                    string             `json:"id"`
	Index                 string             `json:"index"`
	DocDate               string             `json:"docDate"`
	DocNumber             string             `json:"docNumber"`
	Status                bool               `json:"status"`
	StoId                 string             `json:"stoId"`
	WarehouseId           string             `json:"warehouseId"`
	ClientId              string             `json:"clientId"`
	ContractId            string             `json:"contractId"`
	CarId                 string             `json:"carId"`
	GarantWork            string             `json:"garantWork"`
	LegalPersonShortName  string             `json:"legalPersonShortName"`
	ClientIndividualName  string             `json:"clientIndividualName"`
	ClientLegalShortName  string             `json:"clientLegalShortName"`
	ContactPersonFio      string             `json:"contactPersonFio"`
	CarBrand              string             `json:"carBrand"`
	CarEngineCapacity     string             `json:"carEngineCapacity"`
	CarLicensePlateNumber string             `json:"carLicensePlateNumber"`
	CarManufactureYear    string             `json:"carManufactureYear"`
	CarMileage            string             `json:"carMileage"`
	CarModelName          string             `json:"carModelName"`
	CarModification       string             `json:"carModification"`
	CarVin                string             `json:"carVin"`
	CompleteStatusDate    string             `json:"completeStatusDate"`
	CloseStatusDate       string             `json:"closeStatusDate"`
	SumWork               string             `json:"sumWork"`
	SumParts              string             `json:"sumParts"`
	SumReq                string             `json:"sumReq"`
	CreatedUserShortFio   string             `json:"createdUserShortFio"`
	NormoTimeFact         string             `json:"normoTimeFact"`
	NormoTimePlan         string             `json:"normoTimePlan"`
	PercentVat            float64            `json:"percentVat"`
	SumReqWithoutVat      string             `json:"sumReqWithoutVat"`
	SumVat                string             `json:"sumVat"`
	ReasonForPetition     string             `json:"reasonForPetition"`
	Positions             []RequestPosition  `json:"positions"`
	Works                 []RequestWork      `json:"works"`
	PerformersRequest     []RequestPerformer `json:"performersRequest"`
}

type RequestPosition struct {
	ID              string  `json:"id"`
	ProductId       string  `json:"productId"`
	Amount          string  `json:"amount"`
	Price           string  `json:"price"`
	PriceWithoutVat string  `json:"priceWithoutVat"`
	SumWithoutVat   string  `json:"sumWithoutVat"`
	PercentVat      float64 `json:"percentVat"`
	SumVat          string  `json:"sumVat"`
	Sum             string  `json:"sum"`
	ClientParts     bool    `json:"clientParts"`
	UsedParts       bool    `json:"usedParts"`
}

type RequestWork struct {
	ID              string  `json:"id"`
	WorkName        string  `json:"workName"`
	Amount          string  `json:"amount"`
	Price           string  `json:"price"`
	PriceWithoutVat string  `json:"priceWithoutVat"`
	SumWithoutVat   string  `json:"sumWithoutVat"`
	PercentVat      float64 `json:"percentVat"`
	SumVat          string  `json:"sumVat"`
	Sum             string  `json:"sum"`
}

type RequestPerformer struct {
	ID                  string `json:"id"`
	WorkId              string `json:"workId"`
	Name                string `json:"name"`
	CostPerHour         string `json:"costPerHour"`
	TotalEarnings       string `json:"totalEarnings"`
	ProductivityPercent string `json:"productivityPercent"`
}

func SaveRequest(db *sql.DB, data []byte, entityHash string, usedProductIDs map[string]bool) error {
	var entity Request
	if err := json.Unmarshal(data, &entity); err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
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
		INSERT OR REPLACE INTO request (
			id, ext_index, doc_date, doc_number, status, sto_id, warehouse_id, 
			client_id, contract_id, car_id, garant_work, legal_person_short_name, 
			client_individual_name, client_legal_short_name, contact_person_fio, 
			car_brand, car_engine_capacity, car_license_plate_number, 
			car_manufacture_year, car_mileage, car_model_name, car_modification, 
			car_vin, complete_status_date, close_status_date, sum_work, sum_parts, 
			sum, created_user_short_fio, normo_time_fact, normo_time_plan, 
			percent_vat, sum_without_vat, sum_vat, reason_for_petition, 
			updated_at, is_updated
		) VALUES (
		 	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?, ?
		)`,
		entity.ID,
		entity.Index,
		entity.DocDate,
		entity.DocNumber,
		entity.Status,
		entity.StoId,
		entity.WarehouseId,
		entity.ClientId,
		entity.ContractId,
		entity.CarId,
		entity.GarantWork,
		entity.LegalPersonShortName,
		entity.ClientIndividualName,
		entity.ClientLegalShortName,
		entity.ContactPersonFio,
		entity.CarBrand,
		entity.CarEngineCapacity,
		entity.CarLicensePlateNumber,
		entity.CarManufactureYear,
		entity.CarMileage,
		entity.CarModelName,
		entity.CarModification,
		entity.CarVin,
		entity.CompleteStatusDate,
		entity.CloseStatusDate,
		entity.SumWork,
		entity.SumParts,
		entity.SumReq,
		entity.CreatedUserShortFio,
		entity.NormoTimeFact,
		entity.NormoTimePlan,
		entity.PercentVat,
		entity.SumReqWithoutVat,
		entity.SumVat,
		entity.ReasonForPetition,
		time.Now().Format(time.RFC3339),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to save request: %w", err)
	}

	// Удаляем старые строки ТЧ.Товары
	_, err = tx.Exec(`
		DELETE FROM request_position WHERE request_id = ?`,
		entity.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete request position: %w", err)
	}

	_, err = tx.Exec(`
		DELETE FROM request_work WHERE request_id = ?`,
		entity.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete request work: %w", err)
	}

	_, err = tx.Exec(`
		DELETE FROM request_performer WHERE request_id = ?`,
		entity.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete request performer: %w", err)
	}

	// Сохраняем данные - ТЧ.Товары
	for _, pos := range entity.Positions {

		_, err = tx.Exec(`
			INSERT OR REPLACE INTO request_position (
				id, request_id, product_id, amount, price, price_without_vat, 
				sum_without_vat, percent_vat, sum_vat, sum, client_parts, used_parts
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
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
			pos.ClientParts,
			pos.UsedParts,
		)
		if err != nil {
			return fmt.Errorf("failed to save request position: %w", err)
		}

		// заполняем карту используемых элементов Product
		usedProductIDs[pos.ProductId] = true
	}

	// Сохраняем данные - ТЧ.Работы
	for _, work := range entity.Works {

		_, err = tx.Exec(`
			INSERT OR REPLACE INTO request_work (
				id, request_id, work_name, amount, price, price_without_vat, 
				sum_without_vat, percent_vat, sum_vat, sum
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			work.ID,
			entity.ID,
			work.WorkName,
			work.Amount,
			work.Price,
			work.PriceWithoutVat,
			work.SumWithoutVat,
			work.PercentVat,
			work.SumVat,
			work.Sum,
		)
		if err != nil {
			return fmt.Errorf("failed to save request work: %w", err)
		}
	}

	// Сохраняем данные - ТЧ.Исполнители
	for _, performer := range entity.PerformersRequest {

		_, err = tx.Exec(`
			INSERT OR REPLACE INTO request_performer (
				id, request_id, work_id, name, cost_per_hour, 
				total_earnings, productivity_percent
			) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			performer.ID,
			entity.ID,
			performer.WorkId,
			performer.Name,
			performer.CostPerHour,
			performer.TotalEarnings,
			performer.ProductivityPercent,
		)
		if err != nil {
			return fmt.Errorf("failed to save request performer: %w", err)
		}
	}

	// Сохраняем хэш
	if err := storage.SaveEntityHash(tx, "request", entity.ID, entityHash); err != nil {
		return fmt.Errorf("failed to save request hash: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
