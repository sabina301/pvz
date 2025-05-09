// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	models "pvz/internal/models"

	mock "github.com/stretchr/testify/mock"

	repositories "pvz/internal/repositories"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

type UserRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *UserRepository) EXPECT() *UserRepository_Expecter {
	return &UserRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: q, user
func (_m *UserRepository) Create(q repositories.Querier, user models.User) (string, error) {
	ret := _m.Called(q, user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(repositories.Querier, models.User) (string, error)); ok {
		return rf(q, user)
	}
	if rf, ok := ret.Get(0).(func(repositories.Querier, models.User) string); ok {
		r0 = rf(q, user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(repositories.Querier, models.User) error); ok {
		r1 = rf(q, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - q repositories.Querier
//   - user models.User
func (_e *UserRepository_Expecter) Create(q interface{}, user interface{}) *UserRepository_Create_Call {
	return &UserRepository_Create_Call{Call: _e.mock.On("Create", q, user)}
}

func (_c *UserRepository_Create_Call) Run(run func(q repositories.Querier, user models.User)) *UserRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(repositories.Querier), args[1].(models.User))
	})
	return _c
}

func (_c *UserRepository_Create_Call) Return(_a0 string, _a1 error) *UserRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_Create_Call) RunAndReturn(run func(repositories.Querier, models.User) (string, error)) *UserRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetByEmail provides a mock function with given fields: q, email
func (_m *UserRepository) GetByEmail(q repositories.Querier, email string) (*models.User, error) {
	ret := _m.Called(q, email)

	if len(ret) == 0 {
		panic("no return value specified for GetByEmail")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(repositories.Querier, string) (*models.User, error)); ok {
		return rf(q, email)
	}
	if rf, ok := ret.Get(0).(func(repositories.Querier, string) *models.User); ok {
		r0 = rf(q, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(repositories.Querier, string) error); ok {
		r1 = rf(q, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEmail'
type UserRepository_GetByEmail_Call struct {
	*mock.Call
}

// GetByEmail is a helper method to define mock.On call
//   - q repositories.Querier
//   - email string
func (_e *UserRepository_Expecter) GetByEmail(q interface{}, email interface{}) *UserRepository_GetByEmail_Call {
	return &UserRepository_GetByEmail_Call{Call: _e.mock.On("GetByEmail", q, email)}
}

func (_c *UserRepository_GetByEmail_Call) Run(run func(q repositories.Querier, email string)) *UserRepository_GetByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(repositories.Querier), args[1].(string))
	})
	return _c
}

func (_c *UserRepository_GetByEmail_Call) Return(_a0 *models.User, _a1 error) *UserRepository_GetByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetByEmail_Call) RunAndReturn(run func(repositories.Querier, string) (*models.User, error)) *UserRepository_GetByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
