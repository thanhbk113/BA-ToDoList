package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"todolist/configs"
	"todolist/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var todoListCollection *mongo.Collection = configs.GetCollection(configs.DB, "todoList") //initial mongo CLient and create database name is todos

var validate = validator.New()

func CreateToDo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var todo models.Todo // todo use to get from client sent
		var user models.User //user to check user exist or not

		var userId = c.Param("user_id")
		if err := c.BindJSON(&todo); err != nil { // check get json client send success or not
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		//create new todo get from todo sent from client
		newTodo := models.Todo{
			Id:         primitive.NewObjectID(),
			Task:       todo.Task,
			Done:       todo.Done,
			Created_at: todo.Created_at,
		}
		//validtaion json
		if validationErr := validate.Struct(&newTodo); validationErr != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		//check user id exist
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		defer cancel()
		//if todo list id not exist
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user id not exist"})
			return
		}

		var foundTodoList models.TodoList

		err = todoListCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundTodoList)
		if err != nil { //check if todo not exist when user create first then create one
			newTodoList := models.TodoList{
				User_id:    userId,
				Created_at: todo.Created_at,
			}
			todoListCollection.InsertOne(ctx, newTodoList)
			query := bson.M{"user_id": userId}
			update := bson.M{"$push": bson.M{"todo_list": newTodo}}

			todoListCollection.UpdateOne(ctx, query, update)
			c.JSON(http.StatusOK, gin.H{"message": "create success new todo list for your account"})
		} else {
			// //query
			query := bson.M{"user_id": userId}
			update := bson.M{"$push": bson.M{"todo_list": newTodo}}

			todoListCollection.UpdateOne(ctx, query, update)
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		}

	}
}

func GetAllTodoList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		userId := c.Param("user_id")
		var todo_list struct {
			User_id    string        `json:"user_id"` //related SQL
			Created_at time.Time     `json:"created_at"`
			Todos      []models.Todo `bson:"todo_list"`
		}

		defer cancel()

		err := todoListCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&todo_list)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error user not exist": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": todo_list})
	}
}

func UpdateATodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		todoId := c.Param("todo_id")
		var todo models.TodoList
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error json input": err.Error()})
		}
		objId, _ := primitive.ObjectIDFromHex(todoId)
		fmt.Println(todoId, "-->", objId)
		defer cancel()

		_, err := todoListCollection.UpdateOne(ctx, bson.M{"todo_list.task": "Drinking Tea"}, bson.M{"$set": todo})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error update todo": err.Error()})
			fmt.Println(err.Error())
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, todo)

	}
}