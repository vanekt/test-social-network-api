package main

import (
	"github.com/gin-gonic/gin"
	"vanekt/test-social-network-api/controller"
)

func configureRouter(
	r *gin.Engine,
	authController *controller.AuthController,
	userController *controller.UserController,
	messagesController *controller.MessagesController,
) {
	v1 := r.Group("/api/v1")

	// -------------------------------
	auth := v1.Group("/auth").Use()

	loginHandler := authController.Login()
	auth.POST("/login", loginHandler)

	logoutHandler := authController.Logout()
	auth.GET("/logout", logoutHandler)

	checkAuthHandler := authController.CheckAuth()
	auth.GET("/check", checkAuthHandler)

	// -------------------------------
	users := v1.Group("/users").Use()

	getUserHandler := userController.GetUserById()
	users.GET("/:id", getUserHandler)

	// -------------------------------
	dialogs := v1.Group("/dialogs").Use()

	getDialogsHandler := messagesController.GetDialogsByUserId()
	dialogs.GET("/:id", getDialogsHandler)

	// -------------------------------
	messages := v1.Group("/messages").Use()

	getDialogMessagesHandler := messagesController.GetDialogMessages()
	messages.GET("/:userId/:peerId", getDialogMessagesHandler)
}
