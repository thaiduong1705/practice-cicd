package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"todolist/internal/models"
)

// ConnectDB uses GORM to open (and create if missing) a Postgres database
// It also runs an AutoMigrate for the Task model.
func ConnectDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh", host, user, password, dbName, port)

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres with gorm: %w", err)
	}

	// Initialize schema
	if err := gdb.AutoMigrate(&models.Task{}); err != nil {
		return nil, fmt.Errorf("auto-migrate: %w", err)
	}

	return gdb, nil
}

func CloseDB(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB from gorm.DB: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close sql.DB: %w", err)
	}
	return nil
}
