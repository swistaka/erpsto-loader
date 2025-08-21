package processor

import (
	"arm/erpsto-loader/internal/api"
	"arm/erpsto-loader/internal/config"
	"arm/erpsto-loader/internal/logger"
	"arm/erpsto-loader/internal/storage"
	"arm/erpsto-loader/internal/storage/entities"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

func Process(cfg *config.Config, db *sql.DB, logger *logger.FileLogger) error {
	var filePath string
	var err error

	switch cfg.SourceType {
	case "file":
		filePath = cfg.FileName
	case "api":
		// Загружаем данные API во временный файл
		api_cli := api.NewAPIClient(cfg.ApiUrl, cfg.ApiKey)
		filePath, err = api_cli.DownloadToTempFile()
		if err != nil {
			return fmt.Errorf("API download failed: %w", err)
		}
		defer os.Remove(filePath)
	default:
		return fmt.Errorf("unsupported source type: %s", cfg.SourceType)
	}

	// Далее обработка как для файла
	file_data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("file read failed: %w", err)
	}

	// Парсинг JSON
	var input struct {
		DateFrom string            `json:"dateFrom"`
		DateTo   string            `json:"dateTo"`
		Data     []json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(file_data, &input); err != nil {
		return fmt.Errorf("error unmarshaling main packet: %w", err)
	}

	// product будем сохранять только те, которые встречаются в документах
	var rawEntitiesProduct []json.RawMessage
	usedProductIDs := make(map[string]bool)

	// Обработка каждой сущности
	for _, rawEntity := range input.Data {
		// считаем тип объекта и id - для дальнейшей маршрутизации обработки
		var baseEntity struct {
			ObjectType string `json:"objectType"`
			ID         string `json:"id"`
		}
		if err := json.Unmarshal(rawEntity, &baseEntity); err != nil {
			logger.LogError(fmt.Errorf("error unmarshaling entity: %w", err))
			continue
		}

		// Вычисляем хэш сущности
		entityHash := computeHash(rawEntity)

		// Проверяем, нужно ли обновлять сущность
		needsUpdate, err := storage.CheckIfNeedsUpdate(db, baseEntity.ObjectType, baseEntity.ID, entityHash)
		if err != nil {
			logger.LogError(fmt.Errorf("error checking entity update status: %w", err))
			continue
		}

		if !needsUpdate {
			continue // Сущность не изменилась, пропускаем
		}

		switch baseEntity.ObjectType {
		case "sto":
			err = entities.SaveSTO(db, rawEntity, entityHash)
		case "warehouse":
			err = entities.SaveWarehouse(db, rawEntity, entityHash)
		case "client":
			err = entities.SaveClient(db, rawEntity, entityHash)
		case "contract":
			err = entities.SaveContract(db, rawEntity, entityHash)
		case "car":
			err = entities.SaveCar(db, rawEntity, entityHash)
		case "product":
			rawEntitiesProduct = append(rawEntitiesProduct, rawEntity)
			//err = entities.SaveProduct(db, rawEntity, entityHash)
		case "invoice":
			err = entities.SaveInvoice(db, rawEntity, entityHash, usedProductIDs)
		case "realization":
			err = entities.SaveRealization(db, rawEntity, entityHash, usedProductIDs)
		case "refund":
			err = entities.SaveRefund(db, rawEntity, entityHash, usedProductIDs)
		case "moving":
			err = entities.SaveMoving(db, rawEntity, entityHash, usedProductIDs)
		case "dismantling":
			err = entities.SaveDismantling(db, rawEntity, entityHash, usedProductIDs)
		case "inventory":
			err = entities.SaveInventory(db, rawEntity, entityHash, usedProductIDs)
		case "request":
			err = entities.SaveRequest(db, rawEntity, entityHash, usedProductIDs)
		default:
			logger.LogError(fmt.Errorf("unknown entity type: %s", baseEntity.ObjectType))
			continue
		}

		if err != nil {
			logger.LogError(fmt.Errorf("save %s failed: %w", baseEntity.ObjectType, err))
		}
	}

	// вторым проходом - обработка product
	for _, rawEntity := range rawEntitiesProduct {
		// считаем тип объекта и id - для дальнейшей маршрутизации обработки
		var baseEntity struct {
			ObjectType string `json:"objectType"`
			ID         string `json:"id"`
		}
		if err := json.Unmarshal(rawEntity, &baseEntity); err != nil {
			logger.LogError(fmt.Errorf("error unmarshaling entity: %w", err))
			continue
		}

		if !usedProductIDs[baseEntity.ID] {
			continue
		}

		// Вычисляем хэш сущности
		entityHash := computeHash(rawEntity)

		switch baseEntity.ObjectType {
		case "product":
			err = entities.SaveProduct(db, rawEntity, entityHash)
		}

		if err != nil {
			logger.LogError(fmt.Errorf("save %s failed: %w", baseEntity.ObjectType, err))
		}
	}

	return nil
}

// computeHash вычисляет SHA256 хэш для данных сущности
func computeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
