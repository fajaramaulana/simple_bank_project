package controller

import (
	"log"
	"net/http"

	"github.com/fajaramaulana/simple_bank_project/internal/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (a *AuthController) Login(ctx *gin.Context) {
	var req request.AuthLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		massage, data := helper.GlobalCheckingErrorBindJson(err.Error(), req)
		log.Printf("Error Bind: %s", massage)
		helper.ReturnJSONError(ctx, http.StatusBadRequest, massage, nil, data)
		return
	}

	responseService, err := a.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		log.Printf("Error: %s", err.Error())

		if err == service.ErrUserNotFound {
			helper.ReturnJSONError(ctx, http.StatusNotFound, "User not found", nil, nil)
			return
		} else if err == service.ErrorInvalidPassword {
			helper.ReturnJSONError(ctx, http.StatusUnauthorized, "Invalid password", nil, nil)
			return
		} else {
			helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, err.Error())
			return
		}
	}

	helper.ReturnJSON(ctx, http.StatusOK, "Login success", responseService)
}
