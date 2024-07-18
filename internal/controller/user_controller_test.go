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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserController_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a new mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock store
	store := mockdb.NewMockStore(ctrl)

	// Create a new instance of UserService with the mock store
	userService := service.NewUserService(store)

	// Create a new instance of UserController with the UserService
	userController := controller.NewUserController(userService)

	// Define test user data
	email := "test@example.com"
	username := "testuser"
	fullName := "Test User"
	password := "password"
	currency := "USD"

	// Test cases
	tests := []struct {
		name           string
		body           gin.H
		mockSetup      func()
		expectedStatus int
		expectedBody   string
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Success - Create User",
			body: gin.H{
				"email":     email,
				"username":  username,
				"password":  password,
				"full_name": fullName,
				"currency":  currency,
			},
			mockSetup: func() {
				// Mock setup for CreateUser
				_ = db.CreateUserParams{
					Email:    email,
					Username: username,
					FullName: fullName,
				}
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetUserByEmailRow{}, sql.ErrNoRows) // Simulate user not found

				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Times(1).Return(db.GetUserByUsernameRow{}, sql.ErrNoRows)

				store.EXPECT().
					CreateUserWithAccountTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserWithAccountResult{}, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)

				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				message := responseBody["meta"].(map[string]interface{})["message"].(string)
				require.Equal(t, "Account created", message)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "Account created",
		},
		// {
		// 	name: "Conflict - User Already Exists",
		// 	body: gin.H{
		// 		"email":     email,
		// 		"username":  username,
		// 		"password":  password,
		// 		"full_name": fullName,
		// 		"currency":  currency,
		// 	},
		// 	mockSetup: func() {
		// 		// Mock setup for GetUserByEmail to simulate user already exists
		// 		store.EXPECT().
		// 			GetUserByEmail(gomock.Any(), email).
		// 			Times(1).
		// 			Return(db.GetUserByEmailRow{}, nil) // Simulate user already exists
		// 	},
		// 	checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusConflict, w.Code)
		// 		var responseBody map[string]interface{}
		// 		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		// 		require.NoError(t, err)

		// 		message := responseBody["meta"].(map[string]interface{})["message"].(string)

		// 		require.Equal(t, "User Already Exist", message)
		// 	},
		// 	expectedStatus: http.StatusConflict,
		// 	expectedBody:   "User Already Exists",
		// },
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build stubs
			tt.mockSetup()

			// Create a new recorder to simulate HTTP responses
			w := httptest.NewRecorder()

			// Create a new HTTP request with the specified body
			body, err := json.Marshal(tt.body)
			require.NoError(t, err)
			req, err := http.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
			require.NoError(t, err)

			// Create a new Gin context from the request and response recorder
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			// Call the CreateUser method of UserController
			userController.CreateUser(ctx)

			// Assert the HTTP status code
			// require.Equal(t, tt.expectedStatus, w.Code)

			tt.checkResponse(t, w)
		})
	}
}
