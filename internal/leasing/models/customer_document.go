package models

import "gorm.io/gorm"

// CustomerDocument maps to the `customer_documents` table.
// Column names (PascalCase) match the Laravel saveFiles() method:
//   Description, Path, Customer_idCustomer
type CustomerDocument struct {
	gorm.Model
	CustomerID  uint   `gorm:"column:Customer_idCustomer" json:"customer_id"`
	Description string `gorm:"column:Description" json:"description"` // Category / document name
	Path        string `gorm:"column:Path" json:"path"`               // Relative path under /uploads
}

func (CustomerDocument) TableName() string {
	return "customer_documents"
}
