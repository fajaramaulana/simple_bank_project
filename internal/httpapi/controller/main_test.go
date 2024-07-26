package controller_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/router"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/service"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/setup"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func NewTestServer(t *testing.T, store db.Store) *router.Router {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	configEnv := setup.CheckingEnv(config)

	// configToken
	configToken := map[string]string{
		"token_secret":           configEnv.TokenSymmetricKey,
		"access_token_duration":  configEnv.AccessTokenDuration.String(),
		"refresh_token_duration": configEnv.RefreshTokenDuration.String(),
	}
	fmt.Printf("%# v\n", configToken)
	// account
	accountService := service.NewAccountService(store)
	accountController := controller.NewAccountController(accountService)

	// transfer
	transferService := service.NewTransactionService(store)
	transferController := controller.NewTransactionController(transferService)

	// user
	userService := service.NewUserService(store)
	userController := controller.NewUserController(userService)

	// auth
	authService := service.NewAuthService(store, configToken)
	authController := controller.NewAuthController(authService)

	server, err := router.NewRouter(accountController, transferController, userController, authController, configToken)
	require.NoError(t, err)

	return server
}
