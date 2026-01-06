package models

import (
	// "time"
	"errors"

	"gorm.io/gorm"
)

// User represents the system users for authentication 👤
type User struct {
	gorm.Model
	Name         string `gorm:"not null" json:"name"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	Password     string `gorm:"not null" json:"-"`
	NIC          string `gorm:"uniqueIndex" json:"nic"`
	MobileNo     string `json:"mobile_no"`
	Address      string `json:"address"`
	ProfileImage string `json:"profile_image"` // Stores relative path to image
	RoleID       *uint  `json:"role_id"`       // Pointer to allow nulls if no role is assigned
	BranchID     *uint  `json:"branch_id"`     // Pointer to allow nulls
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
	Photo                string `json:"photo"`        // Base64 string
	ImageFormat          string `json:"image_format"` // e.g., "png"
	ImageName            string `json:"image_name"`   // e.g., "my_photo.png"
	RoleID               string `json:"role_id"`      // Received as string from payload
	BranchID             *uint  `json:"branch_id"`    // Can be null
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
