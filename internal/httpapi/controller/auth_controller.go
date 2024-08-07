package controller

import (
	"log"
	"net/http"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/service"
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

	// get user agent
	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		log.Printf("User-Agent is required")
		helper.ReturnJSONError(ctx, http.StatusBadRequest, "User-Agent is required", nil, nil)
		return
	}

	// get client ip
	clientIp := ctx.ClientIP()
	if clientIp == "" {
		log.Printf("Client IP is required")
		helper.ReturnJSONError(ctx, http.StatusBadRequest, "Client IP is required", nil, nil)
		return
	}

	responseService, err := a.authService.Login(ctx, req.Username, req.Password, userAgent, clientIp)
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

func (a *AuthController) RefreshToken(ctx *gin.Context) {
	var req request.RefreshRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		massage, data := helper.GlobalCheckingErrorBindJson(err.Error(), req)
		log.Printf("Error Bind: %s", massage)
		helper.ReturnJSONError(ctx, http.StatusBadRequest, massage, nil, data)
		return
	}

	// get user agent
	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		log.Printf("User-Agent is required")
		helper.ReturnJSONError(ctx, http.StatusBadRequest, "User-Agent is required", nil, nil)
		return
	}

	// get client ip
	clientIp := ctx.ClientIP()
	if clientIp == "" {
		log.Printf("Client IP is required")
		helper.ReturnJSONError(ctx, http.StatusBadRequest, "Client IP is required", nil, nil)
		return
	}

	responseService, err := a.authService.RefreshToken(ctx, req.RefreshToken, userAgent, clientIp)

	if err != nil {
		log.Printf("Error: %s", err.Error())

		if err == service.ErrUserNotFound {
			helper.ReturnJSONError(ctx, http.StatusNotFound, "User not found", nil, nil)
			return
		} else if err == service.ErrorInvalidPassword {
			helper.ReturnJSONError(ctx, http.StatusUnauthorized, "Invalid password", nil, nil)
			return
		} else if err.Error() == "session is blocked" {
			helper.ReturnJSONError(ctx, http.StatusUnauthorized, "Session is blocked", nil, nil)
			return
		} else if err.Error() == "session is expired" {
			helper.ReturnJSONError(ctx, http.StatusUnauthorized, "session is expired", nil, nil)
			return
		} else {
			helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, err.Error())
			return
		}
	}

	helper.ReturnJSON(ctx, http.StatusOK, "Refresh token success", responseService)
}
