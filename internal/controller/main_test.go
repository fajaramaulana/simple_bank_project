package controller_test

import (
	"os"
	"testing"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/router"
	"github.com/fajaramaulana/simple_bank_project/internal/service"
	"github.com/fajaramaulana/simple_bank_project/internal/setup"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func NewTestServer(t *testing.T, store db.Store) *router.Router {
	// config := util.Config{
	// 	TokenSymmetricKey:   util.RandomString(32),
	// 	AccessTokenDuration: time.Minute,
	// }

	// server, err := routerNewServer(config, store)
	setup.CheckingEnv()

	// configToken
	configToken := map[string]string{
		"token_secret":          os.Getenv("TOKEN_SYMMETRIC_KEY"),
		"access_token_duration": os.Getenv("ACCESS_TOKEN_DURATION"),
	}
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
