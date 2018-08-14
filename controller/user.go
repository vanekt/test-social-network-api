package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"net/http"

	"database/sql"
	"strconv"
	"vanekt/test-social-network-api/error"
	"vanekt/test-social-network-api/model"
)

type UserController struct {
	logger    *logging.Logger
	userModel *model.UserModel
}

func NewUserController(logger *logging.Logger, userModel *model.UserModel) *UserController {
	return &UserController{logger, userModel}
}

func (c *UserController) GetUserById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
}
