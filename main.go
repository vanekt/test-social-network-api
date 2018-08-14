package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"os"
	"vanekt/test-social-network-api/controller"
	"vanekt/test-social-network-api/model"
)

func main() {
	port := os.Getenv("PORT")
	logger := NewLogger("test-social-network-api", logging.DEBUG)
	db := sqlx.MustConnect("mysql", os.Getenv("SQL_DB_DSN"))

	// init models
	userModel := model.NewUserModel(db, logger)

	// init controllers
	authController := controller.NewAuthController(logger, userModel)
	userController := controller.NewUserController(logger, userModel)

	// init gin
	r := gin.Default()
	configureRouter(r, authController, userController)

	r.Run(":" + port)
}
