// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3/src/dbnode/storage/series (interfaces: DatabaseSeries,QueryableBlockRetriever)

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

// Package series is a generated GoMock package.
package series

import (
	"reflect"
	"time"

	"github.com/m3db/m3/src/dbnode/namespace"
	"github.com/m3db/m3/src/dbnode/persist"
	"github.com/m3db/m3/src/dbnode/storage/block"
	"github.com/m3db/m3/src/dbnode/ts"
	"github.com/m3db/m3/src/dbnode/x/xio"
	"github.com/m3db/m3/src/x/context"
	"github.com/m3db/m3/src/x/ident"
	time0 "github.com/m3db/m3/src/x/time"

	"github.com/golang/mock/gomock"
)

// MockDatabaseSeries is a mock of DatabaseSeries interface
type MockDatabaseSeries struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseSeriesMockRecorder
}

// MockDatabaseSeriesMockRecorder is the mock recorder for MockDatabaseSeries
type MockDatabaseSeriesMockRecorder struct {
	mock *MockDatabaseSeries
}

// NewMockDatabaseSeries creates a new mock instance
func NewMockDatabaseSeries(ctrl *gomock.Controller) *MockDatabaseSeries {
	mock := &MockDatabaseSeries{ctrl: ctrl}
	mock.recorder = &MockDatabaseSeriesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDatabaseSeries) EXPECT() *MockDatabaseSeriesMockRecorder {
	return m.recorder
}

// Bootstrap mocks base method
func (m *MockDatabaseSeries) Bootstrap(arg0 block.DatabaseSeriesBlocks) (BootstrapResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bootstrap", arg0)
	ret0, _ := ret[0].(BootstrapResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Bootstrap indicates an expected call of Bootstrap
func (mr *MockDatabaseSeriesMockRecorder) Bootstrap(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bootstrap", reflect.TypeOf((*MockDatabaseSeries)(nil).Bootstrap), arg0)
}

// Close mocks base method
func (m *MockDatabaseSeries) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockDatabaseSeriesMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDatabaseSeries)(nil).Close))
}

// ColdFlushBlockStarts mocks base method
func (m *MockDatabaseSeries) ColdFlushBlockStarts(arg0 map[time0.UnixNano]BlockState) OptimizedTimes {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ColdFlushBlockStarts", arg0)
	ret0, _ := ret[0].(OptimizedTimes)
	return ret0
}

// ColdFlushBlockStarts indicates an expected call of ColdFlushBlockStarts
func (mr *MockDatabaseSeriesMockRecorder) ColdFlushBlockStarts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ColdFlushBlockStarts", reflect.TypeOf((*MockDatabaseSeries)(nil).ColdFlushBlockStarts), arg0)
}

// FetchBlocks mocks base method
func (m *MockDatabaseSeries) FetchBlocks(arg0 context.Context, arg1 []time.Time, arg2 namespace.Context) ([]block.FetchBlockResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchBlocks", arg0, arg1, arg2)
	ret0, _ := ret[0].([]block.FetchBlockResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchBlocks indicates an expected call of FetchBlocks
func (mr *MockDatabaseSeriesMockRecorder) FetchBlocks(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchBlocks", reflect.TypeOf((*MockDatabaseSeries)(nil).FetchBlocks), arg0, arg1, arg2)
}

// FetchBlocksForColdFlush mocks base method
func (m *MockDatabaseSeries) FetchBlocksForColdFlush(arg0 context.Context, arg1 time.Time, arg2 int, arg3 namespace.Context) ([]xio.BlockReader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchBlocksForColdFlush", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]xio.BlockReader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchBlocksForColdFlush indicates an expected call of FetchBlocksForColdFlush
func (mr *MockDatabaseSeriesMockRecorder) FetchBlocksForColdFlush(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchBlocksForColdFlush", reflect.TypeOf((*MockDatabaseSeries)(nil).FetchBlocksForColdFlush), arg0, arg1, arg2, arg3)
}

// FetchBlocksMetadata mocks base method
func (m *MockDatabaseSeries) FetchBlocksMetadata(arg0 context.Context, arg1, arg2 time.Time, arg3 FetchBlocksMetadataOptions) (block.FetchBlocksMetadataResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchBlocksMetadata", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(block.FetchBlocksMetadataResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchBlocksMetadata indicates an expected call of FetchBlocksMetadata
func (mr *MockDatabaseSeriesMockRecorder) FetchBlocksMetadata(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchBlocksMetadata", reflect.TypeOf((*MockDatabaseSeries)(nil).FetchBlocksMetadata), arg0, arg1, arg2, arg3)
}

// ID mocks base method
func (m *MockDatabaseSeries) ID() ident.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(ident.ID)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockDatabaseSeriesMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockDatabaseSeries)(nil).ID))
}

