package models

import "gorm.io/gorm"

// CustomerBankAccount maps to the `customer_has_bank` table.
// FK is `cus_id` (confirmed by saveBank method in the Laravel controller).
type CustomerBankAccount struct {
	gorm.Model
	CustomerID    uint   `gorm:"column:cus_id" json:"customer_id"`
	BankName      string `gorm:"column:bank_name" json:"bank"`
	AccountName   string `gorm:"column:account_name" json:"beneficiary"`
	AccountNumber string `gorm:"column:account_number" json:"accountNumber"`
	Branch        string `gorm:"column:branch" json:"branch"`
	BankCode      string `gorm:"column:bank_code" json:"bankCode"`
	BranchCode    string `gorm:"column:branch_code" json:"branchCode"`
	AccountType   string `gorm:"column:account_type" json:"type"` // Savings / Current
}

func (CustomerBankAccount) TableName() string {
	return "customer_has_bank"
}
