package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"net/http"
	"os"
	"strconv"
	"vanekt/test-social-network-api/entity"
	"vanekt/test-social-network-api/error"
	"vanekt/test-social-network-api/model"
	"vanekt/test-social-network-api/util"
)

type MessagesController struct {
	logger        *logging.Logger
	messagesModel *model.MessagesModel
}

func NewMessagesController(logger *logging.Logger, messagesModel *model.MessagesModel) *MessagesController {
	return &MessagesController{logger, messagesModel}
}

func (c *MessagesController) GetDialogsByUserId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strId := ctx.Param("id")
		userId, err := strconv.Atoi(strId)
		if err != nil {
			c.logger.Errorf("[MessagesController GetDialogsByUserId] parse request userId error: %v", err.Error())
			ctx.JSON(http.StatusBadRequest, error.HttpResponseErrorBadRequest)
			return
		}

		tokenString, err := ctx.Cookie(os.Getenv("TOKEN_COOKIE_NAME"))
		if err != nil {
			c.logger.Errorf("[MessagesController GetDialogsByUserId] Cannot get token from cookies. Trace %s", err.Error())
			ctx.JSON(http.StatusUnauthorized, error.HttpResponseErrorBadRequest)
			return
		}

		authUserId, err := util.GetUserIdFromToken(tokenString)
		if err != nil {
			c.logger.Errorf("[MessagesController GetDialogsByUserId] Cannot fetch authUserId from token. Trace %s", err.Error())
			ctx.JSON(http.StatusUnauthorized, error.HttpResponseErrorBadRequest)
			return
		}

		if userId != int(authUserId) {
			c.logger.Error("[MessagesController GetDialogsByUserId] userId != authUserId")
			ctx.JSON(http.StatusForbidden, error.HttpResponseErrorForbidden)
			return
		}

		dialogs := make([]entity.Dialog, 0)
		dialogs, err = c.messagesModel.GetDialogsByUserId(userId)
		if err != nil {
			c.logger.Errorf("[MessagesController GetDialogsByUserId] messagesModel.GetDialogsByUserId err. Trace %s", err.Error())
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"payload": dialogs,
		})
		return
	}
}
