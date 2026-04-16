package models

import "gorm.io/gorm"

// CustomerSavingAccount maps to the `customer_saving_accounts` table.
// Referenced in the Customer model's savingAccounts() relationship and
// createDefaultSavingAccount() in the Laravel controller.
type CustomerSavingAccount struct {
	gorm.Model
	CustomerID    uint    `gorm:"column:customer_id" json:"customer_id"`
	AccountNo     string  `gorm:"column:account_no" json:"account_no"`
	OpenDate      string  `gorm:"column:open_date" json:"open_date"`
	Balance       float64 `gorm:"column:balance" json:"balance"`
	InterestRate  float64 `gorm:"column:interest_rate" json:"interest_rate"`
	InterestType  string  `gorm:"column:interest_type" json:"interest_type"`
	AccountStatus string  `gorm:"column:account_status" json:"account_status"` // active | inactive
	BranchID      *uint   `gorm:"column:branch_id" json:"branch_id"`
	CreatedBy     *uint   `gorm:"column:created_by" json:"created_by"`
}

func (CustomerSavingAccount) TableName() string {
	return "customer_saving_accounts"
}
