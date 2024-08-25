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

func GetCustomers() gin.HandlerFunc {
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
