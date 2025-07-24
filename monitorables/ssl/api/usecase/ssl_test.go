package usecase

import (
	"errors"
	"testing"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ssl/api"
	"github.com/monitoror/monitoror/monitorables/ssl/api/mocks"
	"github.com/monitoror/monitoror/monitorables/ssl/api/models"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestUsecase_SSL_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	cert := &api.Certificate{NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(48 * time.Hour), Issuer: "issuer", Subject: "subject"}
	mockRepo.On("FetchCertificate", AnythingOfType("string"), AnythingOfType("int")).Return(cert, nil)

	usecase := NewSSLUsecase(mockRepo)

	param := &models.SSLParams{Domain: "example.com", Port: 443, WarnDays: 1}

	eTile := coreModels.NewTile(api.SSLTileType).WithMetrics(coreModels.RawUnit)
	eTile.Label = "example.com:443"
	eTile.Status = coreModels.SuccessStatus
	remaining := int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
	eTile.Message = buildMessage("full", cert, remaining)
	eTile.Metrics.Values = []string{cert.NotBefore.Format(time.RFC3339), cert.NotAfter.Format(time.RFC3339), "issuer", "subject"}

	rTile, err := usecase.SSL(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "FetchCertificate", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_SSL_Warn(t *testing.T) {
	mockRepo := new(mocks.Repository)
	cert := &api.Certificate{NotBefore: time.Now(), NotAfter: time.Now().Add(12 * time.Hour), Issuer: "issuer", Subject: "subject"}
	mockRepo.On("FetchCertificate", AnythingOfType("string"), AnythingOfType("int")).Return(cert, nil)

	usecase := NewSSLUsecase(mockRepo)

	param := &models.SSLParams{Domain: "example.com", Port: 443, WarnDays: 1}

	eTile := coreModels.NewTile(api.SSLTileType).WithMetrics(coreModels.RawUnit)
	eTile.Label = "example.com:443"
	eTile.Status = coreModels.WarningStatus
	remaining := int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
	eTile.Message = buildMessage("full", cert, remaining)
	eTile.Metrics.Values = []string{cert.NotBefore.Format(time.RFC3339), cert.NotAfter.Format(time.RFC3339), "issuer", "subject"}

	rTile, err := usecase.SSL(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "FetchCertificate", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_SSL_Fail(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("FetchCertificate", AnythingOfType("string"), AnythingOfType("int")).Return(nil, errors.New("fail"))

	usecase := NewSSLUsecase(mockRepo)

	param := &models.SSLParams{Domain: "example.com", Port: 443, WarnDays: 1}

	eTile := coreModels.NewTile(api.SSLTileType)
	eTile.Label = "example.com:443"
	eTile.Status = coreModels.FailedStatus
	eTile.Message = "fail"

	rTile, err := usecase.SSL(param)

	if assert.NoError(t, err) {
		assert.Equal(t, eTile, rTile)
		mockRepo.AssertNumberOfCalls(t, "FetchCertificate", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_SSL_CustomDisplay(t *testing.T) {
	mockRepo := new(mocks.Repository)
	cert := &api.Certificate{NotBefore: time.Now(), NotAfter: time.Now().Add(48 * time.Hour), Issuer: "issuer", Subject: "subject"}
	mockRepo.On("FetchCertificate", AnythingOfType("string"), AnythingOfType("int")).Return(cert, nil)

	usecase := NewSSLUsecase(mockRepo)

	param := &models.SSLParams{Domain: "example.com", Port: 443, WarnDays: 1, Display: "issuer"}

	tile, err := usecase.SSL(param)

	if assert.NoError(t, err) {
		assert.Equal(t, "issuer", tile.Message)
		mockRepo.AssertNumberOfCalls(t, "FetchCertificate", 1)
		mockRepo.AssertExpectations(t)
	}
}
