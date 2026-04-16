package models

import "gorm.io/gorm"

// Customer maps to the Laravel `customers` table.
// The fillable columns from Laravel are the source of truth for what's stored locally.
// Extended fields (DOB, gender, address lines etc.) that Laravel fetches from Account Center
// are stored directly here since this is a standalone Go backend.
type Customer struct {
	gorm.Model
	// Core identity — from $fillable
	CustomerCode     string  `gorm:"column:customer_code" json:"customer_code"`
	Title            string  `gorm:"column:title" json:"title"`
	FullName         string  `gorm:"column:full_name" json:"full_name"`
	NameWithInitials string  `gorm:"column:name_with_initials" json:"name_with_initials"`
	FirstName        string  `gorm:"column:first_name" json:"first_name"`
	LastName         string  `gorm:"column:last_name" json:"last_name"`
	Email            string  `gorm:"column:email" json:"email"`
	ContactNo        string  `gorm:"column:contact_no" json:"contact_no"`          // primary mobile
	ContactNo2       string  `gorm:"column:contact_no_02" json:"contact_no_02"`    // secondary mobile
	Landline         string  `gorm:"column:landline" json:"landline"`
	NewNic           string  `gorm:"column:new_nic" json:"new_nic"`
	OldNic           string  `gorm:"column:old_nic" json:"old_nic"`
	Status           string  `gorm:"column:status" json:"status"`
	Latitude         float64 `gorm:"column:latitude" json:"latitude"`
	Longitude        float64 `gorm:"column:longitude" json:"longitude"`

	// FK references — from $fillable
	RouteID    *uint `gorm:"column:route_id" json:"route_id"`
	CenterID   *uint `gorm:"column:center_id" json:"center_id"`
	GroupID    *uint `gorm:"column:group_id" json:"group_id"`
	BranchID   *uint `gorm:"column:branch_id" json:"branch_id"`
	AccCenterCusID *uint `gorm:"column:acc_center_cus_id" json:"acc_center_cus_id"`

	// Audit columns — from $fillable
	CreatedBy *uint `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *uint `gorm:"column:updated_by" json:"updated_by"`
	DeletedBy *uint `gorm:"column:deleted_by" json:"deleted_by"`

	// Extended fields — NOT in Laravel $fillable (they come from Account Center there),
	// but stored locally here since this Go backend is standalone.
	DOB                string `gorm:"column:dob" json:"dob"`
	Gender             string `gorm:"column:gender" json:"gender"`
	Remarks            string `gorm:"column:remarks" json:"remarks"`
	Province           string `gorm:"column:province" json:"province"`
	City               string `gorm:"column:city" json:"city"`
	PostalProvince     string `gorm:"column:postal_province" json:"postal_province"`
	PostalCity         string `gorm:"column:postal_city" json:"postal_city"`
	PerAddressLine1    string `gorm:"column:per_address_line_1" json:"permanent_address_line1"`
	PerAddressLine2    string `gorm:"column:per_address_line_2" json:"permanent_address_line2"`
	PerAddressLine3    string `gorm:"column:per_address_line_3" json:"permanent_address_line3"`
	PostalAddressLine1 string `gorm:"column:postal_address_line_1" json:"postal_address_line1"`
	PostalAddressLine2 string `gorm:"column:postal_address_line_2" json:"postal_address_line2"`
	PostalAddressLine3 string `gorm:"column:postal_address_line_3" json:"postal_address_line3"`

	// Relationships
	Occupations    []CustomerOccupation    `gorm:"foreignKey:CustomerID" json:"occupations"`
	BankAccounts   []CustomerBankAccount   `gorm:"foreignKey:CustomerID" json:"bank_accounts"`
	Documents      []CustomerDocument      `gorm:"foreignKey:CustomerID" json:"documents"`
	SavingAccounts []CustomerSavingAccount `gorm:"foreignKey:CustomerID" json:"saving_accounts"`
}

func (Customer) TableName() string {
	return "customers"
}
