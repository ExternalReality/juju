// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/upgradeseries (interfaces: Facade,Logger,AgentService,ServiceAccess)

// Package upgradeseries_test is a generated GoMock package.
package upgradeseries_test

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/juju/juju/core/model"
	watcher "github.com/juju/juju/watcher"
	upgradeseries "github.com/juju/juju/worker/upgradeseries"
	loggo "github.com/juju/loggo"
	reflect "reflect"
)

// MockFacade is a mock of Facade interface
type MockFacade struct {
	ctrl     *gomock.Controller
	recorder *MockFacadeMockRecorder
}

// MockFacadeMockRecorder is the mock recorder for MockFacade
type MockFacadeMockRecorder struct {
	mock *MockFacade
}

// NewMockFacade creates a new mock instance
func NewMockFacade(ctrl *gomock.Controller) *MockFacade {
	mock := &MockFacade{ctrl: ctrl}
	mock.recorder = &MockFacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFacade) EXPECT() *MockFacadeMockRecorder {
	return m.recorder
}

// MachineStatus mocks base method
func (m *MockFacade) MachineStatus() (model.UpgradeSeriesStatus, error) {
	ret := m.ctrl.Call(m, "MachineStatus")
	ret0, _ := ret[0].(model.UpgradeSeriesStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MachineStatus indicates an expected call of MachineStatus
func (mr *MockFacadeMockRecorder) MachineStatus() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MachineStatus", reflect.TypeOf((*MockFacade)(nil).MachineStatus))
}

// SetMachineStatus mocks base method
func (m *MockFacade) SetMachineStatus(arg0 model.UpgradeSeriesStatus) error {
	ret := m.ctrl.Call(m, "SetMachineStatus", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetMachineStatus indicates an expected call of SetMachineStatus
func (mr *MockFacadeMockRecorder) SetMachineStatus(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMachineStatus", reflect.TypeOf((*MockFacade)(nil).SetMachineStatus), arg0)
}

// SetUpgradeSeriesStatus mocks base method
func (m *MockFacade) SetUpgradeSeriesStatus(arg0 string, arg1 model.UpgradeSeriesStatusType) error {
	ret := m.ctrl.Call(m, "SetUpgradeSeriesStatus", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUpgradeSeriesStatus indicates an expected call of SetUpgradeSeriesStatus
func (mr *MockFacadeMockRecorder) SetUpgradeSeriesStatus(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUpgradeSeriesStatus", reflect.TypeOf((*MockFacade)(nil).SetUpgradeSeriesStatus), arg0, arg1)
}

// UpgradeSeriesStatus mocks base method
func (m *MockFacade) UpgradeSeriesStatus(arg0 model.UpgradeSeriesStatusType) ([]string, error) {
	ret := m.ctrl.Call(m, "UpgradeSeriesStatus", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpgradeSeriesStatus indicates an expected call of UpgradeSeriesStatus
func (mr *MockFacadeMockRecorder) UpgradeSeriesStatus(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradeSeriesStatus", reflect.TypeOf((*MockFacade)(nil).UpgradeSeriesStatus), arg0)
}

// WatchUpgradeSeriesNotifications mocks base method
func (m *MockFacade) WatchUpgradeSeriesNotifications() (watcher.NotifyWatcher, error) {
	ret := m.ctrl.Call(m, "WatchUpgradeSeriesNotifications")
	ret0, _ := ret[0].(watcher.NotifyWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchUpgradeSeriesNotifications indicates an expected call of WatchUpgradeSeriesNotifications
func (mr *MockFacadeMockRecorder) WatchUpgradeSeriesNotifications() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchUpgradeSeriesNotifications", reflect.TypeOf((*MockFacade)(nil).WatchUpgradeSeriesNotifications))
}

// MockLogger is a mock of Logger interface
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Errorf mocks base method
func (m *MockLogger) Errorf(arg0 string, arg1 ...interface{}) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf
func (mr *MockLoggerMockRecorder) Errorf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockLogger)(nil).Errorf), varargs...)
}

// Logf mocks base method
func (m *MockLogger) Logf(arg0 loggo.Level, arg1 string, arg2 ...interface{}) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Logf", varargs...)
}

// Logf indicates an expected call of Logf
func (mr *MockLoggerMockRecorder) Logf(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logf", reflect.TypeOf((*MockLogger)(nil).Logf), varargs...)
}

// Warningf mocks base method
func (m *MockLogger) Warningf(arg0 string, arg1 ...interface{}) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warningf", varargs...)
}

// Warningf indicates an expected call of Warningf
func (mr *MockLoggerMockRecorder) Warningf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warningf", reflect.TypeOf((*MockLogger)(nil).Warningf), varargs...)
}

// MockAgentService is a mock of AgentService interface
type MockAgentService struct {
	ctrl     *gomock.Controller
	recorder *MockAgentServiceMockRecorder
}

// MockAgentServiceMockRecorder is the mock recorder for MockAgentService
type MockAgentServiceMockRecorder struct {
	mock *MockAgentService
}

// NewMockAgentService creates a new mock instance
func NewMockAgentService(ctrl *gomock.Controller) *MockAgentService {
	mock := &MockAgentService{ctrl: ctrl}
	mock.recorder = &MockAgentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAgentService) EXPECT() *MockAgentServiceMockRecorder {
	return m.recorder
}

// Running mocks base method
func (m *MockAgentService) Running() (bool, error) {
	ret := m.ctrl.Call(m, "Running")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Running indicates an expected call of Running
func (mr *MockAgentServiceMockRecorder) Running() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Running", reflect.TypeOf((*MockAgentService)(nil).Running))
}

// Start mocks base method
func (m *MockAgentService) Start() error {
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockAgentServiceMockRecorder) Start() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockAgentService)(nil).Start))
}

// Stop mocks base method
func (m *MockAgentService) Stop() error {
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockAgentServiceMockRecorder) Stop() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockAgentService)(nil).Stop))
}

// MockServiceAccess is a mock of ServiceAccess interface
type MockServiceAccess struct {
	ctrl     *gomock.Controller
	recorder *MockServiceAccessMockRecorder
}

// MockServiceAccessMockRecorder is the mock recorder for MockServiceAccess
type MockServiceAccessMockRecorder struct {
	mock *MockServiceAccess
}

// NewMockServiceAccess creates a new mock instance
func NewMockServiceAccess(ctrl *gomock.Controller) *MockServiceAccess {
	mock := &MockServiceAccess{ctrl: ctrl}
	mock.recorder = &MockServiceAccessMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockServiceAccess) EXPECT() *MockServiceAccessMockRecorder {
	return m.recorder
}

// DiscoverService mocks base method
func (m *MockServiceAccess) DiscoverService(arg0 string) (upgradeseries.AgentService, error) {
	ret := m.ctrl.Call(m, "DiscoverService", arg0)
	ret0, _ := ret[0].(upgradeseries.AgentService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DiscoverService indicates an expected call of DiscoverService
func (mr *MockServiceAccessMockRecorder) DiscoverService(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiscoverService", reflect.TypeOf((*MockServiceAccess)(nil).DiscoverService), arg0)
}

// ListServices mocks base method
func (m *MockServiceAccess) ListServices() ([]string, error) {
	ret := m.ctrl.Call(m, "ListServices")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListServices indicates an expected call of ListServices
func (mr *MockServiceAccessMockRecorder) ListServices() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListServices", reflect.TypeOf((*MockServiceAccess)(nil).ListServices))
}
