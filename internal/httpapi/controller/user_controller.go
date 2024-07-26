package controller

import (
	"log"
	"net/http"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/request"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/service"
	"github.com/gin-gonic/gin"
)

// UserController handles HTTP requests related to accounts.
type UserController struct {
	userService *service.UserService
}

// NewUserController creates a new instance of the UserController struct.
// It takes an userService parameter of type *service.UserService and returns a pointer to the UserController.
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) CreateUser(ctx *gin.Context) {
	var request request.CreateUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		message, data := helper.GlobalCheckingErrorBindJson(err.Error(), request)
		log.Printf("Error Bind: %s", message)
		helper.ReturnJSONError(ctx, http.StatusBadRequest, message, nil, data)
		return
	}

	res := helper.DoValidation(&request)

	if len(res) > 0 {
		log.Println("Error: Validation error")
		helper.ReturnJSONError(ctx, http.StatusBadRequest, "Validation error", nil, res)
		return
	}
	// check request data

	user, err := u.userService.CreateUser(ctx.Request.Context(), &request)
	if err != nil {
		if user.Email != "" {
			log.Println("Error: User already exists")
			helper.ReturnJSONError(ctx, http.StatusConflict, "User Already Exist", nil, nil)
			return
		} else {
			log.Printf("Error: %s", err.Error())
			helper.ReturnJSONError(ctx, http.StatusInternalServerError, "Internal server error", nil, err.Error())
			return
		}
	}

	// checking account already exist
	helper.ReturnJSON(ctx, http.StatusCreated, "Account created", user)
}
