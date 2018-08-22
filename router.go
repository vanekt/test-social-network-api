package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vanekt/test-social-network-api/controller"
)

func configureRouter(
	r *gin.Engine,
	authController *controller.AuthController,
	userController *controller.UserController,
	messagesController *controller.MessagesController,
) {
	v1 := r.Group("/api/v1")

	// -------------------------------
	auth := v1.Group("/auth")
	auth.POST("/login", authController.Login)
	auth.GET("/logout", authController.Logout)
	auth.GET("/check", authController.CheckAuth)

	// -------------------------------
	users := v1.Group("/users")
	users.GET("/", userController.GetAll)
	users.GET("/:id", userController.GetUserById)

	// -------------------------------
	dialogs := v1.Group("/dialogs")
	dialogs.GET("/:id", messagesController.GetDialogsByUserId)

	// -------------------------------
	messages := v1.Group("/messages")
	messages.GET("/:userId/:peerId", messagesController.GetDialogMessages)
}
