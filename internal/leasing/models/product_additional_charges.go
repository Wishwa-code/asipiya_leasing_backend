package models

import "gorm.io/gorm"

// ProductAdditionalCharges maps to the `product_additional_charges` table
type ProductAdditionalCharges struct {
	gorm.Model // Adds ID, CreatedAt, UpdatedAt, DeletedAt
	ProductID     uint    `gorm:"column:product_id" json:"product_id"`
	Description   string  `gorm:"column:description" json:"description"`
	ValueType     string  `gorm:"column:value_type" json:"value_type"`
	Value         float64 `gorm:"column:value" json:"value"`
	DeductionType string  `gorm:"column:deduction_type" json:"deduction_type"`
}

// TableName explicitly overrides pluralization to match exactly
func (ProductAdditionalCharges) TableName() string {
	return "product_additional_charges"
}
