package models

import "gorm.io/gorm"

type InsuranceCompany struct {
	gorm.Model
	CompanyCode          string  `gorm:"column:company_code" json:"company_code"`
	CompanyName          string  `gorm:"column:company_name" json:"company_name"`
	HeadOfficeAddress    string  `gorm:"column:head_office_address" json:"head_office_address"`
	ContactPerson        string  `gorm:"column:contact_person" json:"contact_person"`
	ContactMobile        string  `gorm:"column:contact_mobile" json:"contact_mobile"`
	ContactEmail         string  `gorm:"column:contact_email" json:"contact_email"`
	ContactPerson2       string  `gorm:"column:contact_person2" json:"contact_person2"`
	ContactPerson2Mobile string  `gorm:"column:contact_person2_mobile" json:"contact_person2_mobile"`
	ContactPerson2Email  string  `gorm:"column:contact_person2_email" json:"contact_person2_email"`
	CommisionRate        float64 `gorm:"column:commision_rate" json:"commision_rate"`
	BankAccountNo        string  `gorm:"column:bank_account_no" json:"bank_account_no"`
	BankAccountName      string  `gorm:"column:bank_account_name" json:"bank_account_name"`
	BankName             string  `gorm:"column:bank_name" json:"bank_name"`
	Status               string  `gorm:"column:status;default:'Active'" json:"status"`
	CreatedBy            uint    `gorm:"column:created_by" json:"created_by"`
}

func (InsuranceCompany) TableName() string {
	return "insurance_companies"
}
