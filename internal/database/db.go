package database

import (
	// "database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"garment-management-backend/internal/garmentOperation/operationModels"
	leasingModels "garment-management-backend/internal/leasing/models"
	"garment-management-backend/internal/models"
	_ "github.com/lib/pq"
)

// DB is now a GORM instance 🔗
var DB *gorm.DB

func Connect() {
	// Your Neon PostgreSQL connection string
	dsn := "postgresql://neondb_owner:npg_UyWhNfBkqC21@ep-spring-rain-a1gdk8hj-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

	var err error
	// Open connection using GORM 🚀	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
		// Add GORM config here if needed (e.g., Logger)
	})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&operationModels.DailyReport{},
		&operationModels.DailyAmount{},
		&operationModels.Style{},
		&operationModels.Employee{},
		&leasingModels.Product{},
		&leasingModels.ProductHasItem{},
		&leasingModels.ProductAdditionalCharges{},
		&leasingModels.ProductRequiredDocuments{},
	)

	if err != nil {
		log.Fatalf("❌ Migration Failed: %v", err)
	}

	// Optional: Configure Connection Pool 🏊‍♂️
	sqlDB, err := DB.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
	}

	fmt.Println("✅ GORM connected to PostgreSQL successfully")
}
