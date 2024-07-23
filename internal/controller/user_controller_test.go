package controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/fajaramaulana/simple_bank_project/db/mock"
	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/controller"
	"github.com/fajaramaulana/simple_bank_project/internal/service"
	"github.com/fajaramaulana/simple_bank_project/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserController_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mock data
	user := randomUser(t)
	password := util.RandomName()

	tests := []struct {
		name           string
		body           gin.H
		mockSetup      func(store *mockdb.MockStore)
		expectedStatus int
		expectedBody   string
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success create",
			body: gin.H{
				"email":     user.Email,
				"password":  password,
				"username":  user.Username,
				"full_name": user.FullName,
				"currency":  "USD",
			},
			mockSetup: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(db.GetUserByEmailRow{}, sql.ErrNoRows) // Simulate user does not exist

				store.EXPECT().GetUserByUsername(gomock.Any(), user.Username).Times(1).Return(db.GetUserByUsernameRow{}, sql.ErrNoRows)

				store.EXPECT().
					CreateUserWithAccountTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserWithAccountResult{}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "Account created",
		},
		{
			name: "Conflict - Email Already Exists",
			body: gin.H{
				"email":     user.Email,
				"password":  password,
				"username":  user.Username,
				"full_name": user.FullName,
				"currency":  "USD",
			},
			mockSetup: func(store *mockdb.MockStore) {
				// Simulate that the user already exists by returning a non-nil result for GetUserByEmail
				store.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(db.GetUserByEmailRow{
						Email: user.Email,
					}, nil) // Simulate user exists

				// Simulate that the username does not exist
				store.EXPECT().
					GetUserByUsername(gomock.Any(), user.Username).
					Times(0)

				// No call to CreateUserWithAccountTx should be made
				store.EXPECT().
					CreateUserWithAccountTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   "User Already Exist",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, w.Code)
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				meta, ok := responseBody["meta"].(map[string]interface{})
				require.True(t, ok)
				message, ok := meta["message"].(string)
				require.True(t, ok)
				require.Equal(t, "User Already Exist", message)
			},
		},
		{
			name: "Conflict - Username Already Exists",
			body: gin.H{
				"email":     user.Email,
				"password":  password,
				"username":  user.Username,
				"full_name": user.FullName,
				"currency":  "USD",
			},
			mockSetup: func(store *mockdb.MockStore) {
				// Simulate that the user already exists by returning a non-nil result for GetUserByEmail
				store.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(db.GetUserByEmailRow{}, sql.ErrNoRows)

				// Simulate that the username does not exist
				store.EXPECT().
					GetUserByUsername(gomock.Any(), user.Username).
					Times(1).Return(db.GetUserByUsernameRow{
					Email:    user.Email,
					Username: user.Username,
				}, nil)

				// No call to CreateUserWithAccountTx should be made
				store.EXPECT().
					CreateUserWithAccountTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   "User Already Exist",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, w.Code)
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				meta, ok := responseBody["meta"].(map[string]interface{})
				require.True(t, ok)
				message, ok := meta["message"].(string)
				require.True(t, ok)
				require.Equal(t, "User Already Exist", message)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"email":     user.Email,
				"password":  password,
				"username":  user.Username,
				"full_name": user.FullName,
				"currency":  "USD",
			},
			mockSetup: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(db.GetUserByEmailRow{}, sql.ErrNoRows) // Simulate user does not exist

				store.EXPECT().GetUserByUsername(gomock.Any(), user.Username).Times(1).Return(db.GetUserByUsernameRow{}, sql.ErrNoRows)

				store.EXPECT().
					CreateUserWithAccountTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserWithAccountResult{}, sql.ErrConnDone)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal server error",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				meta, ok := responseBody["meta"].(map[string]interface{})
				require.True(t, ok)
				message, ok := meta["message"].(string)
				require.True(t, ok)
				require.Equal(t, "Internal server error", message)
			},
		},
		{
			name: "password too short",
			body: gin.H{
				"email":     user.Email,
				"password":  "123",
				"username":  user.Username,
				"full_name": user.FullName,
				"currency":  "USD",
			},
			mockSetup: func(store *mockdb.MockStore) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Validation error",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "invalid-email",
			body: gin.H{
				"email":     "invalidEmail",
				"password":  user.Email,
				"username":  user.Username,
				"full_name": user.FullName,
				"currency":  "USD",
			},
			mockSetup: func(store *mockdb.MockStore) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Validation error",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "invalid-username",
			body: gin.H{
				"email":     user.Email,
				"password":  user.Email,
				"username":  "user",
				"full_name": user.FullName,
				"currency":  "USD",
			},
			mockSetup:      func(store *mockdb.MockStore) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Validation error",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tt.mockSetup(store)

			userService := service.NewUserService(store)
			userController := controller.NewUserController(userService)

			bodyJSON, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("Failed to marshal body: %v", err)
			}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(bodyJSON))

			userController.CreateUser(ctx)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !bytes.Contains(w.Body.Bytes(), []byte(tt.expectedBody)) {
				t.Errorf("Expected body to contain %q", tt.expectedBody)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func randomUser(t *testing.T) db.CreateUserParams {
	passHash, err := util.MakePasswordBcrypt(util.RandomName())

	require.NoError(t, err)
	user := db.CreateUserParams{
		Username:       util.RandomUsername(),
		Email:          util.RandomEmail(),
		FullName:       util.RandomName(),
		HashedPassword: passHash,
	}

	return user
}

func randomUser2(t *testing.T) (db.GetDetailLoginByUsernameRow, string) {
	password := "Password123!"
	passHash, err := util.MakePasswordBcrypt(password)

	require.NoError(t, err)
	user := db.GetDetailLoginByUsernameRow{
		Username:       util.RandomUsername(),
		Email:          util.RandomEmail(),
		FullName:       util.RandomName(),
		HashedPassword: passHash,
	}

	return user, password
}

func randomUser3() db.GetUserByUserUUIDRow {
	user := db.GetUserByUserUUIDRow{
		Username: util.RandomUsername(),
		Email:    util.RandomEmail(),
		FullName: util.RandomName(),
		UserUuid: util.RandomUUID(),
		Role:     util.RandomRole(),
	}

	return user
}
