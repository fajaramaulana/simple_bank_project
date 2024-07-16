package router

import (
	"github.com/fajaramaulana/simple_bank_project/internal/controller"
	"github.com/gin-gonic/gin"
)

type router struct {
	router  *gin.Engine
	account *controller.AccountController
}

func NewRouter(account *controller.AccountController) *router {
	return &router{
		router:  gin.Default(),
		account: account,
	}
}

func (r *router) SetupRouter(port string) {
	v1 := r.router.Group("/api/v1")

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

	r.router.Run(":" + port)
}