// IsBootstrapped mocks base method
func (m *MockDatabaseSeries) IsBootstrapped() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsBootstrapped")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsBootstrapped indicates an expected call of IsBootstrapped
func (mr *MockDatabaseSeriesMockRecorder) IsBootstrapped() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsBootstrapped", reflect.TypeOf((*MockDatabaseSeries)(nil).IsBootstrapped))
}

// IsEmpty mocks base method
func (m *MockDatabaseSeries) IsEmpty() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmpty")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsEmpty indicates an expected call of IsEmpty
func (mr *MockDatabaseSeriesMockRecorder) IsEmpty() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmpty", reflect.TypeOf((*MockDatabaseSeries)(nil).IsEmpty))
}

// NumActiveBlocks mocks base method
func (m *MockDatabaseSeries) NumActiveBlocks() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NumActiveBlocks")
	ret0, _ := ret[0].(int)
	return ret0
}

// NumActiveBlocks indicates an expected call of NumActiveBlocks
func (mr *MockDatabaseSeriesMockRecorder) NumActiveBlocks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NumActiveBlocks", reflect.TypeOf((*MockDatabaseSeries)(nil).NumActiveBlocks))
}

// OnEvictedFromWiredList mocks base method
func (m *MockDatabaseSeries) OnEvictedFromWiredList(arg0 ident.ID, arg1 time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnEvictedFromWiredList", arg0, arg1)
}

// OnEvictedFromWiredList indicates an expected call of OnEvictedFromWiredList
func (mr *MockDatabaseSeriesMockRecorder) OnEvictedFromWiredList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnEvictedFromWiredList", reflect.TypeOf((*MockDatabaseSeries)(nil).OnEvictedFromWiredList), arg0, arg1)
}

// OnRetrieveBlock mocks base method
func (m *MockDatabaseSeries) OnRetrieveBlock(arg0 ident.ID, arg1 ident.TagIterator, arg2 time.Time, arg3 ts.Segment, arg4 namespace.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnRetrieveBlock", arg0, arg1, arg2, arg3, arg4)
}

// OnRetrieveBlock indicates an expected call of OnRetrieveBlock
func (mr *MockDatabaseSeriesMockRecorder) OnRetrieveBlock(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnRetrieveBlock", reflect.TypeOf((*MockDatabaseSeries)(nil).OnRetrieveBlock), arg0, arg1, arg2, arg3, arg4)
}

// ReadEncoded mocks base method
func (m *MockDatabaseSeries) ReadEncoded(arg0 context.Context, arg1, arg2 time.Time, arg3 namespace.Context) ([][]xio.BlockReader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadEncoded", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([][]xio.BlockReader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadEncoded indicates an expected call of ReadEncoded
func (mr *MockDatabaseSeriesMockRecorder) ReadEncoded(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadEncoded", reflect.TypeOf((*MockDatabaseSeries)(nil).ReadEncoded), arg0, arg1, arg2, arg3)
}

// Reset mocks base method
func (m *MockDatabaseSeries) Reset(arg0 ident.ID, arg1 ident.Tags, arg2 QueryableBlockRetriever, arg3 block.OnRetrieveBlock, arg4 block.OnEvictedFromWiredList, arg5 Options) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset", arg0, arg1, arg2, arg3, arg4, arg5)
}

// Reset indicates an expected call of Reset
func (mr *MockDatabaseSeriesMockRecorder) Reset(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockDatabaseSeries)(nil).Reset), arg0, arg1, arg2, arg3, arg4, arg5)
}

// Snapshot mocks base method
func (m *MockDatabaseSeries) Snapshot(arg0 context.Context, arg1 time.Time, arg2 persist.DataFn, arg3 namespace.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Snapshot", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Snapshot indicates an expected call of Snapshot
func (mr *MockDatabaseSeriesMockRecorder) Snapshot(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Snapshot", reflect.TypeOf((*MockDatabaseSeries)(nil).Snapshot), arg0, arg1, arg2, arg3)
}

// Tags mocks base method
func (m *MockDatabaseSeries) Tags() ident.Tags {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tags")
	ret0, _ := ret[0].(ident.Tags)
	return ret0
}

// Tags indicates an expected call of Tags
func (mr *MockDatabaseSeriesMockRecorder) Tags() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tags", reflect.TypeOf((*MockDatabaseSeries)(nil).Tags))
}

