// Code generated by mockery v2.20.0. DO NOT EDIT.

package example

import (
	context "context"
	domainexample "demo/internal/models/example"

	mock "github.com/stretchr/testify/mock"
)

// MockExampleUsecase is an autogenerated mock type for the ExampleUsecase type
type MockExampleUsecase struct {
	mock.Mock
}

// ListItems provides a mock function with given fields: ctx, username, item
func (_m *MockExampleUsecase) ListItems(ctx context.Context, username string, item string) ([]*domainexample.Item, error) {
	ret := _m.Called(ctx, username, item)

	var r0 []*domainexample.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]*domainexample.Item, error)); ok {
		return rf(ctx, username, item)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*domainexample.Item); ok {
		r0 = rf(ctx, username, item)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domainexample.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, item)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, username, password
func (_m *MockExampleUsecase) Login(ctx context.Context, username string, password string) error {
	ret := _m.Called(ctx, username, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, username, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockExampleUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockExampleUsecase creates a new instance of MockExampleUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockExampleUsecase(t mockConstructorTestingTNewMockExampleUsecase) *MockExampleUsecase {
	mock := &MockExampleUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
