package router

import (
	"github.com/fajaramaulana/simple_bank_project/internal/controller"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router  *gin.Engine
	account *controller.AccountController
}

func NewRouter(account *controller.AccountController) *Router {
	return &Router{
		Router:  gin.Default(),
		account: account,
	}
}

// SetupRouter sets up the router for the application and starts the server on the specified port.
func (r *Router) SetupRouter(port string) {
	v1 := r.Router.Group("/api/v1")

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

	r.Router.Run(":" + port)
}
