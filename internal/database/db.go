package database

import (
	// "database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"

	adminModels "garment-management-backend/internal/admin/models"
	leasingModels "garment-management-backend/internal/leasing/models"
	"garment-management-backend/internal/models"
	_ "github.com/lib/pq"
)

// DB is now a GORM instance 🔗
var DB *gorm.DB

func Connect() {
	// Your Neon PostgreSQL connection string
	dsn := "postgresql://neondb_owner:npg_UyWhNfBkqC21@ep-spring-rain-a1gdk8hj-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require&connect_timeout=30"

	var err error
	// Open connection using GORM 🚀
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
		// Add GORM config here if needed (e.g., Logger)
	})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Configure connection pool FIRST before any operations
	// This is critical for Neon serverless - keeps idle connections alive
	sqlDB, poolErr := DB.DB()
	if poolErr == nil {
		sqlDB.SetMaxIdleConns(2)           // Low idle conns for serverless (Neon suspends when idle)
		sqlDB.SetMaxOpenConns(10)           // Keep open connections low for Neon free tier
		sqlDB.SetConnMaxLifetime(5 * time.Minute) // Recycle connections before Neon suspends them
		sqlDB.SetConnMaxIdleTime(2 * time.Minute) // Release idle connections quickly
	}

	// Ping with retry to wake up Neon from cold start (Neon serverless suspends when idle)
	fmt.Println("⏳ Waiting for database connection (Neon may be waking up)...")
	for attempt := 1; attempt <= 5; attempt++ {
		if sqlDB != nil {
			if pingErr := sqlDB.Ping(); pingErr == nil {
				fmt.Printf("✅ Database ping successful on attempt %d\n", attempt)
				break
			} else if attempt == 5 {
				log.Fatalf("❌ Could not reach database after %d attempts: %v", attempt, pingErr)
			} else {
				fmt.Printf("⚠️  Ping attempt %d failed, retrying in 3s... (%v)\n", attempt, pingErr)
				time.Sleep(3 * time.Second)
			}
		}
	}

	err = DB.AutoMigrate(
		&models.User{},
		&leasingModels.Product{},
		&leasingModels.ProductHasItem{},
		&leasingModels.ProductAdditionalCharges{},
		&leasingModels.ProductRequiredDocuments{},
		&leasingModels.Customer{},
		&leasingModels.CustomerOccupation{},
		&leasingModels.CustomerBankAccount{},
		&leasingModels.CustomerDocument{},
		&leasingModels.CustomerSavingAccount{},
		&leasingModels.CustomerLoan{},
		&leasingModels.Bank{},
		&leasingModels.InsuranceCompany{},
		&adminModels.VehicleType{},
		&adminModels.VehicleMake{},
		&adminModels.VehicleModel{},
		&adminModels.Color{},
		&leasingModels.Supplier{},
		&leasingModels.Seizer{},
		&leasingModels.Introducer{},
		&leasingModels.ValuationCompany{},
		&leasingModels.AuctionCompany{},
		&leasingModels.VehicleYard{},

		// Leasing Application (Stepper) Models
		&leasingModels.LeasingApplication{},
		&leasingModels.LeasingVehicle{},
		&leasingModels.LeasingVehicleDocumentImage{},
		&leasingModels.LeasingVehicleAudit{},
		&leasingModels.LeasingLoan{},
		&leasingModels.LeasingGuarantor{},
		&leasingModels.PdcSecurity{},
		&leasingModels.PdcChequeDetail{},
		&leasingModels.PdcCrBookDetail{},
		&leasingModels.PdcDeedDetail{},
		&leasingModels.LeasingChequeDefine{},
		&leasingModels.LeasingChequeDefineItem{},
	)

	if err != nil {
		log.Fatalf("❌ Migration Failed: %v", err)
	}

	// Seed database with default values 🚀
	if err := SeedDatabase(DB); err != nil {
		log.Printf("⚠️ Database Seeding Failed: %v", err)
	}

	fmt.Println("✅ GORM connected to PostgreSQL successfully")
}
