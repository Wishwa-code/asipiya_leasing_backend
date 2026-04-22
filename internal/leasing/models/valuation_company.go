package models

import "gorm.io/gorm"

type ValuationCompany struct {
	gorm.Model
	CompanyName          string `gorm:"column:company_name" json:"company_name"`
	ContactPerson1Name   string `gorm:"column:contact_person_1_name" json:"contact_person_1_name"`
	ContactPerson2Name   string `gorm:"column:contact_person_2_name" json:"contact_person_2_name"`
	ContactPerson1Mobile string `gorm:"column:contact_person_1_mobile" json:"contact_person_1_mobile"`
	ContactPerson2Mobile string `gorm:"column:contact_person_2_mobile" json:"contact_person_2_mobile"`
	Address              string `gorm:"column:address" json:"address"`
	Note                 string `gorm:"column:note" json:"note"`
}

func (ValuationCompany) TableName() string {
	return "valuation_companies"
}
