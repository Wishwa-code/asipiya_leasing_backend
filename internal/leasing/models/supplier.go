package models

import "gorm.io/gorm"

// Supplier maps to the `suppliers` table.
type Supplier struct {
	gorm.Model
	Name         string  `gorm:"column:name" json:"name"`
	NIC          string  `gorm:"column:nic" json:"nic"`
	Latitude     float64 `gorm:"column:latitude" json:"latitude"`
	Longitude    float64 `gorm:"column:longitude" json:"longitude"`
	Address      string  `gorm:"column:address" json:"address"`
	ContactNo    string  `gorm:"column:contact_no" json:"contact_no"`
	Occupation   string  `gorm:"column:occupation" json:"occupation"`
	Income       float64 `gorm:"column:income" json:"income"`
	NameInCheque string  `gorm:"column:name_in_cheque" json:"name_in_cheque"`
}

func (Supplier) TableName() string {
	return "suppliers"
}
