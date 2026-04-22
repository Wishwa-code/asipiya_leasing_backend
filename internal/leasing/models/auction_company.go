package models

import "gorm.io/gorm"

type AuctionCompany struct {
	gorm.Model
	Name          string `gorm:"column:name" json:"name"`
	ContactNo1    string `gorm:"column:contact_no_1" json:"contact_no_1"`
	ContactNo2    string `gorm:"column:contact_no_2" json:"contact_no_2"`
	ContactPerson string `gorm:"column:contact_person" json:"contact_person"`
	Address       string `gorm:"column:address" json:"address"`
	Note          string `gorm:"column:note" json:"note"`
}

func (AuctionCompany) TableName() string {
	return "auction_companies"
}
