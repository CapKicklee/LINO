// Code generated by mockery v1.0.0. DO NOT EDIT.

package extract

import mock "github.com/stretchr/testify/mock"

// MockStepList is an autogenerated mock type for the StepList type
type MockStepList struct {
	mock.Mock
}

// Len provides a mock function with given fields:
func (_m *MockStepList) Len() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// Step provides a mock function with given fields: _a0
func (_m *MockStepList) Step(_a0 uint) Step {
	ret := _m.Called(_a0)

	var r0 Step
	if rf, ok := ret.Get(0).(func(uint) Step); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Step)
		}
	}

	return r0
}
