package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/SiddharthaKR/golang-jwt-project/database"
	"github.com/SiddharthaKR/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var interactionCollection *mongo.Collection = database.OpenCollection(database.Client, "interaction")

func CreateMeeting() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var interaction models.Interaction

		// Bind JSON request to interaction model
		if err := c.BindJSON(&interaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ensure CustomerID, UserID, and CompanyID are valid ObjectIDs
		customerID, err := primitive.ObjectIDFromHex(interaction.CustomerID.Hex())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Customer ID"})
			return
		}

		userID := c.GetString("uid")
		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
			return
		}
		interaction.UserID = userObjID

		// Get the company ID from the request parameters
        companyIDParam := c.Param("company_id")
        companyID, err := primitive.ObjectIDFromHex(companyIDParam)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Company ID"})
            return
        }

		// Check if the CompanyID exists
		var company models.Company
		err = companyCollection.FindOne(ctx, bson.M{"_id": companyID}).Decode(&company)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking company"})
			return
		}

		interaction.CustomerID = customerID
		interaction.CompanyID = companyID
		interaction.ID = primitive.NewObjectID()
		interaction.CreatedAt = time.Now()
		interaction.UpdatedAt = time.Now()
 
		result, err := interactionCollection.InsertOne(ctx, interaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while creating interaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Meeting created successfully", "interaction_id": result.InsertedID})
	}
}


func UpdateInteractionStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		interactionID := c.Param("interaction_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Define a struct to bind the request body
		var requestBody struct {
			Status string `json:"status" validate:"required,oneof=OPEN RESOLVED"`
		}

		// Bind JSON request to requestBody struct
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate status
		status := requestBody.Status
		if status != "OPEN" && status != "RESOLVED" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		// Convert interactionID to ObjectID
		interactionObjectID, err := primitive.ObjectIDFromHex(interactionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Interaction ID"})
			return
		}

		// Update interaction status
		result, err := interactionCollection.UpdateOne(
			ctx,
			bson.M{"_id": interactionObjectID},
			bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating interaction"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Interaction not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Interaction status updated successfully"})
	}
}


func GetCustomerInteractions() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerIDParam := c.Param("customer_id")
		customerID, err := primitive.ObjectIDFromHex(customerIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Customer ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var interactions []models.Interaction
		cursor, err := interactionCollection.Find(ctx, bson.M{"customerID": customerID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while retrieving interactions"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &interactions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding interactions"})
			return
		}

		c.JSON(http.StatusOK, interactions)
	}
}


func RaiseTicket() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var interaction models.Interaction

        // Bind JSON request to interaction model
        if err := c.BindJSON(&interaction); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Get the customer ID from the token (uid)
        customerID := c.GetString("uid")

        // Convert customer ID to ObjectID
        customerObjID, err := primitive.ObjectIDFromHex(customerID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Customer ID"})
            return
        }

        // Get the company ID from the request parameters
        companyIDParam := c.Param("company_id")
        companyID, err := primitive.ObjectIDFromHex(companyIDParam)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Company ID"})
            return
        }

		// Check if the CompanyID exists
		var company models.Company
		err = companyCollection.FindOne(ctx, bson.M{"_id": companyID}).Decode(&company)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking company"})
			return
		}

        interaction.CustomerID = customerObjID
        interaction.UserID = customerObjID // Since the customer is raising the ticket
        interaction.CompanyID = companyID
        interaction.ID = primitive.NewObjectID()
        interaction.Type = "TICKET"
        interaction.Status = "OPEN"
        interaction.CreatedAt = time.Now()
        interaction.UpdatedAt = time.Now()

        // Insert the ticket into the database
        result, err := interactionCollection.InsertOne(ctx, interaction)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while raising the ticket"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Ticket raised successfully", "ticket_id": result.InsertedID})
    }
}

