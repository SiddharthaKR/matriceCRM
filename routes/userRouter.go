package routes

import (
	controller "github.com/SiddharthaKR/golang-jwt-project/controllers"
	"github.com/SiddharthaKR/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetAllUsers())
	incomingRoutes.PUT("/users/:user_id", controller.UpdateUser())
	incomingRoutes.DELETE("/users/:user_id", controller.DeleteUser())  // Delete a user by ID
}
