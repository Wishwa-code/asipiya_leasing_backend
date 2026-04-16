package models

import "gorm.io/gorm"

// CustomerLoan maps to the `customer_loans` table.
// Referenced in the Customer model's loans() relationship and getCustomerLoans() in the controller.
type CustomerLoan struct {
	gorm.Model
	CustomerID       uint    `gorm:"column:customer_id" json:"customer_id"`
	LoanNo           string  `gorm:"column:loan_no" json:"loan_no"`
	LoanStatus       string  `gorm:"column:loan_status" json:"loan_status"` // active | disbursed | approved | closed
	TotalLoanAmount  float64 `gorm:"column:total_loan_amount" json:"total_loan_amount"`
	TotalLoanBalance float64 `gorm:"column:total_loan_balance" json:"total_loan_balance"`
	BranchID         *uint   `gorm:"column:branch_id" json:"branch_id"`
	RouteID          *uint   `gorm:"column:route_id" json:"route_id"`
	ProductID        *uint   `gorm:"column:product_id" json:"product_id"`
	CreatedBy        *uint   `gorm:"column:created_by" json:"created_by"`
	UpdatedBy        *uint   `gorm:"column:updated_by" json:"updated_by"`
}

func (CustomerLoan) TableName() string {
	return "customer_loans"
}
