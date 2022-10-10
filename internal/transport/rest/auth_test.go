package rest

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/Arkosh744/simpleREST_blog/internal/service"
	"github.com/Arkosh744/simpleREST_blog/internal/transport/rest/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"golang.org/x/net/context"
	"net/http/httptest"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mocks.MockUsers, ctx context.Context, user domain.SignUpInput)
	tests := []struct {
		name                 string
		inputBody            string
		inputUser            domain.SignUpInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name": "TestName", "email": "username@gmail.com", "password": "qwerty"}`,
			inputUser: domain.SignUpInput{
				Name:     "TestName",
				Email:    "username@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(r *mocks.MockUsers, ctx context.Context, inp domain.SignUpInput) {
				r.EXPECT().SignUp(gomock.Any(), inp).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `"OK"`,
		},
		{
			name:      "Wrong Input",
			inputBody: `{"name": "username"}`,
			inputUser: domain.SignUpInput{},
			mockBehavior: func(r *mocks.MockUsers, ctx context.Context, inp domain.SignUpInput) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"name": "username", "email": "username@gmail.com", "password": "qwerty"}`,
			inputUser: domain.SignUpInput{
				Name:     "username",
				Email:    "username@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(r *mocks.MockUsers, ctx context.Context, inp domain.SignUpInput) {
				r.EXPECT().SignUp(gomock.Any(), inp).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("---------------- Start %s ----------------\n", test.name)
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mocks.NewMockUsers(c)
			test.mockBehavior(auth, context.Background(), test.inputUser)
			fmt.Printf("test.inputUser: %v\n", test.inputUser)
			handler := NewHandler(&service.Posts{}, auth)
			// Init Endpoint
			r := gin.Default()
			r.POST("/auth/sign-up", handler.signUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-up",
				bytes.NewBufferString(test.inputBody))
			// Make Request
			r.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mocks.MockUsers, ctx context.Context, user domain.SignInInput)
	//var getContext = func(id int) *gin.Context {
	//	ctx := &gin.Context{}
	//	ctx.Set("id", id)
	//	return ctx
	//}

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            domain.SignInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"email": "username@gmail.com", "password": "qwerty"}`,
			inputUser: domain.SignInInput{
				Email:    "username@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(r *mocks.MockUsers, ctx context.Context, inp domain.SignInInput) {
				r.EXPECT().SignIn(gomock.Any(), inp).Return("mocked_token", "mocked_refresh_token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"mocked_token"}`,
		},
		{
			name:      "Invalid Credentials",
			inputBody: `{"email": "username@gmail.com", "password": "qwertyloggg"}`,
			inputUser: domain.SignInInput{
				Email:    "username@gmail.com",
				Password: "qwertyloggg",
			},
			mockBehavior: func(r *mocks.MockUsers, ctx context.Context, inp domain.SignInInput) {
				r.EXPECT().SignIn(gomock.Any(), inp).Return("", "", errors.New("invalid credentials"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"email": "username@gmail.com", "password": "qwerty"}`,
			inputUser: domain.SignInInput{
				Email:    "username@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(r *mocks.MockUsers, ctx context.Context, inp domain.SignInInput) {
				r.EXPECT().SignIn(gomock.Any(), inp).Return("", "", errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("---------------- Start %s ----------------\n", test.name)
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mocks.NewMockUsers(c)
			test.mockBehavior(auth, context.Background(), test.inputUser)
			fmt.Printf("test.inputUser: %v\n", test.inputUser)
			handler := NewHandler(&service.Posts{}, auth)
			// Init Endpoint
			r := gin.Default()
			r.POST("/auth/sign-in", handler.signIn)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-in",
				bytes.NewBufferString(test.inputBody))
			// Make Request
			r.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
