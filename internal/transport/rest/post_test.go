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
		mockCookie           *http.Cookie
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
			mockCookie: &http.Cookie{
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
			name:       "W/o Cookie",
			inputBody:  `{"title": "TestTitle", "body": "TestBody"}`,
			mockCookie: &http.Cookie{},
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
			mockCookie: &http.Cookie{
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
			name:      "Service Error",
			inputBody: `{"title": "TestTitle", "body": "TestBody"}`,
			inputPost: domain.Post{
				Title:    "TestTitle",
				Body:     "TestBody",
				AuthorId: AuthorId,
			},
			mockCookie: &http.Cookie{
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
			fmt.Printf("test.inputPost: %v\n", test.inputPost)
			handler := NewHandler(post, auth)
			// Init Endpoint
			r := gin.Default()
			r.POST("/post/", handler.Create)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/post/",
				bytes.NewBufferString(test.inputBody))
			mockCookie := test.mockCookie
			req.AddCookie(mockCookie)
			// Make Request
			r.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_List(t *testing.T) {
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

func TestHandler_GetById(t *testing.T) {
	type mockBehavior func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers,
		ctx context.Context, id int64, responsePost domain.Post)
	var AuthorId int64 = 1
	tests := []struct {
		name                 string
		inputID              int64
		mockCookie           *http.Cookie
		responsePost         domain.Post
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Ok",
			inputID: 1,
			mockCookie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			responsePost: domain.Post{
				Id:       1,
				Title:    "TestTitle",
				Body:     "TestBody",
				AuthorId: 1,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inputID int64, responsePost domain.Post) {
				mockUser.EXPECT().GetIdByToken(gomock.Any(), "refreshToken").Return(AuthorId, nil)
				mockPost.EXPECT().GetById(gomock.Any(), inputID, AuthorId).Return(responsePost, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"title":"TestTitle","body":"TestBody","AuthorId":1,"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:       "W/o Cookie",
			inputID:    1,
			mockCookie: &http.Cookie{},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inputID int64, responsePost domain.Post) {
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"http: named cookie not present"}`,
		},
		{
			name:    "Wrong input",
			inputID: 1,
			mockCookie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inputID int64, responsePost domain.Post) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input post id"}`,
		},
		{
			name:    "Service Error",
			inputID: 1,
			mockCookie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inputID int64, responsePost domain.Post) {
				mockUser.EXPECT().GetIdByToken(gomock.Any(), "refreshToken").Return(AuthorId, nil)
				mockPost.EXPECT().GetById(gomock.Any(), inputID, AuthorId).Return(domain.Post{}, errors.New("something went wrong"))
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
			test.mockBehavior(post, auth, context.Background(), test.inputID, test.responsePost)
			fmt.Printf("test.inputPost: %v\n", test.inputID)
			handler := NewHandler(post, auth)
			// Init Endpoint
			r := gin.Default()
			reqID := fmt.Sprintf("/post/%d", test.inputID)
			if test.name == "Wrong input" {
				reqID = fmt.Sprintf("/post/%v", "wrong")
			}
			r.GET("post/:id", handler.GetById)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", reqID, nil)
			mockCookie := test.mockCookie
			req.AddCookie(mockCookie)
			// Make Request
			r.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_UpdateById(t *testing.T) {
	type mockBehavior func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers,
		ctx context.Context, inp domain.UpdatePost)
	var AuthorId int64 = 1
	tests := []struct {
		name                 string
		inputBody            string
		inputPost            domain.UpdatePost
		mockCookie           *http.Cookie
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"id":1, "title": "TestTitleNew", "body": "TestBodyNew"}`,
			inputPost: domain.UpdatePost{
				Id:    1,
				Title: "TestTitleNew",
				Body:  "TestBodyNew",
			},
			mockCookie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.UpdatePost) {
				mockUser.EXPECT().GetIdByToken(gomock.Any(), "refreshToken").Return(AuthorId, nil)
				mockPost.EXPECT().Update(gomock.Any(), inp.Id, &inp, AuthorId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"updated"}`,
		},
		{
			name:       "W/o Cookie",
			inputBody:  `{"id":1, "title": "TestTitleNew", "body": "TestBodyNew"}`,
			mockCookie: &http.Cookie{},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.UpdatePost) {
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"http: named cookie not present"}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"name": "username"}`,
			inputPost: domain.UpdatePost{},
			mockCookie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.UpdatePost) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input post body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"id":1, "title": "TestTitleNew", "body": "TestBodyNew"}`,
			inputPost: domain.UpdatePost{
				Id:    1,
				Title: "TestTitleNew",
				Body:  "TestBodyNew",
			},
			mockCookie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.UpdatePost) {
				mockUser.EXPECT().GetIdByToken(gomock.Any(), "refreshToken").Return(AuthorId, nil)
				mockPost.EXPECT().Update(gomock.Any(), inp.Id, &inp, AuthorId).Return(errors.New("something went wrong"))
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
			fmt.Printf("test.inputPost: %v\n", test.inputPost)
			handler := NewHandler(post, auth)
			// Init Endpoint
			r := gin.Default()
			r.PUT("/post/", handler.UpdateById)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/post/",
				bytes.NewBufferString(test.inputBody))
			mockCookie := test.mockCookie
			req.AddCookie(mockCookie)
			// Make Request
			r.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_DeleteById(t *testing.T) {
	type mockBehavior func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers,
		ctx context.Context, inp domain.UpdatePost)
	var AuthorId int64 = 1
	tests := []struct {
		name                 string
		inputBody            string
		inputPost            domain.UpdatePost
		mockCookie           *http.Cookie
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"id":1}`,
			inputPost: domain.UpdatePost{
				Id: 1,
			},
			mockCookie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.UpdatePost) {
				mockUser.EXPECT().GetIdByToken(gomock.Any(), "refreshToken").Return(AuthorId, nil)
				mockPost.EXPECT().Delete(gomock.Any(), inp.Id, AuthorId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"deleted"}`,
		},
		{
			name:       "W/o Cookie",
			inputBody:  `{"id":1}`,
			mockCookie: &http.Cookie{},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.UpdatePost) {
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"http: named cookie not present"}`,
		},
		{
			name:      "Wrong input",
			inputBody: `{"name": "username"}`,
			inputPost: domain.UpdatePost{},
			mockCookie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.UpdatePost) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input post body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"id":1}`,
			inputPost: domain.UpdatePost{
				Id: 1,
			},
			mockCookie: &http.Cookie{
				Name:   "refresh-token",
				Value:  "refreshToken",
				MaxAge: 2592000,
			},
			mockBehavior: func(mockPost *mocks.MockPosts, mockUser *mocks.MockUsers, ctx context.Context, inp domain.UpdatePost) {
				mockUser.EXPECT().GetIdByToken(gomock.Any(), "refreshToken").Return(AuthorId, nil)
				mockPost.EXPECT().Delete(gomock.Any(), inp.Id, AuthorId).Return(errors.New("something went wrong"))
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
			fmt.Printf("test.inputPost: %v\n", test.inputPost)
			handler := NewHandler(post, auth)
			// Init Endpoint
			r := gin.Default()
			r.DELETE("/post/", handler.DeleteById)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/post/",
				bytes.NewBufferString(test.inputBody))
			mockCookie := test.mockCookie
			req.AddCookie(mockCookie)
			// Make Request
			r.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
