package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"net/http"

	"database/sql"
	"github.com/vanekt/test-social-network-api/error"
	"github.com/vanekt/test-social-network-api/model"
	"strconv"
)

type UserController struct {
	logger    *logging.Logger
	userModel *model.UserModel
}

func NewUserController(logger *logging.Logger, userModel *model.UserModel) *UserController {
	return &UserController{logger, userModel}
}

func (c *UserController) GetUserById(ctx *gin.Context) {
	strId := ctx.Param("id")
	userId, err := strconv.Atoi(strId)
	if err != nil {
		c.logger.Errorf("[UserController GetUserById] parse request userId error: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, error.HttpResponseErrorBadRequest)
		return
	}

	user, err := c.userModel.GetUserById(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.logger.Errorf("[UserController GetUserById] User not found: id=%d", userId)
			ctx.JSON(http.StatusNotFound, error.HttpResponseErrorNotFound)
			return
		}
		c.logger.Errorf("[UserController GetUserById] GetUserById error: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, error.HttpResponseErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"payload": user,
	})
	return
}

func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.userModel.GetAll()
	if err != nil {
		c.logger.Errorf("[UserController GetAll] error: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, error.HttpResponseErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"payload": users,
	})
	return
}
