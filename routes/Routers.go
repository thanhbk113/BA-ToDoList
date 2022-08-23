package routes

import (
	"todolist/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signupbk", controllers.SignUp())
	incomingRoutes.POST("/login", controllers.Login())
}

func TodoRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.POST("todo/:user_id", controllers.CreateToDo())
	incommingRoutes.GET("todos/:user_id", controllers.GetAllTodoList())
	incommingRoutes.PUT("todo/:user_id", controllers.UpdateATodo())
	incommingRoutes.DELETE("todo/:user_id", controllers.DeleteATodo())
}
