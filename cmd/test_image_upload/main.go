package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type UserRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	NIC                  string `json:"nic"`
	MobileNo             string `json:"mobile_no"`
	Address              string `json:"address"`
	Photo                string `json:"photo"`
	ImageFormat          string `json:"image_format"`
	ImageName            string `json:"image_name"`
}

func main() {
	baseURL := "http://localhost:8080"

	// 1. Login
	loginReq := LoginRequest{
		Username: "wishwa@gmail.com",
		Password: "12345678",
	}
	loginBody, _ := json.Marshal(loginReq)
	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Login failed with status %d: %s\n", resp.StatusCode, string(body))
		return
	}

	var tokens TokenResponse
	json.NewDecoder(resp.Body).Decode(&tokens)
	fmt.Println("Login successful, token received.")

	// 2. Create User with Image
	// Small red dot png base64
	b64Image := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg=="

	uniqueID := time.Now().Unix()
	newUser := UserRequest{
		Name:                 "Test User",
		Email:                fmt.Sprintf("testuser%d@example.com", uniqueID),
		Password:             "password123",
		PasswordConfirmation: "password123",
		NIC:                  fmt.Sprintf("%d", uniqueID),
		MobileNo:             "0771234567",
		Address:              "Test Address",
		Photo:                b64Image,
		ImageFormat:          "png",
		ImageName:            "test_pixel.png",
	}

	userBody, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/users", bytes.NewBuffer(userBody))
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Create user failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Create User Response: %s\n", string(body))

	if resp.StatusCode == 201 {
		fmt.Println("User created successfully!")
	} else {
		fmt.Printf("Failed to create user, status: %d\n", resp.StatusCode)
	}
}
