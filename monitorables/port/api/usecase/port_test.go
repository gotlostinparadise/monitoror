package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api"
	"github.com/monitoror/monitoror/monitorables/port/api/mocks"
	"github.com/monitoror/monitoror/monitorables/port/api/models"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestUsecase_CheckPort_Success(t *testing.T) {
	// Init
	mockRepo := new(mocks.Repository)
	mockRepo.On("OpenSocket", AnythingOfType("string"), AnythingOfType("int"), AnythingOfType("string"), Anything).Return(true, "SSH-2.0-openSSH", time.Millisecond*50, nil)
	usecase := NewPortUsecase(mockRepo)

	// Params
	param := &models.PortParams{
		Hostname: "monitoror.example.com",
		Port:     1234,
	}

	// Expected
	eTile := coreModels.NewTile(api.PortTileType).WithMetrics(coreModels.MillisecondUnit)
	eTile.Label = fmt.Sprintf("%s:%d", param.Hostname, param.Port)
	eTile.Status = coreModels.SuccessStatus
	eTile.Message = "responding"
	eTile.Metrics.Values = []string{"50"}

	// Test
	rTile, err := usecase.Port(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "OpenSocket", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_CheckPort_Fail(t *testing.T) {
	// Init
	mockRepo := new(mocks.Repository)
	mockRepo.On("OpenSocket", AnythingOfType("string"), AnythingOfType("int"), AnythingOfType("string"), Anything).Return(false, "", time.Millisecond*0, errors.New("port error"))
	usecase := NewPortUsecase(mockRepo)

	// Params
	param := &models.PortParams{
		Hostname: "monitoror.example.com",
		Port:     1234,
	}

	// Expected
	eTile := coreModels.NewTile(api.PortTileType)
	eTile.Label = fmt.Sprintf("%s:%d", param.Hostname, param.Port)
	eTile.Status = coreModels.FailedStatus

	// Test
	rTile, err := usecase.Port(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "OpenSocket", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_CheckPort_WithPayload(t *testing.T) {
	mockRepo := new(mocks.Repository)
	payload := []byte{0xde, 0xad}
	mockRepo.On("OpenSocket", "monitoror.example.com", 1234, "udp", payload).Return(true, "pong", time.Millisecond*10, nil)
	usecase := NewPortUsecase(mockRepo)

	param := &models.PortParams{
		Hostname: "monitoror.example.com",
		Port:     1234,
		Type:     models.UDPPortType,
		Payload:  "dead",
	}

	eTile := coreModels.NewTile(api.PortTileType).WithMetrics(coreModels.MillisecondUnit)
	eTile.Label = fmt.Sprintf("%s:%d", param.Hostname, param.Port)
	eTile.Status = coreModels.SuccessStatus
	eTile.Message = "responding"
	eTile.Metrics.Values = []string{"10"}

	rTile, err := usecase.Port(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "OpenSocket", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_CheckPort_WithRawPayload(t *testing.T) {
	mockRepo := new(mocks.Repository)
	rawPayload := []byte{0x0a, 0x0d, 0x0a, 0x0d}
	mockRepo.On("OpenSocket", "monitoror.example.com", 22022, "tcp", rawPayload).Return(true, "pong", time.Millisecond*1, nil)
	usecase := NewPortUsecase(mockRepo)

	param := &models.PortParams{
		Hostname: "monitoror.example.com",
		Port:     22022,
		Payload:  string([]byte{0x0a, 0x0d, 0x0a, 0x0d}),
	}

	eTile := coreModels.NewTile(api.PortTileType).WithMetrics(coreModels.MillisecondUnit)
	eTile.Label = fmt.Sprintf("%s:%d", param.Hostname, param.Port)
	eTile.Status = coreModels.SuccessStatus
	eTile.Message = "responding"
	eTile.Metrics.Values = []string{"1"}

	rTile, err := usecase.Port(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "OpenSocket", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_CheckPort_Display(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("OpenSocket", "monitoror.example.com", 22, "tcp", Anything).Return(true, "SSH-2.0-OpenSSH_7.9", time.Millisecond*5, nil)
	usecase := NewPortUsecase(mockRepo)

	param := &models.PortParams{
		Hostname: "monitoror.example.com",
		Port:     22,
		Display:  `SSH-2.0-(.*)`,
	}

	eTile := coreModels.NewTile(api.PortTileType).WithMetrics(coreModels.MillisecondUnit)
	eTile.Label = "monitoror.example.com:22"
	eTile.Status = coreModels.SuccessStatus
	eTile.Message = "OpenSSH_7.9"
	eTile.Metrics.Values = []string{"5"}

	rTile, err := usecase.Port(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "OpenSocket", 1)
		mockRepo.AssertExpectations(t)
	}
}
