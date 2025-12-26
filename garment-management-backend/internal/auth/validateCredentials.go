package auth

import "fmt"

// CheckCredentials verifies the username and password 🔍
// Moving this here allows main.go to stay focused on routing.
func CheckCredentials(username, password string) bool {
	fmt.Printf("CheckCredentials called for user: %s\n", username)
	
	// TODO: Replace with database lookup or bcrypt comparison
	return username == "admin" && password == "password"
}