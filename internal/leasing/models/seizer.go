package models

import "gorm.io/gorm"

type Seizer struct {
	gorm.Model
	SeizerType          string `gorm:"column:seizer_type" json:"seizer_type"`
	CompanyName         string `gorm:"column:company_name" json:"company_name"`
	CompanyRegistration string `gorm:"column:company_registration" json:"company_registration"`
	CompanyContactNo    string `gorm:"column:company_contact_no" json:"company_contact_no"`
	NIC                 string `gorm:"column:nic" json:"nic"`
	SeizerContactNo     string `gorm:"column:seizer_contact_no" json:"seizer_contact_no"`
	MobileNo            string `gorm:"column:mobile_no" json:"mobile_no"`
	Address             string `gorm:"column:address" json:"address"`
	Remarks             string `gorm:"column:remarks" json:"remarks"`
	Status              string `gorm:"column:status;default:'Active'" json:"status"`
}

func (Seizer) TableName() string {
	return "seizers"
}
