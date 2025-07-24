package usecase

import (
	"errors"
	"fmt"
	"testing"

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
	mockRepo.On("OpenSocket", AnythingOfType("string"), AnythingOfType("int"), AnythingOfType("string"), Anything).Return(nil)
	usecase := NewPortUsecase(mockRepo)

	// Params
	param := &models.PortParams{
		Hostname: "monitoror.example.com",
		Port:     1234,
	}

	// Expected
	eTile := coreModels.NewTile(api.PortTileType)
	eTile.Label = fmt.Sprintf("%s:%d", param.Hostname, param.Port)
	eTile.Status = coreModels.SuccessStatus

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
	mockRepo.On("OpenSocket", AnythingOfType("string"), AnythingOfType("int"), AnythingOfType("string"), Anything).Return(errors.New("port error"))
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
	mockRepo.On("OpenSocket", "monitoror.example.com", 1234, "udp", payload).Return(nil)
	usecase := NewPortUsecase(mockRepo)

	param := &models.PortParams{
		Hostname: "monitoror.example.com",
		Port:     1234,
		Type:     models.UDPPortType,
		Payload:  "dead",
	}

	eTile := coreModels.NewTile(api.PortTileType)
	eTile.Label = fmt.Sprintf("%s:%d", param.Hostname, param.Port)
	eTile.Status = coreModels.SuccessStatus

	rTile, err := usecase.Port(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "OpenSocket", 1)
		mockRepo.AssertExpectations(t)
	}
}
