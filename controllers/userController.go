package controllers

import (
	"context"
	"log"
	"net/http"

	// "strconv"
	"time"

	"github.com/PranavMasekar/go-jwt-project/database"
	"github.com/PranavMasekar/go-jwt-project/helpers"
	"github.com/PranavMasekar/go-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "golang.org/x/crypto/bcrypt"
)

// Getting the collection of user
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func VerifyPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(c, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking emails"})
			return
		}

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
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {

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
