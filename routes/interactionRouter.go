package routes

import (
	controller "github.com/SiddharthaKR/golang-jwt-project/controllers"
	"github.com/SiddharthaKR/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func InteractionRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/interactions/:company_id/ticket", controller.RaiseTicket())
	incomingRoutes.POST("/interactions/:company_id/meeting", controller.CreateMeeting())
    incomingRoutes.PUT("/interactions/:interaction_id/status", controller.UpdateInteractionStatus())
    incomingRoutes.GET("/customers/:customer_id/interactions", controller.GetCustomerInteractions())
}
