package models

import "gorm.io/gorm"

type Introducer struct {
	gorm.Model
	IntroducerType   string  `gorm:"column:introducer_type" json:"introducer_type"`
	Name             string  `gorm:"column:name" json:"name"`
	RegistrationNo   string  `gorm:"column:registration_no" json:"registration_no"`
	ContactPerson    string  `gorm:"column:contact_person" json:"contact_person"`
	PrimaryContact   string  `gorm:"column:primary_contact" json:"primary_contact"`
	SecondaryContact string  `gorm:"column:secondary_contact" json:"secondary_contact"`
	Email            string  `gorm:"column:email" json:"email"`
	Address          string  `gorm:"column:address" json:"address"`
	CommissionRate   float64 `gorm:"column:commission_rate" json:"commission_rate"`
	BankDetails      string  `gorm:"column:bank_details;type:jsonb" json:"bank_details"`
	Remarks          string  `gorm:"column:remarks" json:"remarks"`
	Status           string  `gorm:"column:status;default:'Active'" json:"status"`
	CreatedBy        uint    `gorm:"column:created_by" json:"created_by"`
}

func (Introducer) TableName() string {
	return "introducers"
}
