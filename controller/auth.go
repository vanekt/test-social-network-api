package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"net/http"
	"os"
	"vanekt/test-social-network-api/error"
	"vanekt/test-social-network-api/model"
	"vanekt/test-social-network-api/util"
)

type AuthController struct {
	logger    *logging.Logger
	userModel *model.UserModel
}

func NewAuthController(logger *logging.Logger, userModel *model.UserModel) *AuthController {
	return &AuthController{logger, userModel}
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *AuthController) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req LoginRequest
		if err := ctx.BindJSON(&req); err != nil {
			c.logger.Errorf("[AuthController Login] parse request data error: %v", err.Error())
			ctx.JSON(http.StatusBadRequest, error.HttpResponseErrorBadRequest)
			return
		}

		user, err := c.userModel.GetUserByLogin(req.Login)
		if err != nil {
			if err == sql.ErrNoRows {
				c.logger.Errorf("[AuthController Login] User not found for login=%s", req.Login)
				ctx.JSON(http.StatusNotFound, error.HttpResponseErrorNotFound)
				return
			}
			c.logger.Errorf("[AuthController Login] GetUserByLogin error: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, error.HttpResponseErrorInternalServerError)
			return
		}

		if ok := util.CheckPasswordHash(req.Password, user.Password); !ok {
			c.logger.Errorf("[AuthController Login] Wrong password for userId=%d", user.Id)
			ctx.JSON(http.StatusNotFound, error.HttpResponseErrorNotFound)
			return
		}

		tokenString, err := util.CreateToken(user.Id)
		if err != nil {
			c.logger.Errorf("[AuthController Login] Cannot create token. Trace %s", err.Error())
			ctx.JSON(http.StatusInternalServerError, error.HttpResponseErrorInternalServerError)
			return
		}

		ctx.SetCookie(
			os.Getenv("TOKEN_COOKIE_NAME"), // name
			tokenString,                    // value
			0,                              // maxAge // TODO
			"/",                            // path
			os.Getenv("TOKEN_COOKIE_DOMAIN"), // domain
			false, // secure
			true,  // httpOnly
		)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
		return
	}
}

func (c *AuthController) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetCookie(
			os.Getenv("TOKEN_COOKIE_NAME"), // name
			"",  // value
			-1,  // maxAge
			"/", // path
			os.Getenv("TOKEN_COOKIE_DOMAIN"), // domain
			false, // secure
			true,  // httpOnly
		)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
		return
	}
}

func (c *AuthController) CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie(os.Getenv("TOKEN_COOKIE_NAME"))
		if err != nil {
			c.logger.Errorf("[AuthController CheckAuth] Cannot get token from cookies. Trace %s", err.Error())
			ctx.JSON(http.StatusOK, gin.H{
				"success": false,
			})
			return
		}

		userId, err := util.GetUserIdFromToken(tokenString)
		if err != nil {
			c.logger.Errorf("[AuthController CheckAuth] Cannot fetch userId from token. Trace %s", err.Error())
			ctx.JSON(http.StatusOK, gin.H{
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"payload": userId,
		})
		return
	}
}