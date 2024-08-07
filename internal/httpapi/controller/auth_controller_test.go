package controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/fajaramaulana/simple_bank_project/db/mock"
	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/service"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user, password := randomUser2(t)

	testCase := []struct {
		name           string
		body           gin.H
		userAgent      string
		clientIp       string
		mockSetup      func(store *mockdb.MockStore)
		expectedStatus int
		expectedBody   string
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			userAgent: "test",
			clientIp:  "1.1.1.1",
			mockSetup: func(store *mockdb.MockStore) {
				store.EXPECT().GetDetailLoginByUsername(gomock.Any(), user.Username).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1).Return(db.Session{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Login success",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				var responseBody map[string]interface{}

				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				meta, ok := responseBody["meta"].(map[string]interface{})
				require.True(t, ok)
				message, ok := meta["message"].(string)
				require.True(t, ok)
				require.Equal(t, "Login success", message)
			},
		},
		{
			name: "Invalid Password",
			body: gin.H{
				"username": user.Username,
				"password": "invalid password",
			},
			userAgent: "test",
			clientIp:  "1.1.1.1",
			mockSetup: func(store *mockdb.MockStore) {
				store.EXPECT().GetDetailLoginByUsername(gomock.Any(), user.Username).Times(1).Return(user, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid password",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
				var responseBody map[string]interface{}

				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				meta, ok := responseBody["meta"].(map[string]interface{})
				require.True(t, ok)
				message, ok := meta["message"].(string)
				require.True(t, ok)
				require.Equal(t, "Invalid password", message)
			},
		},
		{
			name: "Password Too Short",
			body: gin.H{
				"username": user.Username,
				"password": "short",
			},
			userAgent: "test",
			clientIp:  "1.1.1.1",
			mockSetup: func(store *mockdb.MockStore) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Validation error",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Username Too Short",
			body: gin.H{
				"username": "short",
				"password": "short",
			},
			userAgent: "test",
			clientIp:  "1.1.1.1",
			mockSetup: func(store *mockdb.MockStore) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Validation error",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Empty user agent",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			userAgent: "",
			clientIp:  "1.1.1.1",
			mockSetup: func(store *mockdb.MockStore) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Validation error",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "User Not Found",
			body: gin.H{
				"username": "notfound",
				"password": "password",
			},
			userAgent:      "test",
			clientIp:       "1.1.1.1",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "User not found",
			mockSetup: func(store *mockdb.MockStore) {
				store.EXPECT().GetDetailLoginByUsername(gomock.Any(), "notfound").Times(1).Return(db.GetDetailLoginByUsernameRow{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
				var responseBody map[string]interface{}

				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				meta, ok := responseBody["meta"].(map[string]interface{})
				require.True(t, ok)
				message, ok := meta["message"].(string)
				require.True(t, ok)
				require.Equal(t, "User not found", message)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			userAgent:      "test",
			clientIp:       "1.1.1.1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal server error",
			mockSetup: func(store *mockdb.MockStore) {
				store.EXPECT().GetDetailLoginByUsername(gomock.Any(), gomock.Any()).Times(1).Return(db.GetDetailLoginByUsernameRow{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.mockSetup(store)

			// config token

			configToken := map[string]string{
				"token_secret":           util.RandomString(32),
				"access_token_duration":  time.Minute.String(),
				"refresh_token_duration": (15 * time.Minute).String(),
			}

			service := service.NewAuthService(store, configToken)
			controller := controller.NewAuthController(service)

			bodyJSON, err := json.Marshal(tt.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(bodyJSON))

			ctx.Request.Header.Set("User-Agent", tt.userAgent)

			controller.Login(ctx)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
