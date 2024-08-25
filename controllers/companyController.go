package controllers

import (
	"context"
	// "fmt"
	"log"
	"net/http"
	"time"

	"github.com/SiddharthaKR/golang-jwt-project/database"
	"github.com/SiddharthaKR/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var companyCollection *mongo.Collection = database.OpenCollection(database.Client, "company")

// CreateCompany creates a new company and updates user company IDs
func CreateCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from context
		userID, exists := c.Get("uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var company models.Company

		if err := c.BindJSON(&company); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		company.CreatedAt = time.Now()
		company.UpdatedAt = time.Now()
		company.ID = primitive.NewObjectID()

		// Insert the company
		resultInsertionNumber, insertErr := companyCollection.InsertOne(ctx, company)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Company item was not created"})
			return
		}

		// Update the current user's document to include the new company ID
		_, err := userCollection.UpdateOne(
			ctx,
			bson.M{"user_id": userID},
			bson.M{"$push": bson.M{"CompanyIDs": company.ID}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating user with the new company ID"})
			return
		}

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}


// DeleteCompany removes a company and updates user company IDs
func DeleteCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		companyIDParam := c.Param("company_id") // Keep this as a string
		userID := c.GetString("uid")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Convert companyIDParam to ObjectID
		companyObjectID, err := primitive.ObjectIDFromHex(companyIDParam) // Convert to ObjectID
		if err != nil {
			log.Println("Invalid company ID format:", companyIDParam)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
			return
		}

		// Check if the user has access to this company
		if !checkUserAccessToCompany(userID, companyObjectID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this company"})
			return
		}

		// Attempt to delete the company using the ObjectID
		result, err := companyCollection.DeleteOne(ctx, bson.M{"_id": companyObjectID})
		if err != nil {
			log.Println("Error occurred while deleting company:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deleting company"})
			return
		}

		if result.DeletedCount == 0 {
			log.Println("No company found with ID:", companyIDParam)
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}

		// Remove the company ID from all users
		updateResult, err := userCollection.UpdateMany(
			ctx,
			bson.M{},
			bson.M{"$pull": bson.M{"CompanyIDs": companyObjectID}}, // Use companyObjectID here
		)
		if err != nil {
			log.Println("Error occurred while updating users after company deletion:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating users after company deletion"})
			return
		}

		log.Println("Company deleted successfully. Company ID:", companyIDParam)
		log.Println("Users updated. Matched count:", updateResult.MatchedCount, "Modified count:", updateResult.ModifiedCount)

		c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully", "deleted_count": result.DeletedCount})
	}
}


// GetCompany retrieves a single company by ID
func GetCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID := c.Param("company_id")
		userID := c.GetString("uid")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Convert companyIDParam to ObjectID
		companyObjectID, err := primitive.ObjectIDFromHex(companyID) // Convert to ObjectID
		if err != nil {
			log.Println("Invalid company ID format:", companyID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
			return
		}

        // Check if the user has access to this company
		if !checkUserAccessToCompany(userID, companyObjectID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this company"})
			return
		}
		var company models.Company
		err = companyCollection.FindOne(ctx, bson.M{"_id": companyID}).Decode(&company)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching company"})
			return
		}
		c.JSON(http.StatusOK, company)
	}
}

// UpdateCompany updates a company's details
func UpdateCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID := c.Param("company_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var company models.Company

		if err := c.BindJSON(&company); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := bson.M{
			"$set": company,
		}

		result, err := companyCollection.UpdateOne(ctx, bson.M{"_id": companyID}, update)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while updating company"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetCompanies retrieves a list of companies
func GetCompanies() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var companies []models.Company
		cursor, err := companyCollection.Find(ctx, bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing companies"})
			return
		}

		if err = cursor.All(ctx, &companies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, companies)
	}
}
