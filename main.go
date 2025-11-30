package main

import (
	"Go/controllers/auth"
	"Go/controllers/post"
	"Go/controllers/user"
	"Go/initializers"
	"Go/middlewares"
	"Go/services"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvironmentVariables()
	initializers.ConnectToDatabase()
}

func main() {

	r := gin.Default()

	userService := services.NewUserService(initializers.DB)
	postService := services.NewPostService(initializers.DB)

	authController := auth.NewAuthController(userService)
	postController := post.NewPostController(postService)

	//Auth
	r.POST("/auth/register", authController.Register)
	r.POST("/auth/login", authController.Login)

	api := r.Group("/api")

	//User
	api.Use(middlewares.JWTAuthMiddleware())
	{
		api.GET("/user/:id/view", user.View)
		api.GET("/user/list", user.List)
		api.POST("/user/:id/delete", user.Delete)
		api.POST("/user/create", user.Create)
	}

	//Posts
	api.GET("/post/view/:uuid", postController.View)

	api.Use(middlewares.JWTAuthMiddleware())
	{
		api.POST("/post/create", postController.Create)
		api.POST("/post/delete/:uuid", postController.Delete)
		api.POST("/post/update", postController.Update)
		api.POST("/post/list", postController.List)
	}

	r.Run(":6969")
}
