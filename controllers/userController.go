package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PranavMasekar/go-jwt-project/database"
	"github.com/PranavMasekar/go-jwt-project/helpers"
	"github.com/PranavMasekar/go-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Getting the collection of user
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email of password is incorrect")
		check = false
	}
	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		// SignUp means we have to create user so get the request body and convert to struct user
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate the new user with validate properties we set in model
		validationErr := validate.Struct(user)

		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		// Checking if user aleady exits in DB via Email
		count, err := userCollection.CountDocuments(c, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking emails"})
			return
		}
		// Hashing the password to store in db
		password := HashPassword(*user.Password)
		user.Password = &password

		// Checking if user aleady exits in DB via Phone
		count, err = userCollection.CountDocuments(c, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking phone numbers"})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number alredy exits"})
		}
		// Complete the user model with rokens and time stamps
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken
		// Adding user to DB
		resultInsertionNumber, insertedError := userCollection.InsertOne(c, user)

		if insertedError != nil {
			msg := fmt.Sprintf("user item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		// Response
		ctx.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User
		// Body of request
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Find the collection with email and store in foundUser
		err := userCollection.FindOne(c, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password is incorrect"})
			return
		}
		// Verify the password
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !passwordIsValid {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		// Generate the tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(c, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Only admin are allowed
		if err := helpers.CheckUserType(ctx, "ADMIN"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recoredPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))
		if err != nil || recoredPerPage < 1 {
			recoredPerPage = 10
		}
		page, err1 := strconv.Atoi(ctx.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recoredPerPage

		startIndex, _ = strconv.Atoi(ctx.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		// Group the unique id in one group and count them
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recoredPerPage}}}}}}}

		result, err := userCollection.Aggregate(c, mongo.Pipeline{matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}

		var allUsers []bson.M
		if err = result.All(c, &allUsers); err != nil {
			log.Fatal(err)
		}

		ctx.JSON(http.StatusOK, allUsers[0])
	}
}
func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Getting the user Id
		userId := ctx.Param("user_id")
		// Matching the user id and user type
		if err := helpers.MatchUserTypeToUid(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Interacting with db
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		// Get particular user from db
		err := userCollection.FindOne(c, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}
