package repository

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/monitoror/monitoror/monitorables/port/config"
	pkgNet "github.com/monitoror/monitoror/pkg/net"
	"github.com/monitoror/monitoror/pkg/net/mocks"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T, dialer pkgNet.Dialer) *portRepository {
	conf := &config.Port{
		Timeout: 1000,
	}
	repository := NewPortRepository(conf)

	systemPortRepository, ok := repository.(*portRepository)
	if assert.True(t, ok) {
		systemPortRepository.dialer = dialer
		return systemPortRepository
	}
	return nil
}

func TestRepository_OpenSocket_Success(t *testing.T) {
	mockConn := new(mocks.Conn)
	mockConn.On("Close").Return(nil)
	mockConn.On("SetDeadline", AnythingOfType("time.Time")).Return(nil)
	mockConn.On("Read", Anything).Return(1, nil)
	mockDialer := new(mocks.Dialer)
	mockDialer.On("Dial", AnythingOfType("string"), AnythingOfType("string")).Return(mockConn, nil)

	repository := initRepository(t, mockDialer)
	if repository != nil {
		responding, _, _, err := repository.OpenSocket("test", 1234, "tcp", nil)
		assert.NoError(t, err)
		assert.True(t, responding)
		mockConn.AssertNumberOfCalls(t, "Close", 1)
		mockConn.AssertExpectations(t)
		mockDialer.AssertNumberOfCalls(t, "Dial", 1)
		mockDialer.AssertExpectations(t)
	}
}

func TestRepository_OpenSocket_Success_WithPayload(t *testing.T) {
	mockConn := new(mocks.Conn)
	mockConn.On("Close").Return(nil)
	mockConn.On("SetDeadline", AnythingOfType("time.Time")).Return(nil)
	payload := []byte{0x01, 0x02}
	mockConn.On("Write", payload).Return(len(payload), nil)
	mockConn.On("Read", Anything).Return(1, nil)

	mockDialer := new(mocks.Dialer)
	mockDialer.On("Dial", AnythingOfType("string"), AnythingOfType("string")).Return(mockConn, nil)

	repository := initRepository(t, mockDialer)
	if repository != nil {
		responding, _, _, err := repository.OpenSocket("test", 1234, "udp", payload)
		assert.NoError(t, err)
		assert.True(t, responding)
		mockConn.AssertCalled(t, "Write", payload)
		mockConn.AssertNumberOfCalls(t, "Close", 1)
		mockDialer.AssertNumberOfCalls(t, "Dial", 1)
	}
}

func TestRepository_OpenSocket_Failed(t *testing.T) {
	mockDialer := new(mocks.Dialer)
	mockDialer.On("Dial", AnythingOfType("string"), AnythingOfType("string")).Return(nil, errors.New("check port failed"))

	repository := initRepository(t, mockDialer)
	if repository != nil {
		_, _, _, err := repository.OpenSocket("test", 1234, "tcp", nil)
		assert.Error(t, err)
		mockDialer.AssertNumberOfCalls(t, "Dial", 1)
		mockDialer.AssertExpectations(t)
	}
}

func TestRepository_OpenSocket_NoResponse(t *testing.T) {
	mockConn := new(mocks.Conn)
	mockConn.On("Close").Return(nil)
	mockConn.On("SetDeadline", AnythingOfType("time.Time")).Return(nil)
	mockConn.On("Read", Anything).Return(0, io.EOF)

	mockDialer := new(mocks.Dialer)
	mockDialer.On("Dial", AnythingOfType("string"), AnythingOfType("string")).Return(mockConn, nil)

	repository := initRepository(t, mockDialer)
	if repository != nil {
		responding, _, _, err := repository.OpenSocket("test", 1234, "tcp", nil)
		assert.NoError(t, err)
		assert.True(t, responding)
		mockConn.AssertNumberOfCalls(t, "Close", 1)
		mockDialer.AssertNumberOfCalls(t, "Dial", 1)
	}
}
