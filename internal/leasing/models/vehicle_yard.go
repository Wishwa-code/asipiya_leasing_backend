package models

import "gorm.io/gorm"

type VehicleYard struct {
	gorm.Model
	YardName      string `gorm:"column:yard_name" json:"yard_name"`
	Address       string `gorm:"column:address" json:"address"`
	Province      string `gorm:"column:province" json:"province"`
	District      string `gorm:"column:district" json:"district"`
	DSD           string `gorm:"column:dsd" json:"dsd"`
	YardContactNo string `gorm:"column:yard_contact_no" json:"yard_contact_no"`
	ContactPerson string `gorm:"column:contact_person" json:"contact_person"`
	MobileNo      string `gorm:"column:mobile_no" json:"mobile_no"`
	Status        string `gorm:"column:status;default:'Active'" json:"status"`
}

func (VehicleYard) TableName() string {
	return "vehicle_yards"
}
