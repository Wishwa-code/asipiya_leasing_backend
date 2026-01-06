package auth

import (
	"fmt"
	"garment-management-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CheckCredentials verifies the username and password 🔍
// Moving this here allows main.go to stay focused on routing.

// func CheckCredentials(username, password string) bool {
// 	fmt.Printf("CheckCredentials called for user: %s\n", username)

// 	TODO: Replace with database lookup or bcrypt comparison
// 	return username == "admin" && password == "password"
// }

// func CheckCredentials(db *gorm.DB, email, password string) (*models.User, bool) {
//     var user models.User

//     if err := db.Where("email = ?", email).First(&user).Error; err != nil {
//         return nil, false
//     }

//     err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
//     if err != nil {
//         return nil, false
//     }

//     return &user, true
// }

func CheckCredentials(db *gorm.DB, email, password string) (*models.User, bool) {
	var user models.User

	// 📝 Log 1: Input arguments
	fmt.Printf("🔍 Auth Attempt - Email: [%s] | Password Length: %d\n", email, len(password))

	// 🔍 Look for the user by email
	result := db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		// 📝 Log 2: Database Error details
		fmt.Printf("❌ DB Error: %v (User not found for email: %s)\n", result.Error, email)
		return nil, false
	}

	// 📝 Log 3: DB Record found
	fmt.Printf("✅ DB Record Found: ID [%d] | Saved Password: [%s]\n", user.ID, user.Password)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, false
	}

	user.Password = "" // Extra layer of security
	return &user, true
}
