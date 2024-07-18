package controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/fajaramaulana/simple_bank_project/db/mock"
	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/service"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMockGetAccountByUUIDController(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	account := randomAccount()

	testCases := []struct {
		name          string
		AccountUuid   uuid.UUID
		paramUuid     string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			AccountUuid: account.AccountUuid,
			paramUuid:   account.AccountUuid.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByUUID(gomock.Any(), gomock.Eq(account.AccountUuid)).Times(1).Return(account, nil)

				store.EXPECT().GetUserByUserUUID(gomock.Any(), gomock.Eq(account.UserUuid)).Times(1).Return(db.GetUserByUserUUIDRow{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:        "NotFound",
			AccountUuid: account.AccountUuid,
			paramUuid:   account.AccountUuid.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByUUID(gomock.Any(), gomock.Eq(account.AccountUuid)).Times(1).Return(db.GetAccountByUUIDRow{}, sql.ErrNoRows)
				store.EXPECT().GetUserByUserUUID(gomock.Any(), gomock.Any()).Times(1).Return(db.GetUserByUserUUIDRow{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "InternalError",
			AccountUuid: account.AccountUuid,
			paramUuid:   account.AccountUuid.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByUUID(gomock.Any(), gomock.Eq(account.AccountUuid)).Times(1).Return(db.GetAccountByUUIDRow{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "BadRequest",
			AccountUuid: uuid.Nil,
			paramUuid:   "invalid-uuid",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByUUID(gomock.Any(), uuid.Nil).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			// Create a mock Gin engine
			router := gin.New()

			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			accountService := service.NewAccountService(store)
			accountController := controller.NewAccountController(accountService)

			// Set up the Gin route and handler
			router.GET("/api/v1/account/:uuid", func(c *gin.Context) {
				accountController.GetAccount(c)
			})

			url := fmt.Sprintf("/api/v1/account/%v", tc.paramUuid)

			request := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAccountsController(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Mock data for accounts
	account := randomAccount()
	account2 := randomAccount()
	account3 := randomAccount()
	account4 := randomAccount()
	account5 := randomAccount()

	expectedAccounts := []db.ListAccountsRow{
		{ID: account.ID, AccountUuid: account.AccountUuid, Owner: account.Owner, UserUuid: account.UserUuid, Currency: account.Currency, Balance: account.Balance, CreatedAt: account.CreatedAt, Status: account.Status},
		{ID: account2.ID, AccountUuid: account2.AccountUuid, Owner: account2.Owner, UserUuid: account2.UserUuid, Currency: account2.Currency, Balance: account2.Balance, CreatedAt: account2.CreatedAt, Status: account2.Status},
		{ID: account3.ID, AccountUuid: account3.AccountUuid, Owner: account3.Owner, UserUuid: account3.UserUuid, Currency: account3.Currency, Balance: account3.Balance, CreatedAt: account3.CreatedAt, Status: account3.Status},
		{ID: account4.ID, AccountUuid: account4.AccountUuid, Owner: account4.Owner, UserUuid: account4.UserUuid, Currency: account4.Currency, Balance: account4.Balance, CreatedAt: account4.CreatedAt, Status: account4.Status},
		{ID: account5.ID, AccountUuid: account5.AccountUuid, Owner: account5.Owner, UserUuid: account5.UserUuid, Currency: account5.Currency, Balance: account5.Balance, CreatedAt: account5.CreatedAt, Status: account5.Status},
	}

	testCases := []struct {
		name          string
		page          int
		limit         int
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			page:  1,
			limit: 5,
			buildStubs: func(store *mockdb.MockStore) {
				mockParams := db.ListAccountsParams{Limit: 5, Offset: 0}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(mockParams)).
					Times(1).
					Return(expectedAccounts, nil)

				store.EXPECT().
					CountAccounts(gomock.Any()).
					Times(1).
					Return(int64(len(expectedAccounts)), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var responseBody map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				require.Equal(t, "Account found", responseBody["meta"].(map[string]interface{})["message"])
				require.Equal(t, len(expectedAccounts), len(responseBody["data"].([]interface{})))
			},
		},
		{
			name:  "Bad Request",
			page:  0,
			limit: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0).
					Return([]db.ListAccountsRow{}, nil)

				store.EXPECT().
					CountAccounts(gomock.Any()).
					Times(0).
					Return(int64(0), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Not Found",
			page:  1,
			limit: 5,
			buildStubs: func(store *mockdb.MockStore) {
				mockParams := db.ListAccountsParams{Limit: 5, Offset: 0}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(mockParams)).
					Times(1).
					Return([]db.ListAccountsRow{}, nil)

				store.EXPECT().
					CountAccounts(gomock.Any()).
					Times(1).
					Return(int64(0), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "Internal Server Error",
			page:  1,
			limit: 5,
			buildStubs: func(store *mockdb.MockStore) {
				mockParams := db.ListAccountsParams{Limit: 5, Offset: 0}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(mockParams)).
					Times(1).
					Return(nil, sql.ErrConnDone) // Replace someError with your mock error

				store.EXPECT().
					CountAccounts(gomock.Any()).
					Times(0).
					Return(int64(0), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

				var responseBody map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				require.Equal(t, "Internal server error", responseBody["meta"].(map[string]interface{})["message"])
				require.Nil(t, responseBody["data"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock Gin engine
			router := gin.New()

			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock the store and service
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			accountService := service.NewAccountService(store)
			accountController := controller.NewAccountController(accountService)

			// Set up the Gin route and handler
			router.GET("/api/v1/accounts", func(c *gin.Context) {
				accountController.GetAccounts(c)
			})

			// Create a mock HTTP request
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/accounts?limit=%d&page=%d", tc.limit, tc.page), nil)
			rec := httptest.NewRecorder()

			// Perform the request against the mock router
			router.ServeHTTP(rec, req)

			// Check the response
			tc.checkResponse(t, rec)
		})
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.GetAccountByUUIDRow) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount map[string]interface{}
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	accountUuid := gotAccount["data"].(map[string]interface{})["account_uuid"]
	owner := gotAccount["data"].(map[string]interface{})["owner"]
	balance := gotAccount["data"].(map[string]interface{})["balance"]
	currency := gotAccount["data"].(map[string]interface{})["currency"]

	require.Equal(t, account.AccountUuid.String(), accountUuid, "AccountUuid mismatch")
	require.Equal(t, account.Owner, owner, "Owner mismatch")
	require.Equal(t, account.Balance, balance, "Balance mismatch")
	require.Equal(t, account.Currency, currency, "Currency mismatch")
}

func randomAccount() db.GetAccountByUUIDRow {
	randomInt, err := util.RandomInt(1000)
	if err != nil {
		return db.GetAccountByUUIDRow{}
	}

	user := db.GetUserByUserUUIDRow{
		UserUuid: uuid.New(),
	}

	r := util.NewRandomMoneyGenerator()
	return db.GetAccountByUUIDRow{
		ID:          int64(randomInt[0]),
		Owner:       util.RandomName(),
		Balance:     util.RandomMoney(r, 10.00, 99999999.00),
		Currency:    util.RandomCurrency(),
		AccountUuid: uuid.New(),
		UserUuid:    user.UserUuid,
	}
}
