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
	"github.com/fajaramaulana/simple_bank_project/internal/router"
	"github.com/fajaramaulana/simple_bank_project/internal/service"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMockGetAccountByUUIDController(t *testing.T) {
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
			paramUuid:   " ",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByUUID(gomock.Any(), uuid.Nil).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		fmt.Printf("%# v\n", i)
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			// build stubs

			// start test server and send request
			accountService := service.NewAccountService(store)
			accountController := controller.NewAccountController(accountService)
			r := router.NewRouter(accountController)
			r.SetupRouter()

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/account/%v", tc.paramUuid)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Serve the HTTP request using the router
			r.Engine.ServeHTTP(recorder, request)
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

	// Manually compare each field if necessary
	require.Equal(t, account.AccountUuid.String(), accountUuid, "AccountUuid mismatch")
	require.Equal(t, account.Owner, owner, "Owner mismatch")
	require.Equal(t, account.Balance, balance, "Balance mismatch")
	require.Equal(t, account.Currency, currency, "Currency mismatch")
	// require.WithinDuration(t, account.CreatedAt, gotAccount.CreatedAt, time.Second, "CreatedAt mismatch")
}

func randomAccount() db.GetAccountByUUIDRow {

	randomInt, err := util.RandomInt(1000)
	if err != nil {
		return db.GetAccountByUUIDRow{}
	}
	r := util.NewRandomMoneyGenerator()
	return db.GetAccountByUUIDRow{
		ID:           int64(randomInt[0]),
		Owner:        util.RandomName(),
		Balance:      util.RandomMoney(r, 10.00, 99999999.00),
		Currency:     util.RandomCurrency(),
		Email:        util.RandomEmail(),
		AccountUuid:  uuid.New(),
		RefreshToken: "refresh_token1",
	}
}
