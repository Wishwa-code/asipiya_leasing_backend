package main

import (
	"fmt"
	"garment-management-backend/internal/database"
	leasingModels "garment-management-backend/internal/leasing/models"
	"log"
)

func main() {
	database.Connect()
	if database.DB == nil {
		log.Fatal("DB is nil")
	}

	var customerCount int64
	database.DB.Model(&leasingModels.Customer{}).Count(&customerCount)
	fmt.Printf("Customer count: %d\n", customerCount)

	var customers []leasingModels.Customer
	database.DB.Limit(5).Find(&customers)
	for _, c := range customers {
		fmt.Printf("- Customer: ID=%d, Name=%s, Code=%s, NIC=%s\n", c.ID, c.FullName, c.CustomerCode, c.NewNic)
	}

	var productCount int64
	database.DB.Model(&leasingModels.Product{}).Count(&productCount)
	fmt.Printf("Product count: %d\n", productCount)
}