// Tick mocks base method
func (m *MockDatabaseSeries) Tick(arg0 map[time0.UnixNano]BlockState, arg1 namespace.Context) (TickResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tick", arg0, arg1)
	ret0, _ := ret[0].(TickResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Tick indicates an expected call of Tick
func (mr *MockDatabaseSeriesMockRecorder) Tick(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tick", reflect.TypeOf((*MockDatabaseSeries)(nil).Tick), arg0, arg1)
}

// WarmFlush mocks base method
func (m *MockDatabaseSeries) WarmFlush(arg0 context.Context, arg1 time.Time, arg2 persist.DataFn, arg3 namespace.Context) (FlushOutcome, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WarmFlush", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(FlushOutcome)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WarmFlush indicates an expected call of WarmFlush
func (mr *MockDatabaseSeriesMockRecorder) WarmFlush(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WarmFlush", reflect.TypeOf((*MockDatabaseSeries)(nil).WarmFlush), arg0, arg1, arg2, arg3)
}

// Write mocks base method
func (m *MockDatabaseSeries) Write(arg0 context.Context, arg1 time.Time, arg2 float64, arg3 time0.Unit, arg4 []byte, arg5 WriteOptions) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write
func (mr *MockDatabaseSeriesMockRecorder) Write(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockDatabaseSeries)(nil).Write), arg0, arg1, arg2, arg3, arg4, arg5)
}

// MockQueryableBlockRetriever is a mock of QueryableBlockRetriever interface
type MockQueryableBlockRetriever struct {
	ctrl     *gomock.Controller
	recorder *MockQueryableBlockRetrieverMockRecorder
}

// MockQueryableBlockRetrieverMockRecorder is the mock recorder for MockQueryableBlockRetriever
type MockQueryableBlockRetrieverMockRecorder struct {
	mock *MockQueryableBlockRetriever
}

// NewMockQueryableBlockRetriever creates a new mock instance
func NewMockQueryableBlockRetriever(ctrl *gomock.Controller) *MockQueryableBlockRetriever {
	mock := &MockQueryableBlockRetriever{ctrl: ctrl}
	mock.recorder = &MockQueryableBlockRetrieverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockQueryableBlockRetriever) EXPECT() *MockQueryableBlockRetrieverMockRecorder {
	return m.recorder
}

// BlockStatesSnapshot mocks base method
func (m *MockQueryableBlockRetriever) BlockStatesSnapshot() map[time0.UnixNano]BlockState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockStatesSnapshot")
	ret0, _ := ret[0].(map[time0.UnixNano]BlockState)
	return ret0
}

// BlockStatesSnapshot indicates an expected call of BlockStatesSnapshot
func (mr *MockQueryableBlockRetrieverMockRecorder) BlockStatesSnapshot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockStatesSnapshot", reflect.TypeOf((*MockQueryableBlockRetriever)(nil).BlockStatesSnapshot))
}

// IsBlockRetrievable mocks base method
func (m *MockQueryableBlockRetriever) IsBlockRetrievable(arg0 time.Time) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsBlockRetrievable", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsBlockRetrievable indicates an expected call of IsBlockRetrievable
func (mr *MockQueryableBlockRetrieverMockRecorder) IsBlockRetrievable(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsBlockRetrievable", reflect.TypeOf((*MockQueryableBlockRetriever)(nil).IsBlockRetrievable), arg0)
}

// RetrievableBlockColdVersion mocks base method
func (m *MockQueryableBlockRetriever) RetrievableBlockColdVersion(arg0 time.Time) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrievableBlockColdVersion", arg0)
	ret0, _ := ret[0].(int)
	return ret0
}

// RetrievableBlockColdVersion indicates an expected call of RetrievableBlockColdVersion
func (mr *MockQueryableBlockRetrieverMockRecorder) RetrievableBlockColdVersion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrievableBlockColdVersion", reflect.TypeOf((*MockQueryableBlockRetriever)(nil).RetrievableBlockColdVersion), arg0)
}

// Stream mocks base method
func (m *MockQueryableBlockRetriever) Stream(arg0 context.Context, arg1 ident.ID, arg2 time.Time, arg3 block.OnRetrieveBlock, arg4 namespace.Context) (xio.BlockReader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stream", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(xio.BlockReader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stream indicates an expected call of Stream
func (mr *MockQueryableBlockRetrieverMockRecorder) Stream(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stream", reflect.TypeOf((*MockQueryableBlockRetriever)(nil).Stream), arg0, arg1, arg2, arg3, arg4)
}
