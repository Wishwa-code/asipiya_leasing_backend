package models

import "gorm.io/gorm"

// LeasingChequeDefine represents the main cheque definitions for the leasing application
type LeasingChequeDefine struct {
	gorm.Model
	LeasingApplicationID *uint  `gorm:"column:leasing_application_id" json:"leasing_application_id"`

	// Relationships
	LeasingApplication *LeasingApplication       `gorm:"foreignKey:LeasingApplicationID" json:"leasing_application,omitempty"`
	Items              []LeasingChequeDefineItem `gorm:"foreignKey:LeasingChequeDefineID" json:"items,omitempty"`
}

func (LeasingChequeDefine) TableName() string {
	return "leasing_cheque_defines"
}

// LeasingChequeDefineItem represents an individual cheque definition item
type LeasingChequeDefineItem struct {
	gorm.Model
	LeasingChequeDefineID *uint   `gorm:"column:leasing_cheque_define_id" json:"leasing_cheque_define_id"`
	PayeeName             string  `gorm:"column:payee_name;type:varchar(255)" json:"payee_name"`
	NicBrNo               string  `gorm:"column:nic_br_no;type:varchar(100)" json:"nic_br_no"`
	Instructions          string  `gorm:"column:instructions;type:text" json:"instructions"`
	PaymentAmount         float64 `gorm:"column:payment_amount" json:"payment_amount"`
	BankName              string  `gorm:"column:bank_name;type:varchar(100)" json:"bank_name"`
	BranchName            string  `gorm:"column:branch_name;type:varchar(100)" json:"branch_name"`
	AccountNumber         string  `gorm:"column:account_number;type:varchar(100)" json:"account_number"`

	ChequeDefine *LeasingChequeDefine `gorm:"foreignKey:LeasingChequeDefineID" json:"cheque_define,omitempty"`
}

func (LeasingChequeDefineItem) TableName() string {
	return "leasing_cheque_define_items"
}
