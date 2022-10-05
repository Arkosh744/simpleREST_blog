// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/Arkosh744/simpleREST_blog/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockPosts is a mock of Posts interface.
type MockPosts struct {
	ctrl     *gomock.Controller
	recorder *MockPostsMockRecorder
}

// MockPostsMockRecorder is the mock recorder for MockPosts.
type MockPostsMockRecorder struct {
	mock *MockPosts
}

// NewMockPosts creates a new mock instance.
func NewMockPosts(ctrl *gomock.Controller) *MockPosts {
	mock := &MockPosts{ctrl: ctrl}
	mock.recorder = &MockPostsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPosts) EXPECT() *MockPostsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPosts) Create(ctx context.Context, post domain.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockPostsMockRecorder) Create(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPosts)(nil).Create), ctx, post)
}

// Delete mocks base method.
func (m *MockPosts) Delete(ctx context.Context, id, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPostsMockRecorder) Delete(ctx, id, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPosts)(nil).Delete), ctx, id, userId)
}

// GetById mocks base method.
func (m *MockPosts) GetById(ctx context.Context, id, userId int64) (domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, id, userId)
	ret0, _ := ret[0].(domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockPostsMockRecorder) GetById(ctx, id, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockPosts)(nil).GetById), ctx, id, userId)
}

// List mocks base method.
func (m *MockPosts) List(ctx context.Context, userId int64) ([]domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, userId)
	ret0, _ := ret[0].([]domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockPostsMockRecorder) List(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPosts)(nil).List), ctx, userId)
}

// Update mocks base method.
func (m *MockPosts) Update(ctx context.Context, id int64, post *domain.UpdatePost, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, post, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPostsMockRecorder) Update(ctx, id, post, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPosts)(nil).Update), ctx, id, post, userId)
}

// MockUsers is a mock of Users interface.
type MockUsers struct {
	ctrl     *gomock.Controller
	recorder *MockUsersMockRecorder
}

// MockUsersMockRecorder is the mock recorder for MockUsers.
type MockUsersMockRecorder struct {
	mock *MockUsers
}

// NewMockUsers creates a new mock instance.
func NewMockUsers(ctrl *gomock.Controller) *MockUsers {
	mock := &MockUsers{ctrl: ctrl}
	mock.recorder = &MockUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsers) EXPECT() *MockUsersMockRecorder {
	return m.recorder
}

// GetIdByToken mocks base method.
func (m *MockUsers) GetIdByToken(ctx context.Context, refreshToken string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIdByToken", ctx, refreshToken)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIdByToken indicates an expected call of GetIdByToken.
func (mr *MockUsersMockRecorder) GetIdByToken(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIdByToken", reflect.TypeOf((*MockUsers)(nil).GetIdByToken), ctx, refreshToken)
}

// ParseToken mocks base method.
func (m *MockUsers) ParseToken(ctx context.Context, token string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", ctx, token)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockUsersMockRecorder) ParseToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockUsers)(nil).ParseToken), ctx, token)
}

// RefreshTokens mocks base method.
func (m *MockUsers) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshTokens", ctx, refreshToken)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RefreshTokens indicates an expected call of RefreshTokens.
func (mr *MockUsersMockRecorder) RefreshTokens(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshTokens", reflect.TypeOf((*MockUsers)(nil).RefreshTokens), ctx, refreshToken)
}

// SignIn mocks base method.
func (m *MockUsers) SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", ctx, inp)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SignIn indicates an expected call of SignIn.
func (mr *MockUsersMockRecorder) SignIn(ctx, inp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockUsers)(nil).SignIn), ctx, inp)
}

// SignUp mocks base method.
func (m *MockUsers) SignUp(ctx context.Context, inp domain.SignUpInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, inp)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUsersMockRecorder) SignUp(ctx, inp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUsers)(nil).SignUp), ctx, inp)
}