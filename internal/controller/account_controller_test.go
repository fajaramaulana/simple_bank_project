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
	"time"

	mockdb "github.com/fajaramaulana/simple_bank_project/db/mock"
	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/middleware"
	"github.com/fajaramaulana/simple_bank_project/internal/handler/token"
	"github.com/fajaramaulana/simple_bank_project/internal/setup"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMockGetAccountByUUIDController(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	user := randomUser3()
	account := randomAccount(t, user.UserUuid)

	testCases := []struct {
		name          string
		AccountUuid   uuid.UUID
		paramUuid     string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			AccountUuid: account.AccountUuid,
			paramUuid:   account.AccountUuid.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, user.Role)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByUUID(gomock.Any(), account.AccountUuid).Times(1).Return(account, nil)
				store.EXPECT().GetUserByUserUUID(gomock.Any(), user.UserUuid).Times(1).Return(user, nil)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, user.Role)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByUUID(gomock.Any(), gomock.Eq(account.AccountUuid)).Times(1).Return(db.GetAccountByUUIDRow{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "InternalError",
			AccountUuid: account.AccountUuid,
			paramUuid:   account.AccountUuid.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, user.Role)
			},
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, user.Role)
			},
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
			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// Initialize the server with the mock store
			server := setup.InitializeAndStartAppTest(t, store)

			// Define the URL path for the request
			urlPath := fmt.Sprintf("/api/v1/account/%s", tc.paramUuid)

			// Initialize the recorder and request
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, urlPath, nil)
			require.NoError(t, err)

			// Set up authentication for the request
			tc.setupAuth(t, request, server.TokenMaker)

			// Serve the request
			server.Engine.ServeHTTP(recorder, request)

			// Check the response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAccountsController(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	user := randomUser3()

	// Mock data for accounts
	account := randomAccount(t, user.UserUuid)
	account2 := randomAccount(t, user.UserUuid)
	account3 := randomAccount(t, user.UserUuid)
	account4 := randomAccount(t, user.UserUuid)
	account5 := randomAccount(t, user.UserUuid)

	expectedAccountsAdmin := []db.ListAccountsRow{
		{ID: account.ID, AccountUuid: account.AccountUuid, Owner: account.Owner, UserUuid: account.UserUuid, Currency: account.Currency, Balance: account.Balance, CreatedAt: account.CreatedAt, Status: account.Status},
		{ID: account2.ID, AccountUuid: account2.AccountUuid, Owner: account2.Owner, UserUuid: account2.UserUuid, Currency: account2.Currency, Balance: account2.Balance, CreatedAt: account2.CreatedAt, Status: account2.Status},
		{ID: account3.ID, AccountUuid: account3.AccountUuid, Owner: account3.Owner, UserUuid: account3.UserUuid, Currency: account3.Currency, Balance: account3.Balance, CreatedAt: account3.CreatedAt, Status: account3.Status},
		{ID: account4.ID, AccountUuid: account4.AccountUuid, Owner: account4.Owner, UserUuid: account4.UserUuid, Currency: account4.Currency, Balance: account4.Balance, CreatedAt: account4.CreatedAt, Status: account4.Status},
		{ID: account5.ID, AccountUuid: account5.AccountUuid, Owner: account5.Owner, UserUuid: account5.UserUuid, Currency: account5.Currency, Balance: account5.Balance, CreatedAt: account5.CreatedAt, Status: account5.Status},
	}

	expectedAccountsCustomer := []db.ListAccountsByUserUUIDRow{
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
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK-admin",
			page:  1,
			limit: 5,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, "admin")
			},
			buildStubs: func(store *mockdb.MockStore) {
				mockParams := db.ListAccountsParams{Limit: 5, Offset: 0}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(mockParams)).
					Times(1).
					Return(expectedAccountsAdmin, nil)

				store.EXPECT().
					CountAccounts(gomock.Any()).
					Times(1).
					Return(int64(len(expectedAccountsAdmin)), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var responseBody map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				require.Equal(t, "Account found", responseBody["meta"].(map[string]interface{})["message"])
				require.Equal(t, len(expectedAccountsAdmin), len(responseBody["data"].([]interface{})))
			},
		},
		{
			name:  "Bad Request-admin",
			page:  0,
			limit: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, "admin")
			},
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
			name:  "Not Found-admin",
			page:  1,
			limit: 5,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, "admin")
			},
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
			name:  "Internal Server Error-admin",
			page:  1,
			limit: 5,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, "admin")
			},
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
		{
			name:  "OK-customer",
			page:  1,
			limit: 5,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, "customer")
			},
			buildStubs: func(store *mockdb.MockStore) {
				mockParams := db.ListAccountsByUserUUIDParams{UserUuid: user.UserUuid, Limit: 5, Offset: 0}
				store.EXPECT().
					ListAccountsByUserUUID(gomock.Any(), gomock.Eq(mockParams)).
					Times(1).
					Return(expectedAccountsCustomer, nil)

				store.EXPECT().
					CountAccountsByUserUUID(gomock.Any(), user.UserUuid).
					Times(1).
					Return(int64(len(expectedAccountsCustomer)), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var responseBody map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				require.Equal(t, "Account found", responseBody["meta"].(map[string]interface{})["message"])
				require.Equal(t, len(expectedAccountsCustomer), len(responseBody["data"].([]interface{})))
			},
		},
		{
			name:  "Bad Request-customer",
			page:  0,
			limit: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, "customer")
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccountsByUserUUID(gomock.Any(), gomock.Any()).
					Times(0).
					Return([]db.ListAccountsByUserUUIDRow{}, nil)

				store.EXPECT().
					CountAccountsByUserUUID(gomock.Any(), user.UserUuid).
					Times(0).
					Return(int64(0), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Not Found-customer",
			page:  1,
			limit: 5,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, "customer")
			},
			buildStubs: func(store *mockdb.MockStore) {
				mockParams := db.ListAccountsByUserUUIDParams{UserUuid: user.UserUuid, Limit: 5, Offset: 0}
				store.EXPECT().
					ListAccountsByUserUUID(gomock.Any(), gomock.Eq(mockParams)).
					Times(1).
					Return([]db.ListAccountsByUserUUIDRow{}, nil)

				store.EXPECT().
					CountAccountsByUserUUID(gomock.Any(), user.UserUuid).
					Times(1).
					Return(int64(0), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "Internal Server Error-customer",
			page:  1,
			limit: 5,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				middleware.AddAuthorizationTestAPI(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.UserUuid.String(), time.Minute, "customer")
			},
			buildStubs: func(store *mockdb.MockStore) {
				mockParams := db.ListAccountsByUserUUIDParams{UserUuid: user.UserUuid, Limit: 5, Offset: 0}
				store.EXPECT().
					ListAccountsByUserUUID(gomock.Any(), gomock.Eq(mockParams)).
					Times(1).
					Return(nil, sql.ErrConnDone) // Replace someError with your mock error

				store.EXPECT().
					CountAccountsByUserUUID(gomock.Any(), user.UserUuid).
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

			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// Initialize the server with the mock store
			server := setup.InitializeAndStartAppTest(t, store)

			// Define the URL path for the request
			urlPath := fmt.Sprintf("/api/v1/accounts?limit=%d&page=%d", tc.limit, tc.page)

			// Initialize the recorder and request
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, urlPath, nil)
			require.NoError(t, err)

			// Set up authentication for the request
			tc.setupAuth(t, request, server.TokenMaker)

			// Serve the request
			server.Engine.ServeHTTP(recorder, request)

			// Check the response
			tc.checkResponse(t, recorder)
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

func randomAccount(t *testing.T, Useruuid uuid.UUID) db.GetAccountByUUIDRow {
	randomInt, err := util.RandomInt(1000)
	require.NoError(t, err)

	r := util.NewRandomMoneyGenerator()
	return db.GetAccountByUUIDRow{
		ID:          int64(randomInt[0]),
		Owner:       util.RandomName(),
		Balance:     util.RandomMoney(r, 10.00, 99999999.00),
		Currency:    util.RandomCurrency(),
		AccountUuid: uuid.New(),
		UserUuid:    Useruuid,
	}
}
