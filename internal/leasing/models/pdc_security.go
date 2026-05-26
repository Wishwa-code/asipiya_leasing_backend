package models

import "gorm.io/gorm"

// PdcSecurity represents the main security grouping for PDC
type PdcSecurity struct {
	gorm.Model
	LeasingApplicationID *uint  `gorm:"column:leasing_application_id" json:"leasing_application_id"`
	PdcSecurityType      string `gorm:"column:pdc_security_type;type:varchar(50)" json:"pdc_security_type"`

	// Relationships
	LeasingApplication *LeasingApplication `gorm:"foreignKey:LeasingApplicationID" json:"leasing_application,omitempty"`
	ChequeDetails      []PdcChequeDetail   `gorm:"foreignKey:PdcSecurityID" json:"cheque_details,omitempty"`
	CrBookDetails      []PdcCrBookDetail   `gorm:"foreignKey:PdcSecurityID" json:"cr_book_details,omitempty"`
	DeedDetails        []PdcDeedDetail     `gorm:"foreignKey:PdcSecurityID" json:"deed_details,omitempty"`
}

func (PdcSecurity) TableName() string {
	return "pdc_securities"
}

// PdcChequeDetail
type PdcChequeDetail struct {
	gorm.Model
	PdcSecurityID    *uint  `gorm:"column:pdc_security_id" json:"pdc_security_id"`
	ChequeStatus     string `gorm:"column:cheque_status;type:varchar(50)" json:"cheque_status"`
	BankID           *uint  `gorm:"column:bank_id" json:"bank_id"`
	ChequeDate       string `gorm:"column:cheque_date;type:varchar(50)" json:"cheque_date"`
	ChequeNo         string `gorm:"column:cheque_no;type:varchar(100)" json:"cheque_no"`
	Ownership        string `gorm:"column:ownership;type:varchar(100)" json:"ownership"`
	ReferenceDetails string `gorm:"column:reference_details;type:text" json:"reference_details"`

	PdcSecurity *PdcSecurity `gorm:"foreignKey:PdcSecurityID" json:"pdc_security,omitempty"`
	Bank        *Bank        `gorm:"foreignKey:BankID" json:"bank,omitempty"`
}

func (PdcChequeDetail) TableName() string {
	return "pdc_cheque_details"
}

// PdcCrBookDetail
type PdcCrBookDetail struct {
	gorm.Model
	PdcSecurityID    *uint  `gorm:"column:pdc_security_id" json:"pdc_security_id"`
	BookDate         string `gorm:"column:book_date;type:varchar(50)" json:"book_date"`
	ReferenceDetails string `gorm:"column:reference_details;type:text" json:"reference_details"`

	PdcSecurity *PdcSecurity `gorm:"foreignKey:PdcSecurityID" json:"pdc_security,omitempty"`
}

func (PdcCrBookDetail) TableName() string {
	return "pdc_cr_book_details"
}

// PdcDeedDetail
type PdcDeedDetail struct {
	gorm.Model
	PdcSecurityID    *uint  `gorm:"column:pdc_security_id" json:"pdc_security_id"`
	ReferenceDetails string `gorm:"column:reference_details;type:text" json:"reference_details"`

	PdcSecurity *PdcSecurity `gorm:"foreignKey:PdcSecurityID" json:"pdc_security,omitempty"`
}

func (PdcDeedDetail) TableName() string {
	return "pdc_deed_details"
}
