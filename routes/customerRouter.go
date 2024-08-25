package routes

import (
	controller "github.com/akhil/golang-jwt-project/controllers"
	"github.com/akhil/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

///// will contain basic meets, call etc....

func CustomerRoutes(incomingRoutes *gin.Engine) {
	// Apply authentication middleware to all customer routes
	incomingRoutes.Use(middleware.Authenticate())

	// Define routes for customer operations
	incomingRoutes.GET("/customers", controller.GetCustomers())               // Get list of all customers
	incomingRoutes.GET("/customers/:customer_id", controller.GetCustomer())   // Get a specific customer by ID
	// incomingRoutes.POST("/customers", controller.CreateCustomer())            // Create a new customer
	// incomingRoutes.PUT("/customers/:customer_id", controller.UpdateCustomer()) // Update an existing customer
	// incomingRoutes.DELETE("/customers/:customer_id", controller.DeleteCustomer()) // Delete a customer by ID
}
