package models

import "gorm.io/gorm"

type InsuranceCompany struct {
	gorm.Model
	CompanyName string `gorm:"column:company_name" json:"company_name"`
	ContactNo   string `gorm:"column:contact_no" json:"contact_no"`
	Email       string `gorm:"column:email" json:"email"`
	Status      string `gorm:"column:status;default:'Active'" json:"status"`
}

func (InsuranceCompany) TableName() string {
	return "insurance_companies"
}
