

package garmentOperation

import (
	// "net/http"
	"github.com/gin-gonic/gin"
    // "gorm.io/gorm"
	"garment-management-backend/internal/database"
    // "garment-management-backend/internal/models"
	"garment-management-backend/internal/garmentOperation/controllers"
)


// RegisterRoutes sets up the plantation specific routes 📍
func RegisterRoutes(rg *gin.RouterGroup) {

    userCtrl := &controllers.UserController{DB: database.DB}
	dailyCtrl := &controllers.DailyAmountController{DB: database.DB}
	styleCtrl := &controllers.StyleController{DB: database.DB}
	employeeCtrl := &controllers.EmployeeController{DB: database.DB}
	opDataCtrl := &controllers.OperationDataController{DB: database.DB}

	v1 := rg.Group("/v1")
	{
		v1.GET("/stats", controllers.GetStatsHandler)

        // User Resource Routes (Laravel-style) 👤
        users := v1.Group("/users")
        {
            users.GET("", userCtrl.Index)          // GET /v1/users
            users.POST("", userCtrl.Store)        // POST /v1/users
            users.GET("/:id", userCtrl.Show)       // GET /v1/users/:id
            users.PUT("/:id", userCtrl.Update)     // PUT /v1/users/:id
            users.DELETE("/:id", userCtrl.Destroy)  // DELETE /v1/users/:id
        }

		// Highlight: New Daily Amount Routes
        daily := v1.Group("/daily-amounts")
        {
            daily.GET("", dailyCtrl.Index)
            daily.POST("", dailyCtrl.Store)
            daily.GET("verify-target/:id", dailyCtrl.VerifyStyleTarget)
        }

		styles := v1.Group("/styles")
        {
            styles.GET("", styleCtrl.Index)
            styles.POST("", styleCtrl.Store)
            styles.GET("/:id", styleCtrl.Show)
            styles.PUT("/:id", styleCtrl.Update)
            styles.DELETE("/:id", styleCtrl.Destroy)
        }

		employees := v1.Group("/employees")
        {
            employees.GET("", employeeCtrl.Index)
            employees.POST("", employeeCtrl.Store)
            employees.GET("/:id", employeeCtrl.Show)
            employees.PUT("/:id", employeeCtrl.Update)
            employees.DELETE("/:id", employeeCtrl.Destroy)
        }

		// Operation Data Routes 📍
        operationData := v1.Group("/operation-data")
        {
            operationData.GET("", opDataCtrl.GetDailyAmounts)
            // operationData.POST("", opDataCtrl.Store)
        }

	}
}

