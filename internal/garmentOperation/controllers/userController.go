package controllers

import (
    "net/http"
	"gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "garment-management-backend/internal/models"
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

    user := models.User{
        Name:     input.Name,
        Email:    input.Email,
        Password: string(hashedPassword),
        NIC:      input.NIC,
        MobileNo: input.MobileNo,
        Address:  input.Address,
        BranchID: input.BranchID,
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

    ctrl.DB.Model(&user).Updates(models.User{
        Name:     input.Name,
        MobileNo: input.MobileNo,
        Address:  input.Address,
    })

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