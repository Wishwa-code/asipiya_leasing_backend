package models

import (
	"gorm.io/gorm"
)

// LeasingApplication maps to the former customer_leasing table, tracking the stepper progress.
type LeasingApplication struct {
	gorm.Model
	LeasingApplicationLoanNo string  `gorm:"column:leasing_application_loan_no;type:varchar(255)" json:"leasing_application_loan_no"`
	CustomerID               uint    `gorm:"column:customer_id;not null" json:"customer_id"`
	IntroducerID             *uint   `gorm:"column:introducer_id" json:"introducer_id"`
	LeasingLoanCode          string  `gorm:"column:leasing_loan_code;type:varchar(255)" json:"leasing_loan_code"`
	Status                   string  `gorm:"column:status;type:varchar(50);default:'draft'" json:"status"`
	BranchID                 *uint   `gorm:"column:branch_id" json:"branch_id"`
	CurrentProgressData      string  `gorm:"column:current_progress_data;type:jsonb" json:"current_progress_data"` // Stores draft JSON

	// Relationships
	Customer           *Customer          `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Introducer         *Introducer        `gorm:"foreignKey:IntroducerID" json:"introducer,omitempty"`
	Vehicle            *LeasingVehicle    `gorm:"foreignKey:LeasingApplicationID" json:"vehicle,omitempty"`
	Loan               *LeasingLoan       `gorm:"foreignKey:LeasingApplicationID" json:"loan,omitempty"`
	Guarantors         []LeasingGuarantor `gorm:"foreignKey:LeasingApplicationID" json:"guarantors,omitempty"`
	PdcSecurity        *PdcSecurity       `gorm:"foreignKey:LeasingApplicationID" json:"pdc_security,omitempty"`
	ChequeDefine       *LeasingChequeDefine `gorm:"foreignKey:LeasingApplicationID" json:"cheque_define,omitempty"`
	
	// Temporary drafts might have uploaded documents linked directly here before the final sub-records are created
	DocumentImages     []LeasingVehicleDocumentImage `gorm:"foreignKey:LeasingApplicationID" json:"document_images,omitempty"`
}

func (LeasingApplication) TableName() string {
	return "leasing_applications"
}
