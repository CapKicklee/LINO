// Code generated by mockery v1.0.0. DO NOT EDIT.

package id

import mock "github.com/stretchr/testify/mock"

// MockRelationReader is an autogenerated mock type for the RelationReader type
type MockRelationReader struct {
	mock.Mock
}

// Read provides a mock function with given fields:
func (_m *MockRelationReader) Read() (RelationList, *Error) {
	ret := _m.Called()

	var r0 RelationList
	if rf, ok := ret.Get(0).(func() RelationList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(RelationList)
		}
	}

	var r1 *Error
	if rf, ok := ret.Get(1).(func() *Error); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*Error)
		}
	}

	return r0, r1
}
