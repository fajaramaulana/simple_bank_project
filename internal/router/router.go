package router

import (
	"fmt"

	"github.com/fajaramaulana/simple_bank_project/internal/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Router struct {
	Engine      *gin.Engine
	account     *controller.AccountController
	transaction *controller.TransactionController
	user        *controller.UserController
	tokenMaker  token.Maker
}

func NewRouter(account *controller.AccountController, transaction *controller.TransactionController, user *controller.UserController, configToken map[string]string) (*Router, error) {
	tokenMaker, err := token.NewPasetoMaker(configToken["token_secret"])
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	return &Router{
		Engine:      gin.Default(),
		account:     account,
		transaction: transaction,
		user:        user,
		tokenMaker:  tokenMaker,
	}, nil
}

// SetupRouter sets up the router for the application.
func (r *Router) SetupRouter() {
	v1 := r.Engine.Group("/api/v1")

	// Register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", helper.CurrencyValidator)
	}

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	// account
	v1.POST("/account", r.account.CreateAccount)
	v1.GET("/account/:uuid", r.account.GetAccount)
	v1.GET("/accounts", r.account.GetAccounts)
	v1.PUT("/account/:uuid", r.account.UpdateAccount)

	// transaction
	v1.POST("/transaction", r.transaction.CreateTransfer)

	// user
	v1.POST("/user", r.user.CreateUser)
}

// StartServer starts the HTTP server on the specified port.
func (r *Router) StartServer(port string) {
	r.Engine.Run(":" + port)
}
