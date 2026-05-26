package leasing

import (
	"garment-management-backend/internal/database"
	"garment-management-backend/internal/leasing/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the leasing specific routes
func RegisterRoutes(rg *gin.RouterGroup) {
	productCtrl := &controllers.ProductController{DB: database.DB}
	customerCtrl := &controllers.CustomerController{DB: database.DB}
	lookupCtrl := &controllers.LookupController{DB: database.DB}
	leaseCtrl := &controllers.LeaseController{DB: database.DB}
	supplierCtrl := &controllers.SupplierController{DB: database.DB}
	seizerCtrl := &controllers.SeizerController{DB: database.DB}
	introducerCtrl := &controllers.IntroducerController{DB: database.DB}
	valuationCompanyCtrl := &controllers.ValuationCompanyController{DB: database.DB}
	insuranceCompanyCtrl := &controllers.InsuranceCompanyController{DB: database.DB}
	auctionCompanyCtrl := &controllers.AuctionCompanyController{DB: database.DB}
	vehicleYardCtrl := &controllers.VehicleYardController{DB: database.DB}

	v1 := rg.Group("/v1")
	{
		leasing := v1.Group("/leasing")
		{
			products := leasing.Group("/products")
			{
				products.GET("", productCtrl.Index)                // GET  /api/v1/leasing/products
				products.POST("", productCtrl.Store)               // POST /api/v1/leasing/products
				products.GET("/:id", productCtrl.Get)              // GET  /api/v1/leasing/products/:id
				products.PUT("/:id", productCtrl.Update)           // PUT  /api/v1/leasing/products/:id
				products.GET("/:id/items", productCtrl.GetItems)   // GET  /api/v1/leasing/products/:id/items
				products.POST("/status", productCtrl.UpdateStatus) // POST /api/v1/leasing/products/status
			}

			leasing.POST("/calculate", leaseCtrl.CalculateSummary) // POST /api/v1/leasing/calculate
		}

		// Supplier routes
		suppliers := v1.Group("/suppliers")
		{
			suppliers.GET("", supplierCtrl.Index)          // GET /api/v1/suppliers
			suppliers.POST("", supplierCtrl.Store)         // POST /api/v1/suppliers
			suppliers.PUT("/:id", supplierCtrl.Update)     // PUT /api/v1/suppliers/:id
			suppliers.DELETE("/:id", supplierCtrl.Destroy) // DELETE /api/v1/suppliers/:id
		}

		seizers := v1.Group("/seizers")
		{
			seizers.GET("", seizerCtrl.Index)
			seizers.POST("", seizerCtrl.Store)
			seizers.PUT("/:id", seizerCtrl.Update)
			seizers.DELETE("/:id", seizerCtrl.Destroy)
		}

		introducers := v1.Group("/introducers")
		{
			introducers.GET("", introducerCtrl.Index)
			introducers.POST("", introducerCtrl.Store)
			introducers.PUT("/:id", introducerCtrl.Update)
			introducers.DELETE("/:id", introducerCtrl.Destroy)
		}

		valuationCompanies := v1.Group("/valuation-companies")
		{
			valuationCompanies.GET("", valuationCompanyCtrl.Index)
			valuationCompanies.POST("", valuationCompanyCtrl.Store)
			valuationCompanies.PUT("/:id", valuationCompanyCtrl.Update)
			valuationCompanies.DELETE("/:id", valuationCompanyCtrl.Destroy)
		}

		insuranceCompanies := v1.Group("/insuarance-companies")
		{
			insuranceCompanies.GET("", insuranceCompanyCtrl.Index)
			insuranceCompanies.POST("", insuranceCompanyCtrl.Store)
			insuranceCompanies.PUT("/:id", insuranceCompanyCtrl.Update)
			insuranceCompanies.DELETE("/:id", insuranceCompanyCtrl.Destroy)
		}

		auctionCompanies := v1.Group("/auction-companies")
		{
			auctionCompanies.GET("", auctionCompanyCtrl.Index)
			auctionCompanies.POST("", auctionCompanyCtrl.Store)
			auctionCompanies.PUT("/:id", auctionCompanyCtrl.Update)
			auctionCompanies.DELETE("/:id", auctionCompanyCtrl.Destroy)
		}

		vehicleYards := v1.Group("/vehicle-yards")
		{
			vehicleYards.GET("", vehicleYardCtrl.Index)
			vehicleYards.POST("", vehicleYardCtrl.Store)
			vehicleYards.PUT("/:id", vehicleYardCtrl.Update)
			vehicleYards.DELETE("/:id", vehicleYardCtrl.Destroy)
		}

		// Lookup routes
		lookups := v1.Group("/lookup")
		{
			lookups.GET("/banks", lookupCtrl.GetBanks)
			lookups.GET("/insurance-companies", lookupCtrl.GetInsuranceCompanies)
			lookups.GET("/vehicle-types", lookupCtrl.GetVehicleTypes)
			lookups.GET("/marketing-executives", lookupCtrl.GetMarketingExecutives)
		}

		// Customer routes
		customers := v1.Group("/customers")
		{
			customers.GET("", customerCtrl.Index)                         // GET /api/v1/customers
			customers.GET("/generate-id", customerCtrl.GenerateID)        // GET  /api/v1/customers/generate-id
			customers.GET("/search", customerCtrl.Search)                 // GET /api/v1/customers/search
			customers.POST("", customerCtrl.Store)                        // POST /api/v1/customers
			customers.GET("/:id", customerCtrl.Get)                       // GET /api/v1/customers/:id
			customers.POST("/:id/status", customerCtrl.UpdateStatus)      // POST /api/v1/customers/:id/status
			customers.POST("/:id/location", customerCtrl.UpdateLocation)  // POST /api/v1/customers/:id/location
			customers.GET("/:id/bank-accounts", customerCtrl.GetBankAccounts) // GET /api/v1/customers/:id/bank-accounts
			customers.POST("/:id/documents", customerCtrl.UploadDocument) // POST /api/v1/customers/:id/documents
		}

		// Leasing Applications (Stepper workflow)
		leasingAppCtrl := &controllers.LeasingApplicationController{DB: database.DB}
		leasingApplications := v1.Group("/leasing-applications")
		{
			leasingApplications.GET("/drafts", leasingAppCtrl.GetDrafts) // Needs to be before /:id
			leasingApplications.GET("/:id", leasingAppCtrl.Get)
			leasingApplications.POST("/draft", leasingAppCtrl.CreateDraft)
			leasingApplications.PUT("/:id/draft", leasingAppCtrl.UpdateDraft)
			leasingApplications.POST("/:id/upload-document", leasingAppCtrl.UploadDocument)
			leasingApplications.POST("/:id/submit", leasingAppCtrl.Submit)
		}

		// Location helper routes
		locations := v1.Group("/locations")
		{
			locations.GET("/cities", customerCtrl.GetCities) // GET /api/v1/locations/cities?province=X
		}
	}
}
