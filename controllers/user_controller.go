package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"todolist/configs"
	"todolist/helpers"
	"todolist/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users") //initial mongo CLient and create database name is todos

//HashPassword is used to encrypt password before stored is in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) //helloword --hashed-->$2a$14$IzdvYkAsnJPF8L/Ok.giyug/j/qDg10YqjIpB6X3hm9kCztJLtmfy
	if err != nil {
		log.Fatal("error hash password:", err)
	}

	return string(bytes)
}

// VerifyPassword checks the input password while veryfying it with the password in the DB.
//VeryfyPassword("Input Password","Password from database hashed")
func VeryfyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = "user or password is incorrect"
		check = false
	}
	//check, msg := VeryfyPassword("helloword", "$2a$14$IzdvYkAsnJPF8L/Ok.giyug/j/qDg10YqjIpB6X3hm9kCztJLtmfy")
	//=>check=true && msg=""
	return check, msg
}

//CreateUser is used to create new user
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil { // check json empty send from client
			c.JSON(http.StatusBadRequest, gin.H{"error signup": err.Error()})
			defer cancel()
			return
		}

		//After BindJSON user become pointe
		//You can access value using *user.value ex:*user.email

		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error json input": validationErr.Error()})
			defer cancel()
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Fatal("error email:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "error occured when cheking email"})
			return
		}
		if count > 0 {
			log.Fatal("error email", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This email already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, err := helpers.GenerateAllTokens(user.User_id)

		if err != nil {
			log.Fatal("error when generate token")
			c.JSON(http.StatusInternalServerError, gin.H{"error when generate token": err.Error()})
			return
		}

		user.Token = &token
		user.Refresh_token = &refreshToken

		// c.SetCookie("access_token", token, 3600, "/", "https://todosbk.netlify.app/", true, true)
		// c.SetCookie("refresh_token", refreshToken, 3600, "/", "https://todosbk.netlify.app/", true, true)

		_, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error when create user": insertErr.Error()})
			return
		}

		defer cancel()
		accessCookie, _ := c.Cookie("access_token")
		fmt.Println("access_token:", accessCookie)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

//Login is the api used to get a single user

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error bind json": err.Error()})
			defer cancel()
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		passwordIsValid, msg := VeryfyPassword(*user.Password, *foundUser.Password) // if password not match
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(foundUser.User_id)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, foundUser) // send json to client
	}
}
