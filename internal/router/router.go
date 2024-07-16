package router

import (
	"github.com/fajaramaulana/simple_bank_project/internal/controller"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine  *gin.Engine
	account *controller.AccountController
}

func NewRouter(account *controller.AccountController) *Router {
	return &Router{
		Engine:  gin.Default(),
		account: account,
	}
}

// SetupRouter sets up the router for the application.
func (r *Router) SetupRouter() {
	v1 := r.Engine.Group("/api/v1")

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
}

// StartServer starts the HTTP server on the specified port.
func (r *Router) StartServer(port string) {
	r.Engine.Run(":" + port)
}
