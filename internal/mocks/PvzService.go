// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	pvz "pvz/internal/models/pvz"

	mock "github.com/stretchr/testify/mock"
)

// PvzService is an autogenerated mock type for the PvzService type
type PvzService struct {
	mock.Mock
}

type PvzService_Expecter struct {
	mock *mock.Mock
}

func (_m *PvzService) EXPECT() *PvzService_Expecter {
	return &PvzService_Expecter{mock: &_m.Mock}
}

// CLoseLastReception provides a mock function with given fields: pvzId
func (_m *PvzService) CLoseLastReception(pvzId string) (pvz.CloseLastProductResponse, error) {
	ret := _m.Called(pvzId)

	if len(ret) == 0 {
		panic("no return value specified for CLoseLastReception")
	}

	var r0 pvz.CloseLastProductResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (pvz.CloseLastProductResponse, error)); ok {
		return rf(pvzId)
	}
	if rf, ok := ret.Get(0).(func(string) pvz.CloseLastProductResponse); ok {
		r0 = rf(pvzId)
	} else {
		r0 = ret.Get(0).(pvz.CloseLastProductResponse)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(pvzId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PvzService_CLoseLastReception_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CLoseLastReception'
type PvzService_CLoseLastReception_Call struct {
	*mock.Call
}

// CLoseLastReception is a helper method to define mock.On call
//   - pvzId string
func (_e *PvzService_Expecter) CLoseLastReception(pvzId interface{}) *PvzService_CLoseLastReception_Call {
	return &PvzService_CLoseLastReception_Call{Call: _e.mock.On("CLoseLastReception", pvzId)}
}

func (_c *PvzService_CLoseLastReception_Call) Run(run func(pvzId string)) *PvzService_CLoseLastReception_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PvzService_CLoseLastReception_Call) Return(_a0 pvz.CloseLastProductResponse, _a1 error) *PvzService_CLoseLastReception_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PvzService_CLoseLastReception_Call) RunAndReturn(run func(string) (pvz.CloseLastProductResponse, error)) *PvzService_CLoseLastReception_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: req
func (_m *PvzService) Create(req pvz.CreateRequest) (pvz.CreateResponse, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 pvz.CreateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(pvz.CreateRequest) (pvz.CreateResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(pvz.CreateRequest) pvz.CreateResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(pvz.CreateResponse)
	}

	if rf, ok := ret.Get(1).(func(pvz.CreateRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PvzService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type PvzService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - req pvz.CreateRequest
func (_e *PvzService_Expecter) Create(req interface{}) *PvzService_Create_Call {
	return &PvzService_Create_Call{Call: _e.mock.On("Create", req)}
}

func (_c *PvzService_Create_Call) Run(run func(req pvz.CreateRequest)) *PvzService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(pvz.CreateRequest))
	})
	return _c
}

func (_c *PvzService_Create_Call) Return(_a0 pvz.CreateResponse, _a1 error) *PvzService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PvzService_Create_Call) RunAndReturn(run func(pvz.CreateRequest) (pvz.CreateResponse, error)) *PvzService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteLastProduct provides a mock function with given fields: pvzId
func (_m *PvzService) DeleteLastProduct(pvzId string) (pvz.DeleteLastProductResponse, error) {
	ret := _m.Called(pvzId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteLastProduct")
	}

	var r0 pvz.DeleteLastProductResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (pvz.DeleteLastProductResponse, error)); ok {
		return rf(pvzId)
	}
	if rf, ok := ret.Get(0).(func(string) pvz.DeleteLastProductResponse); ok {
		r0 = rf(pvzId)
	} else {
		r0 = ret.Get(0).(pvz.DeleteLastProductResponse)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(pvzId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PvzService_DeleteLastProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteLastProduct'
type PvzService_DeleteLastProduct_Call struct {
	*mock.Call
}

// DeleteLastProduct is a helper method to define mock.On call
//   - pvzId string
func (_e *PvzService_Expecter) DeleteLastProduct(pvzId interface{}) *PvzService_DeleteLastProduct_Call {
	return &PvzService_DeleteLastProduct_Call{Call: _e.mock.On("DeleteLastProduct", pvzId)}
}

func (_c *PvzService_DeleteLastProduct_Call) Run(run func(pvzId string)) *PvzService_DeleteLastProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PvzService_DeleteLastProduct_Call) Return(_a0 pvz.DeleteLastProductResponse, _a1 error) *PvzService_DeleteLastProduct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PvzService_DeleteLastProduct_Call) RunAndReturn(run func(string) (pvz.DeleteLastProductResponse, error)) *PvzService_DeleteLastProduct_Call {
	_c.Call.Return(run)
	return _c
}

// ListWithFilterDate provides a mock function with given fields: req
func (_m *PvzService) ListWithFilterDate(req pvz.ListRequest) ([]pvz.ListResponse, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for ListWithFilterDate")
	}

	var r0 []pvz.ListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(pvz.ListRequest) ([]pvz.ListResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(pvz.ListRequest) []pvz.ListResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pvz.ListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(pvz.ListRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PvzService_ListWithFilterDate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListWithFilterDate'
type PvzService_ListWithFilterDate_Call struct {
	*mock.Call
}

// ListWithFilterDate is a helper method to define mock.On call
//   - req pvz.ListRequest
func (_e *PvzService_Expecter) ListWithFilterDate(req interface{}) *PvzService_ListWithFilterDate_Call {
	return &PvzService_ListWithFilterDate_Call{Call: _e.mock.On("ListWithFilterDate", req)}
}

func (_c *PvzService_ListWithFilterDate_Call) Run(run func(req pvz.ListRequest)) *PvzService_ListWithFilterDate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(pvz.ListRequest))
	})
	return _c
}

func (_c *PvzService_ListWithFilterDate_Call) Return(_a0 []pvz.ListResponse, _a1 error) *PvzService_ListWithFilterDate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PvzService_ListWithFilterDate_Call) RunAndReturn(run func(pvz.ListRequest) ([]pvz.ListResponse, error)) *PvzService_ListWithFilterDate_Call {
	_c.Call.Return(run)
	return _c
}

// NewPvzService creates a new instance of PvzService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPvzService(t interface {
	mock.TestingT
	Cleanup(func())
}) *PvzService {
	mock := &PvzService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
