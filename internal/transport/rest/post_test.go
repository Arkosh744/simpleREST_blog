package rest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/Arkosh744/simpleREST_blog/internal/transport/rest/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Create(t *testing.T) {
	type mockBehavior func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.Post)
	var AuthorId int64 = 1

	tests := []struct {
		name                 string
		inputBody            string
		inputPost            domain.Post
		mockCoockie          *http.Cookie
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"title": "TestTitle", "body": "TestBody"}`,
			inputPost: domain.Post{
				Title:    "TestTitle",
				Body:     "TestBody",
				AuthorId: AuthorId,
			},
			mockCoockie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.Post) {
				mockUser.EXPECT().GetIdByToken(gomock.Any(), "refreshToken").Return(AuthorId, nil)
				mockPost.EXPECT().Create(gomock.Any(), inp).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"body":"TestBody","title":"TestTitle"}`,
		},
		{
			name:        "W/o Cookie",
			inputBody:   `{"title": "TestTitle", "body": "TestBody"}`,
			mockCoockie: &http.Cookie{},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.Post) {
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"http: named cookie not present"}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"name": "username"}`,
			inputPost: domain.Post{
				AuthorId: AuthorId,
			},
			mockCoockie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.Post) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input post body"}`,
		},
		{
			name:      "Ok",
			inputBody: `{"title": "TestTitle", "body": "TestBody"}`,
			inputPost: domain.Post{
				Title:    "TestTitle",
				Body:     "TestBody",
				AuthorId: AuthorId,
			},
			mockCoockie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.Post) {
				mockUser.EXPECT().GetIdByToken(gomock.Any(), "refreshToken").Return(AuthorId, nil)
				mockPost.EXPECT().Create(gomock.Any(), inp).Return(errors.New("something went wrong"))
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

			post := mocks.NewMockPosts(c)
			auth := mocks.NewMockUsers(c)
			test.mockBehavior(post, auth, context.Background(), test.inputPost)
			fmt.Println(test.inputPost)
			fmt.Printf("test.inputPost: %v\n", test.inputPost)
			handler := NewHandler(post, auth)
			// Init Endpoint
			r := gin.Default()
			r.POST("/post/", handler.Create)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/post/",
				bytes.NewBufferString(test.inputBody))
			mockCookie := test.mockCoockie
			req.AddCookie(mockCookie)
			// Make Request
			r.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
