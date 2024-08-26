package controllers

import (
	"net/http"
	"strings"

	helpers "github.com/SiddharthaKR/golang-jwt-project/helpers"
	"github.com/gin-gonic/gin"
)

// EmailRequestBody defines the structure of the email request body
type EmailRequestBody struct {
    ToAddr  string `json:"to_addr"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// SendEmail handles sending emails
func SendEmail() gin.HandlerFunc {
    return func(c *gin.Context) {
        var reqBody EmailRequestBody
        if err := c.BindJSON(&reqBody); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
            return
        }

        // Convert comma-separated string to slice of strings
        to := strings.Split(reqBody.ToAddr, ",")

        if err := helpers.SendEmail(to, reqBody.Subject, reqBody.Body); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
    }
}
