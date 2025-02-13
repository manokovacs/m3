// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3/src/cmd/services/m3coordinator/ingest (interfaces: DownsamplerAndWriter)

// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package ingest is a generated GoMock package.
package ingest

import (
	"context"
	"reflect"

	"github.com/m3db/m3/src/query/models"
	"github.com/m3db/m3/src/query/storage"
	"github.com/m3db/m3/src/query/ts"
	"github.com/m3db/m3/src/x/time"

	"github.com/golang/mock/gomock"
)

// MockDownsamplerAndWriter is a mock of DownsamplerAndWriter interface
type MockDownsamplerAndWriter struct {
	ctrl     *gomock.Controller
	recorder *MockDownsamplerAndWriterMockRecorder
}

// MockDownsamplerAndWriterMockRecorder is the mock recorder for MockDownsamplerAndWriter
type MockDownsamplerAndWriterMockRecorder struct {
	mock *MockDownsamplerAndWriter
}

// NewMockDownsamplerAndWriter creates a new mock instance
func NewMockDownsamplerAndWriter(ctrl *gomock.Controller) *MockDownsamplerAndWriter {
	mock := &MockDownsamplerAndWriter{ctrl: ctrl}
	mock.recorder = &MockDownsamplerAndWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDownsamplerAndWriter) EXPECT() *MockDownsamplerAndWriterMockRecorder {
	return m.recorder
}

// Storage mocks base method
func (m *MockDownsamplerAndWriter) Storage() storage.Storage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Storage")
	ret0, _ := ret[0].(storage.Storage)
	return ret0
}

// Storage indicates an expected call of Storage
func (mr *MockDownsamplerAndWriterMockRecorder) Storage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Storage", reflect.TypeOf((*MockDownsamplerAndWriter)(nil).Storage))
}

// Write mocks base method
func (m *MockDownsamplerAndWriter) Write(arg0 context.Context, arg1 models.Tags, arg2 ts.Datapoints, arg3 time.Unit, arg4 WriteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write
func (mr *MockDownsamplerAndWriterMockRecorder) Write(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockDownsamplerAndWriter)(nil).Write), arg0, arg1, arg2, arg3, arg4)
}

// WriteBatch mocks base method
func (m *MockDownsamplerAndWriter) WriteBatch(arg0 context.Context, arg1 DownsampleAndWriteIter) BatchError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteBatch", arg0, arg1)
	ret0, _ := ret[0].(BatchError)
	return ret0
}

// WriteBatch indicates an expected call of WriteBatch
func (mr *MockDownsamplerAndWriterMockRecorder) WriteBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteBatch", reflect.TypeOf((*MockDownsamplerAndWriter)(nil).WriteBatch), arg0, arg1)
}
