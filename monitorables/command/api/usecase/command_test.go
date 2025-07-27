package usecase

import (
	"errors"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/command/api"
	"github.com/monitoror/monitoror/monitorables/command/api/mocks"
	"github.com/monitoror/monitoror/monitorables/command/api/models"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestUsecase_CommandStatus_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("Exec", AnythingOfType("string")).Return("ok", 0, 0, nil)
	usecase := NewCommandUsecase(mockRepo)

	param := &models.CommandParams{Command: "true"}

	eTile := coreModels.NewTile(api.CommandTileType)
	eTile.Label = param.Command
	eTile.Status = coreModels.SuccessStatus
	eTile.Message = "ok"

	rTile, err := usecase.CommandStatus(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "Exec", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_CommandStatus_Fail(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("Exec", AnythingOfType("string")).Return("boom", 1, 0, nil)
	usecase := NewCommandUsecase(mockRepo)

	param := &models.CommandParams{Command: "false"}

	eTile := coreModels.NewTile(api.CommandTileType)
	eTile.Label = param.Command
	eTile.Status = coreModels.FailedStatus
	eTile.Message = "boom"

	rTile, err := usecase.CommandStatus(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "Exec", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_CommandStatus_Error(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("Exec", AnythingOfType("string")).Return("", 0, 0, errors.New("failed"))
	usecase := NewCommandUsecase(mockRepo)

	param := &models.CommandParams{Command: "bad"}

	eTile := coreModels.NewTile(api.CommandTileType)
	eTile.Label = param.Command
	eTile.Status = coreModels.FailedStatus
	eTile.Message = "failed"

	rTile, err := usecase.CommandStatus(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "Exec", 1)
		mockRepo.AssertExpectations(t)
	}
}
