package routes

import(
	controller "github.com/akhil/golang-jwt-project/controllers"
	"github.com/akhil/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	// incomingRoutes.POST("/users", controller.CreateUser())             // Create a new user
	// incomingRoutes.PUT("/users/:user_id", controller.UpdateUser())     // Update an existing user
	// incomingRoutes.DELETE("/users/:user_id", controller.DeleteUser())  // Delete a user by ID
}