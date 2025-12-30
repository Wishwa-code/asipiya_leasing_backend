package models

import (
    // "time"
    "gorm.io/gorm"
    "errors"
)

// User represents the system users for authentication 👤
type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	NIC      string `gorm:"uniqueIndex"`
	MobileNo string
	Address  string
	RoleID   *uint  // Pointer to allow nulls if no role is assigned
	BranchID *uint  // Pointer to allow nulls
}

// UserRequest handles incoming payload validation 📥
type UserRequest struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	NIC                  string `json:"nic"`
	MobileNo             string `json:"mobile_no"`
	Address              string `json:"address"`
	RoleID               string `json:"role_id"`   // Received as string from payload
	BranchID             *uint  `json:"branch_id"` // Can be null
}


func (req *UserRequest) Validate(db *gorm.DB) error {
    var count int64
    // Check Email and NIC in one query for performance 🏎️
    if err := db.Model(&User{}).
        Where("email = ? OR nic = ?", req.Email, req.NIC).
        Count(&count).Error; err != nil {
        return err
    }
    if count > 0 {
        return errors.New("email or NIC already registered")
    }
    return nil
}
// // Product represents your domain data 📦
// type Producttest struct {
//     ID              string         `gorm:"primaryKey"`
//     Title           string         `gorm:"size:255"`
//     Description     string         `gorm:"type:text"`
//     TagOne          string         
//     TagTwo          string
//     ImageURL        string
//     Department      string         `gorm:"default:'mainBuilding'"`
//     MainCategory    string
//     SubCategory     string
//     CreatedAt       time.Time
//     LastModifiedAt  time.Time
// }

// LoginRequest remains as a simple DTO (Not a DB table)
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}