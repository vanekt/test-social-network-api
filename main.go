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
	wsPort := os.Getenv("WS_PORT")
	logger := NewLogger("test-social-network-api", logging.DEBUG)
	db := sqlx.MustConnect("mysql", os.Getenv("SQL_DB_DSN"))

	// init models
	userModel := model.NewUserModel(db, logger)
	messagesModel := model.NewMessagesModel(db, logger)

	// init controllers
	authController := controller.NewAuthController(logger, userModel)
	userController := controller.NewUserController(logger, userModel)
	messagesController := controller.NewMessagesController(logger, messagesModel)

	// init gin
	r := gin.Default()
	configureRouter(r, authController, userController, messagesController)

	ws := NewWebsocket(logger, messagesModel)
	go ws.Run(":" + wsPort)

	r.Run(":" + port)
}
