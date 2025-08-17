package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migration struct {
	Version int
	Name    string
	Content string
}

func runMigrations(db *sql.DB) error {
	// Проверяем существует ли таблица миграций
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Получаем список примененных миграций
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Загружаем все миграции из файловой системы
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Применяем недостающие миграции
	for _, mig := range migrations {
		if _, ok := applied[mig.Version]; !ok {
			if err := applyMigration(db, mig); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", mig.Version, err)
			}
		}
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	return err
}

func getAppliedMigrations(db *sql.DB) (map[int]struct{}, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]struct{})
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = struct{}{}
	}

	return applied, nil
}

func loadMigrations() ([]Migration, error) {
	files, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if file.IsDir() || path.Ext(file.Name()) != ".sql" {
			continue
		}

		// Парсим номер версии из имени файла (001_name.up.sql)
		parts := strings.Split(file.Name(), "_")
		if len(parts) < 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid migration filename format: %s", file.Name())
		}

		content, err := fs.ReadFile(migrationsFS, path.Join("migrations", file.Name()))
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    file.Name(),
			Content: string(content),
		})
	}

	// Сортируем миграции по версии
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func applyMigration(db *sql.DB, mig Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Выполняем SQL-запросы из миграции
	if _, err := tx.Exec(mig.Content); err != nil {
		return fmt.Errorf("migration %d failed: %w", mig.Version, err)
	}

	// Фиксируем факт применения миграции
	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (version, name) VALUES (?, ?)",
		mig.Version, mig.Name,
	); err != nil {
		return err
	}

	return tx.Commit()
}
