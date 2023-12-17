// Code generated by MockGen. DO NOT EDIT.
// Source: internal/user/service.go
//
// Generated by this command:
//
//	mockgen -source=internal/user/service.go -destination internal/user/mock/service_mock.go
//
// Package mock_user is a generated GoMock package.
package mock_user

import (
	context "context"
	multipart "mime/multipart"
	reflect "reflect"

	domain "github.com/blazee5/quizmaster-backend/internal/domain"
	models "github.com/blazee5/quizmaster-backend/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// ChangeAvatar mocks base method.
func (m *MockService) ChangeAvatar(ctx context.Context, userID int, fileHeader *multipart.FileHeader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAvatar", ctx, userID, fileHeader)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeAvatar indicates an expected call of ChangeAvatar.
func (mr *MockServiceMockRecorder) ChangeAvatar(ctx, userID, fileHeader any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAvatar", reflect.TypeOf((*MockService)(nil).ChangeAvatar), ctx, userID, fileHeader)
}

// Delete mocks base method.
func (m *MockService) Delete(ctx context.Context, userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder) Delete(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService)(nil).Delete), ctx, userID)
}

// GetByID mocks base method.
func (m *MockService) GetByID(ctx context.Context, userID int) (models.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, userID)
	ret0, _ := ret[0].(models.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockServiceMockRecorder) GetByID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockService)(nil).GetByID), ctx, userID)
}

// Update mocks base method.
func (m *MockService) Update(ctx context.Context, userID int, input domain.UpdateUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, userID, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockServiceMockRecorder) Update(ctx, userID, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockService)(nil).Update), ctx, userID, input)
}
