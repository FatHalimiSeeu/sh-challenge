package config

import (
	"fmt"
	"inventory-api/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := os.Getenv("SQL_CONNECTION_STRING")
	if dsn == "" {
		log.Fatal("SQL_CONNECTION_STRING environment variable is not set")
	}

	dsnMaster := fmt.Sprintf("%s;database=master", dsn[:len(dsn)-len("inventory_db")])

	maxRetries := 20
	var tempDB *gorm.DB
	for i := 1; i <= maxRetries; i++ {
		tempDB, err = gorm.Open(sqlserver.Open(dsnMaster), &gorm.Config{})
		if err == nil {
			log.Println("Successfully connected to SQL Server (master)!")
			break
		}

		log.Printf("Failed to connect to SQL Server (attempt %d/%d): %v", i, maxRetries, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to SQL Server after multiple attempts:", err)
	}

	dbName := "inventory_db"
	if !databaseExists(tempDB, dbName) {
		createDatabase(tempDB, dbName)
	}

	dsnWithDB := fmt.Sprintf("%s;database=%s", dsn[:len(dsn)-len("inventory_db")], dbName)
	log.Printf("Connecting to inventory_db database: %s", dsnWithDB)

	DB, err = gorm.Open(sqlserver.Open(dsnWithDB), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the inventory_db database:", err)
	}

	if err := DB.AutoMigrate(&models.User{}, &models.InventoryItem{}, &models.Restock{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connection to inventory_db established successfully!")
}

func databaseExists(db *gorm.DB, dbName string) bool {
	var exists int
	query := fmt.Sprintf("SELECT COUNT(*) FROM sys.databases WHERE name = '%s'", dbName)
	db.Raw(query).Scan(&exists)
	return exists > 0
}

func createDatabase(db *gorm.DB, dbName string) {
	query := fmt.Sprintf("CREATE DATABASE %s", dbName)
	if err := db.Exec(query).Error; err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}
	log.Printf("Database %s created successfully!", dbName)
}
