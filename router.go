package main

import (
	"github.com/gin-gonic/gin"
	"vanekt/test-social-network-api/controller"
)

func configureRouter(
	r *gin.Engine,
	authController *controller.AuthController,
) {
	v1 := r.Group("/api/v1")

	// -------------------------------
	auth := v1.Group("/auth").Use()

	loginHandler := authController.Login()
	auth.POST("/login", loginHandler)

	logoutHandler := authController.Logout()
	auth.GET("/logout", logoutHandler)
}
