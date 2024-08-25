package controllers

import (
	"context"
	"fmt"
	"github.com/SiddharthaKR/golang-jwt-project/database"
	helper "github.com/SiddharthaKR/golang-jwt-project/helpers"
	"github.com/SiddharthaKR/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"time"
)

var customerCollection *mongo.Collection = database.OpenCollection(database.Client, "customer")
var validate = validator.New()

type CustomerResponse struct {
	FirstName       *string `json:"first_name"`
	LastName        *string `json:"last_name"`
	Email           *string `json:"email"`
	Phone           *string `json:"phone"`
	Company         *string `json:"company,omitempty"`
	Status          *string `json:"status"`
	Notes           *string `json:"notes,omitempty"`
	CustomerID      string  `json:"customer_id"`
	LastInteraction time.Time `json:"last_interaction,omitempty"`
	CompanyID       string  `json:"company_id"`
}

func CustomerSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var customer models.Customer

		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(customer)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := customerCollection.CountDocuments(ctx, bson.M{"email": customer.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			return
		}

		password := helper.HashPassword(*customer.PasswordHash)
		customer.PasswordHash = &password

		count, err = customerCollection.CountDocuments(ctx, bson.M{"phone": customer.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the phone number"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}

		// Handle CompanyID
		if customer.CompanyID.IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CompanyID is required"})
			return
		}

		customer.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		customer.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		customer.ID = primitive.NewObjectID()
		customer.CustomerID = customer.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*customer.Email, *customer.FirstName, *customer.LastName, "CUSTOMER", *&customer.CustomerID)
		customer.Token = &token
		customer.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := customerCollection.InsertOne(ctx, customer)
		if insertErr != nil {
			msg := fmt.Sprintf("Customer item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func CustomerLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var customer models.Customer
		var foundCustomer models.Customer

		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := customerCollection.FindOne(ctx, bson.M{"email": customer.Email}).Decode(&foundCustomer)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := helper.VerifyPassword(*customer.PasswordHash, *foundCustomer.PasswordHash)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundCustomer.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "customer not found"})
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundCustomer.Email, *foundCustomer.FirstName, *foundCustomer.LastName, "CUSTOMER", foundCustomer.CustomerID)
		helper.UpdateAllTokens(token, refreshToken, foundCustomer.CustomerID)
		err = customerCollection.FindOne(ctx, bson.M{"customer_id": foundCustomer.CustomerID}).Decode(&foundCustomer)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundCustomer)
	}
}


func GetCustomersByCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("uid")
		companyIDParam := c.Param("company_id")

		companyID, err := primitive.ObjectIDFromHex(companyIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
			return
		}

		// Check if the user has access to this company
		if !checkUserAccessToCompany(userID, companyID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this company"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var customers []models.Customer
		cursor, err := customerCollection.Find(ctx, bson.M{"companyID": companyID})
		if err != nil {
			log.Println("Error finding customers:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while retrieving customers"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &customers); err != nil {
			log.Println("Error decoding customers:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding customers"})
			return
		}

		// Transform the customer data
		var response []CustomerResponse
		for _, customer := range customers {
			response = append(response, CustomerResponse{
				FirstName:       customer.FirstName,
				LastName:        customer.LastName,
				Email:           customer.Email,
				Phone:           customer.Phone,
				Company:         customer.Company,
				Status:          customer.Status,
				Notes:           customer.Notes,
				CustomerID:      customer.CustomerID,
				LastInteraction: customer.LastInteraction,
				CompanyID:       customer.CompanyID.Hex(), // Convert ObjectID to string for response
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCompanyCustomerByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("uid")
		companyIDParam := c.Param("company_id")
		customerIDParam := c.Param("customer_id")

		// Convert companyID and customerID to ObjectID
		companyID, err := primitive.ObjectIDFromHex(companyIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
			return
		}

		customerID, err := primitive.ObjectIDFromHex(customerIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
			return
		}

		// Check if the user has access to this company
		if !checkUserAccessToCompany(userID, companyID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this company"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var customer models.Customer
		err = customerCollection.FindOne(ctx, bson.M{"_id": customerID, "companyID": companyID}).Decode(&customer)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			} else {
				log.Println("Error finding customer:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while retrieving customer"})
			}
			return
		}

		// Transform the customer data
		response := CustomerResponse{
			FirstName:       customer.FirstName,
			LastName:        customer.LastName,
			Email:           customer.Email,
			Phone:           customer.Phone,
			Company:         customer.Company,
			Status:          customer.Status,
			Notes:           customer.Notes,
			CustomerID:      customer.CustomerID,
			LastInteraction: customer.LastInteraction,
			CompanyID:       customer.CompanyID.Hex(), // Convert ObjectID to string for response
		}

		c.JSON(http.StatusOK, response)
	}
}


func UpdateCompanyCustomerByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		companyIDParam := c.Param("company_id")
		customerIDParam := c.Param("customer_id")

		companyID, err := primitive.ObjectIDFromHex(companyIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
			return
		}

		customerID, err := primitive.ObjectIDFromHex(customerIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
			return
		}

		var updatedData models.Customer
		if err := c.BindJSON(&updatedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the user has access to this company
		userID := c.GetString("uid")
		if !checkUserAccessToCompany(userID, companyID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this company"})
			return
		}

		update := bson.M{}
		if updatedData.FirstName != nil {
			update["first_name"] = updatedData.FirstName
		}
		if updatedData.LastName != nil {
			update["last_name"] = updatedData.LastName
		}
		if updatedData.Email != nil {
			update["email"] = updatedData.Email
		}
		if updatedData.Phone != nil {
			update["phone"] = updatedData.Phone
		}
		if updatedData.Status != nil {
			update["status"] = updatedData.Status
		}
		if updatedData.Notes != nil {
			update["notes"] = updatedData.Notes
		}
		if !updatedData.CompanyID.IsZero() {
			update["company_id"] = updatedData.CompanyID
		}

		// Perform the update
		result, err := customerCollection.UpdateOne(
			ctx,
			bson.M{"_id": customerID},
			bson.M{"$set": update},
		)
		if err != nil {
			log.Println("Error updating customer:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating customer"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
	}
}


func DeleteComapnyCustomerByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("uid")
		companyIDParam := c.Param("company_id")
		customerIDParam := c.Param("customer_id")

		// Convert companyID and customerID to ObjectID
		companyID, err := primitive.ObjectIDFromHex(companyIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
			return
		}

		customerID, err := primitive.ObjectIDFromHex(customerIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
			return
		}

		// Check if the user has access to this company
		if !checkUserAccessToCompany(userID, companyID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this company"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Delete the customer
		deleteResult, err := customerCollection.DeleteOne(ctx, bson.M{"_id": customerID, "companyID": companyID})
		if err != nil {
			log.Println("Error deleting customer:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deleting customer"})
			return
		}

		if deleteResult.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
	}
}


func GetAllCustomers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"customer_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := customerCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing customer items"})
			return
		}
		var allCustomers []bson.M
		if err = result.All(ctx, &allCustomers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allCustomers[0])
	}
}

func GetCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerId := c.Param("customer_id")

		if err := helper.MatchUserTypeToUid(c, customerId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var customer models.Customer
		err := customerCollection.FindOne(ctx, bson.M{"customer_id": customerId}).Decode(&customer)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, customer)
	}
}


// Check if the user is an Admin or if the company ID is in the user's company ID list.
func checkUserAccessToCompany(userID string, companyID primitive.ObjectID) bool {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		log.Println("Error finding user with ID:", userID, "Error:", err)
		return false
	}

	// Check if user is an Admin
	if *user.UserType == "ADMIN" {
		log.Println("User is an Admin, access granted.")
		return true
	}

	// Check if the company ID is in the user's CompanyIDs list
	for _, id := range user.CompanyIDs {
		if id == companyID {
			log.Println("User has access to company ID:", companyID)
			return true
		}
	}

	// Log if the company ID is not found in the user's CompanyIDs list
	log.Println("User does not have access to company ID:", companyID)
	return false
}
