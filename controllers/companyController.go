package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/akhil/golang-jwt-project/models"
	"github.com/akhil/golang-jwt-project/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var companyCollection *mongo.Collection = database.OpenCollection(database.Client, "company")

// CreateCompany creates a new company and updates user company IDs
func CreateCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
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
		defer cancel()
		if insertErr != nil {
			msg := fmt.Sprintf("Company item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Update user documents to include the new company ID
		_, err := userCollection.UpdateMany(
			ctx,
			bson.M{},
			bson.M{"$push": bson.M{"company_ids": company.ID}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating users with the new company ID"})
			return
		}

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

// DeleteCompany removes a company and updates user company IDs
func DeleteCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID := c.Param("company_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Remove the company
		result, err := companyCollection.DeleteOne(ctx, bson.M{"_id": companyID})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deleting company"})
			return
		}

		// Remove the company ID from all users
		_, err = userCollection.UpdateMany(
			ctx,
			bson.M{},
			bson.M{"$pull": bson.M{"company_ids": companyID}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating users after company deletion"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetCompany retrieves a single company by ID
func GetCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID := c.Param("company_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var company models.Company
		err := companyCollection.FindOne(ctx, bson.M{"_id": companyID}).Decode(&company)
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