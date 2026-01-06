package controllers

import (
	"encoding/base64"
	"fmt"
	"garment-management-backend/internal/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// Store: Create a new User
func (ctrl *UserController) Store(c *gin.Context) {
	var input models.UserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(ctrl.DB); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	branchID := input.BranchID

	// Handle Image Upload 📸
	var profileImagePath string
	if input.Photo != "" {
		// 1. Clean up base64 string
		b64data := input.Photo[strings.IndexByte(input.Photo, ',')+1:]

		// 2. Decode
		dec, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image data"})
			return
		}

		// 3. Create directory
		uploadDir := "uploads/users"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		// 4. Generate Filename
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), input.ImageName)
		if input.ImageName == "" {
			filename = fmt.Sprintf("%d_profile.%s", time.Now().Unix(), input.ImageFormat)
		}
		fullPath := filepath.Join(uploadDir, filename)

		// 5. Save File
		if err := os.WriteFile(fullPath, dec, 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		profileImagePath = fullPath
	}

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     string(hashedPassword),
		NIC:          input.NIC,
		MobileNo:     input.MobileNo,
		Address:      input.Address,
		BranchID:     branchID,
		ProfileImage: profileImagePath,
	}

	if err := ctrl.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	user.Password = "" // Hide password in response

	c.JSON(http.StatusCreated, user)
}

// Index: List all Users
func (ctrl *UserController) Index(c *gin.Context) {
	var users []models.User
	ctrl.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

// Show: Get single User
func (ctrl *UserController) Show(c *gin.Context) {
	var user models.User
	if err := ctrl.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Update: Update User
func (ctrl *UserController) Update(c *gin.Context) {
	var user models.User
	if err := ctrl.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.UserRequest
	c.ShouldBindJSON(&input)

	branchID := input.BranchID

	// Handle Image Upload 📸
	profileImagePath := user.ProfileImage

	if input.Photo != "" && strings.Contains(input.Photo, ",") {
		b64data := input.Photo[strings.IndexByte(input.Photo, ',')+1:]
		dec, err := base64.StdEncoding.DecodeString(b64data)
		if err == nil {
			uploadDir := "uploads/users"
			os.MkdirAll(uploadDir, os.ModePerm)

			filename := fmt.Sprintf("%d_%s", time.Now().Unix(), input.ImageName)
			if input.ImageName == "" {
				filename = fmt.Sprintf("%d_profile.%s", time.Now().Unix(), input.ImageFormat)
			}
			fullPath := filepath.Join(uploadDir, filename)

			if err := os.WriteFile(fullPath, dec, 0644); err == nil {
				profileImagePath = fullPath
			}
		}
	}

	updateData := map[string]interface{}{
		"name":          input.Name,
		"email":         input.Email,
		"nic":           input.NIC,
		"mobile_no":     input.MobileNo,
		"address":       input.Address,
		"branch_id":     branchID,
		"profile_image": profileImagePath,
	}

	// Only update password if the user actually typed a new one 🔑
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err == nil {
			updateData["password"] = string(hashedPassword)
		}
	}

	// Using a map ensures that empty strings like MobileNo: "" are actually updated 🎯
	if err := ctrl.DB.Model(&user).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Destroy: Delete User
func (ctrl *UserController) Destroy(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.User{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Deletion failed"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
