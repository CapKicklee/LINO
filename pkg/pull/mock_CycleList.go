// Code generated by mockery v1.0.0. DO NOT EDIT.

package pull

import mock "github.com/stretchr/testify/mock"

// MockCycleList is an autogenerated mock type for the CycleList type
type MockCycleList struct {
	mock.Mock
}

// Cycle provides a mock function with given fields: idx
func (_m *MockCycleList) Cycle(idx uint) Cycle {
	ret := _m.Called(idx)

	var r0 Cycle
	if rf, ok := ret.Get(0).(func(uint) Cycle); ok {
		r0 = rf(idx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Cycle)
		}
	}

	return r0
}

// Len provides a mock function with given fields:
func (_m *MockCycleList) Len() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}
