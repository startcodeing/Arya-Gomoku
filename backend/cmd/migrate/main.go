package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gomoku-backend/internal/config"

	_ "github.com/lib/pq"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/main.go [up|down|status|reset]")
		os.Exit(1)
	}

	command := os.Args[1]

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// 连接数据库
	var dsn, driver string
	if cfg.Database.Type == "postgres" {
		driver = "postgres"
		dsn = cfg.Database.GetDSN()
	} else {
		// SQLite
		driver = "sqlite"
		dsn = cfg.Database.DBName
	}
	
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	migrator := &Migrator{db: db, cfg: cfg}

	switch command {
	case "up":
		if err := migrator.MigrateUp(); err != nil {
			log.Fatal("Migration up failed:", err)
		}
		fmt.Println("Migration up completed successfully")
	case "down":
		if err := migrator.MigrateDown(); err != nil {
			log.Fatal("Migration down failed:", err)
		}
		fmt.Println("Migration down completed successfully")
	case "status":
		if err := migrator.Status(); err != nil {
			log.Fatal("Failed to get migration status:", err)
		}
	case "reset":
		if err := migrator.Reset(); err != nil {
			log.Fatal("Migration reset failed:", err)
		}
		fmt.Println("Database reset completed successfully")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: up, down, status, reset")
		os.Exit(1)
	}
}

type Migrator struct {
	db  *sql.DB
	cfg *config.Config
}

func (m *Migrator) ensureMigrationsTable() error {
	var query string
	if m.cfg.Database.Type == "sqlite" {
		query = `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version TEXT PRIMARY KEY,
				applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`
	} else {
		query = `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version VARCHAR(255) PRIMARY KEY,
				applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)
		`
	}
	_, err := m.db.Exec(query)
	return err
}

func (m *Migrator) getMigrationFiles() ([]string, error) {
	migrationsDir := "migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("获取迁移文件失败: %w", err)
	}

	var migrationFiles []string
	
	// 如果是SQLite数据库，优先使用SQLite兼容的迁移文件
	if m.cfg.Database.Type == "sqlite" {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), "_sqlite.sql") {
				migrationFiles = append(migrationFiles, file.Name())
			}
		}
		// 如果没有SQLite专用文件，则使用通用文件
		if len(migrationFiles) == 0 {
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") && !strings.Contains(file.Name(), "_sqlite") {
					migrationFiles = append(migrationFiles, file.Name())
				}
			}
		}
	} else {
		// 对于其他数据库，使用非SQLite专用的文件
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") && !strings.Contains(file.Name(), "_sqlite") {
				migrationFiles = append(migrationFiles, file.Name())
			}
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	if err := m.ensureMigrationsTable(); err != nil {
		return nil, err
	}

	rows, err := m.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}

func (m *Migrator) MigrateUp() error {
	migrations, err := m.getMigrationFiles()
	if err != nil {
		return err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		version := strings.TrimSuffix(migration, ".sql")
		if applied[version] {
			fmt.Printf("Migration %s already applied, skipping\n", migration)
			continue
		}

		fmt.Printf("Applying migration %s...\n", migration)
		if err := m.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %v", migration, err)
		}

		// 记录迁移
		_, err := m.db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %v", migration, err)
		}

		fmt.Printf("Migration %s applied successfully\n", migration)
	}

	return nil
}

func (m *Migrator) MigrateDown() error {
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	migrations, err := m.getMigrationFiles()
	if err != nil {
		return err
	}

	// 反向处理迁移
	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]
		version := strings.TrimSuffix(migration, ".sql")
		
		if !applied[version] {
			fmt.Printf("Migration %s not applied, skipping\n", migration)
			continue
		}

		fmt.Printf("Rolling back migration %s...\n", migration)
		
		// 对于简单的迁移，我们只是删除记录
		// 在生产环境中，你可能需要创建专门的回滚脚本
		_, err := m.db.Exec("DELETE FROM schema_migrations WHERE version = $1", version)
		if err != nil {
			return fmt.Errorf("failed to rollback migration %s: %v", migration, err)
		}

		fmt.Printf("Migration %s rolled back successfully\n", migration)
		break // 只回滚最后一个迁移
	}

	return nil
}

func (m *Migrator) Status() error {
	migrations, err := m.getMigrationFiles()
	if err != nil {
		return err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")
	
	for _, migration := range migrations {
		version := strings.TrimSuffix(migration, ".sql")
		status := "Pending"
		if applied[version] {
			status = "Applied"
		}
		fmt.Printf("%-30s %s\n", migration, status)
	}

	return nil
}

func (m *Migrator) Reset() error {
	fmt.Println("WARNING: This will drop all tables and data!")
	fmt.Print("Are you sure you want to continue? (yes/no): ")
	
	var response string
	fmt.Scanln(&response)
	
	if strings.ToLower(response) != "yes" {
		fmt.Println("Reset cancelled")
		return nil
	}

	// 删除所有表
	tables := []string{
		"llm_configs",
		"pvp_room_players", 
		"pvp_games",
		"llm_games",
		"ai_games",
		"user_statistics",
		"user_sessions",
		"users",
		"schema_migrations",
	}

	for _, table := range tables {
		_, err := m.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %v", table, err)
		}
		fmt.Printf("Dropped table %s\n", table)
	}

	// 删除函数
	_, err := m.db.Exec("DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE")
	if err != nil {
		return fmt.Errorf("failed to drop function: %v", err)
	}

	fmt.Println("Database reset completed")
	return nil
}

func (m *Migrator) applyMigration(filename string) error {
	filepath := filepath.Join("migrations", filename)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %v", filename, err)
	}

	// 执行SQL
	_, err = m.db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %v", err)
	}

	return nil
}