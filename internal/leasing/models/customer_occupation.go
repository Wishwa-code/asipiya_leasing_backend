package models

import "gorm.io/gorm"

// CustomerOccupation maps to the `customer_occupations` table.
// Column names match the Laravel controller's occupation()->create() calls.
type CustomerOccupation struct {
	gorm.Model
	CustomerID          uint    `gorm:"column:customer_id" json:"customer_id"`
	Type                string  `gorm:"column:type" json:"engagementType"`
	Designation         string  `gorm:"column:designation" json:"position"`
	BRNo                string  `gorm:"column:br_no" json:"registrationNumber"`
	BusinessName        string  `gorm:"column:business_name" json:"businessName"`
	NatureOfBusiness    string  `gorm:"column:nature_of_business" json:"natureOfBusiness"`
	BusinessAddress1    string  `gorm:"column:business_address_01" json:"businessAddress1"`
	BusinessAddress2    string  `gorm:"column:business_address_02" json:"businessAddress2"`
	BusinessAddress3    string  `gorm:"column:business_address_03" json:"businessAddress3"`
	ContactNo           string  `gorm:"column:contact_no" json:"contactNo"`
	Email               string  `gorm:"column:email" json:"email"`
	EmployerName        string  `gorm:"column:employer_name" json:"employerName"`
	MonthlySalaryIncome float64 `gorm:"column:monthly_salary_income" json:"netMonthlyIncome"`
	FromDate            string  `gorm:"column:from_date" json:"startDate"`
	ToDate              string  `gorm:"column:to_date" json:"endDate"`
	Longitude           float64 `gorm:"column:longitude" json:"longitude"`
	Latitude            float64 `gorm:"column:latitude" json:"latitude"`
}

func (CustomerOccupation) TableName() string {
	return "customer_occupations"
}
