package models

import "gorm.io/gorm"

// ProductRequiredDocuments maps to the `product_required_documents` table
type ProductRequiredDocuments struct {
	gorm.Model // Adds ID, CreatedAt, UpdatedAt, DeletedAt
	ProductID      uint   `gorm:"column:product_id" json:"product_id"`
	Name           string `gorm:"column:name" json:"name"`
	RequiredStatus string `gorm:"column:required_status" json:"required_status"`
}

// TableName explicitly overrides pluralization to match exactly
func (ProductRequiredDocuments) TableName() string {
	return "product_required_documents"
}
