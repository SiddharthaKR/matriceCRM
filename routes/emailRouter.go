package routes

import (
    controller "github.com/SiddharthaKR/golang-jwt-project/controllers"
    "github.com/SiddharthaKR/golang-jwt-project/middleware"
    "github.com/gin-gonic/gin"
)

func EmailRoutes(incomingRoutes *gin.Engine) {
    incomingRoutes.Use(middleware.Authenticate())
    incomingRoutes.POST("/email", controller.SendEmail())  // Route for sending emails
}