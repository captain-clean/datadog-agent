// Code generated by mockery v1.0.0. DO NOT EDIT.

package compliance

import mock "github.com/stretchr/testify/mock"

// MockParser is an autogenerated mock type for the Parser type
type MockParser struct {
	mock.Mock
}

// Parse provides a mock function with given fields: config
func (_m *MockParser) Parse(config string) (*Suite, error) {
	ret := _m.Called(config)

	var r0 *Suite
	if rf, ok := ret.Get(0).(func(string) *Suite); ok {
		r0 = rf(config)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Suite)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}