package router

import (
	"fmt"

	"github.com/fajaramaulana/simple_bank_project/internal/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/helper"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/middleware"
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
	auth        *controller.AuthController
	TokenMaker  token.Maker
}

// NewRouter creates a new instance of the Router struct and initializes its dependencies.
// It takes an AccountController, TransactionController, UserController, AuthController,
// and a map of configuration tokens as parameters. It returns a pointer to the Router
// struct and an error if there was an issue creating the token maker.
func NewRouter(account *controller.AccountController, transaction *controller.TransactionController, user *controller.UserController, auth *controller.AuthController, configToken map[string]string) (*Router, error) {
	tokenMaker, err := token.NewPasetoMaker(configToken["token_secret"])
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	router := &Router{
		Engine:      gin.Default(),
		account:     account,
		transaction: transaction,
		user:        user,
		auth:        auth,
		TokenMaker:  tokenMaker,
	}

	// Register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", helper.CurrencyValidator)
	}

	router.SetupRouter()
	return router, nil
}

// SetupRouter sets up the router for the application.
func (r *Router) SetupRouter() {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	// auth
	v1.POST("/auth/login", r.auth.Login)

	authRoutesV1 := router.Group("/api/v1").Use(middleware.AuthMiddleware(r.TokenMaker))

	// account
	authRoutesV1.POST("/account", r.account.CreateAccount)
	authRoutesV1.GET("/account/:uuid", r.account.GetAccount)
	authRoutesV1.GET("/accounts", r.account.GetAccounts)
	authRoutesV1.PUT("/account/:uuid", r.account.UpdateAccount)

	// transaction
	authRoutesV1.POST("/transaction", r.transaction.CreateTransfer)

	// user
	authRoutesV1.POST("/user", r.user.CreateUser)

	r.Engine = router
}

// StartServer starts the HTTP server on the specified port.
func (r *Router) StartServer(port string) {
	r.Engine.Run(":" + port)
}
