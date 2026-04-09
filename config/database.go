package config

import (
	"fmt"
	"log"

	"project/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.DBTimeZone,
	)

	fmt.Println("DB CONNECTION STRING:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}

	if err := db.AutoMigrate(&models.HierarchyNode{}); err != nil {
		log.Fatal("database migration failed: ", err)
	}

	log.Println("Database connected successfully")
	return db
}