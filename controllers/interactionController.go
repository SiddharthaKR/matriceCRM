package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/SiddharthaKR/golang-jwt-project/database"
	helper "github.com/SiddharthaKR/golang-jwt-project/helpers"
	"github.com/SiddharthaKR/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var interactionCollection *mongo.Collection = database.OpenCollection(database.Client, "interaction")
var leadCollection *mongo.Collection = database.OpenCollection(database.Client, "lead")


func CreateLead() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var lead models.Lead

		// Bind JSON request body to the lead model
		if err := c.BindJSON(&lead); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the incoming data
		validationErr := validate.Struct(lead)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Assign a new ID and set timestamps
		lead.ID = primitive.NewObjectID()
		lead.CreatedAt = time.Now()
		lead.UpdatedAt = time.Now()

		// Insert the new lead into the database
		result, insertErr := leadCollection.InsertOne(ctx, lead)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while adding lead"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Lead created successfully", "lead_id": result.InsertedID})
	}
}

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


func GetInteractionReport() gin.HandlerFunc {
    return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        startDate := c.Query("start_date")
        endDate := c.Query("end_date")
        interactionType := c.Query("type")

        matchStage := bson.D{}
        if startDate != "" || endDate != "" {
            dateFilter := bson.M{}
            if startDate != "" {
                start, _ := time.Parse(time.RFC3339, startDate)
                dateFilter["$gte"] = start
            }
            if endDate != "" {
                end, _ := time.Parse(time.RFC3339, endDate)
                dateFilter["$lte"] = end
            }
            matchStage = append(matchStage, bson.E{"created_at", dateFilter})
        }
        if interactionType != "" {
            matchStage = append(matchStage, bson.E{"type", interactionType})
        }

        groupStage := bson.D{
            {"$group", bson.D{
                {"_id", bson.D{
                    {"type", "$type"},
                    {"status", "$status"},
                    {"day", bson.D{
                        {"$dateToString", bson.D{
                            {"format", "%Y-%m-%d"},
                            {"date", "$created_at"},
                        }},
                    }},
                }},
                {"count", bson.D{{"$sum", 1}}},
            }},
        }

        pipeline := mongo.Pipeline{bson.D{{"$match", matchStage}}, groupStage}
        cursor, err := interactionCollection.Aggregate(ctx, pipeline)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching interaction reports"})
            return
        }

        var results []bson.M
        if err = cursor.All(ctx, &results); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding interaction reports"})
            return
        }

        c.JSON(http.StatusOK, results)
    }
}


func GetConversionRateReport() gin.HandlerFunc {
    return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        startDate := c.Query("start_date")
        endDate := c.Query("end_date")

        matchStage := bson.D{}
        if startDate != "" || endDate != "" {
            dateFilter := bson.M{}
            if startDate != "" {
                start, _ := time.Parse(time.RFC3339, startDate)
                dateFilter["$gte"] = start
            }
            if endDate != "" {
                end, _ := time.Parse(time.RFC3339, endDate)
                dateFilter["$lte"] = end
            }
            matchStage = append(matchStage, bson.E{"created_at", dateFilter})
        }

        // Count total leads
        leadCount, err := leadCollection.CountDocuments(ctx, matchStage)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting leads"})
            return
        }

        // Count total customers
        customerCount, err := customerCollection.CountDocuments(ctx, matchStage)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting customers"})
            return
        }

        conversionRate := 0.0
        if leadCount > 0 {
            conversionRate = (float64(customerCount) / float64(leadCount)) * 100
        }

        c.JSON(http.StatusOK, gin.H{
            "total_leads":     leadCount,
            "total_customers": customerCount,
            "conversion_rate": conversionRate,
        })
    }
}

