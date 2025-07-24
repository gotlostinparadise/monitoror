package usecase

import (
    "errors"
    "testing"

    coreModels "github.com/monitoror/monitoror/models"
    "github.com/monitoror/monitoror/monitorables/dns/api"
    "github.com/monitoror/monitoror/monitorables/dns/api/mocks"
    "github.com/monitoror/monitoror/monitorables/dns/api/models"

    "github.com/stretchr/testify/assert"
    . "github.com/stretchr/testify/mock"
)

func TestUsecase_DNS_Success(t *testing.T) {
    mockRepo := new(mocks.Repository)
    mockRepo.On("Lookup", AnythingOfType("string"), AnythingOfType("string")).Return([]string{"1.2.3.4"}, nil)
    usecase := NewDNSUsecase(mockRepo)

    params := &models.DNSParams{RecordType: "A", Name: "example.com", ExpectedValue: "1.2.3.4"}

    expected := coreModels.NewTile(api.DNSTileType)
    expected.Label = "A example.com"
    expected.Status = coreModels.SuccessStatus

    tile, err := usecase.DNS(params)
    if assert.NoError(t, err) {
        assert.Equal(t, expected, tile)
        mockRepo.AssertNumberOfCalls(t, "Lookup", 1)
    }
}

func TestUsecase_DNS_Fail(t *testing.T) {
    mockRepo := new(mocks.Repository)
    mockRepo.On("Lookup", AnythingOfType("string"), AnythingOfType("string")).Return(nil, errors.New("dns error"))
    usecase := NewDNSUsecase(mockRepo)

    params := &models.DNSParams{RecordType: "A", Name: "example.com", ExpectedValue: "1.2.3.4"}

    expected := coreModels.NewTile(api.DNSTileType)
    expected.Label = "A example.com"
    expected.Status = coreModels.FailedStatus

    tile, err := usecase.DNS(params)
    if assert.NoError(t, err) {
        assert.Equal(t, expected, tile)
        mockRepo.AssertNumberOfCalls(t, "Lookup", 1)
    }
}
