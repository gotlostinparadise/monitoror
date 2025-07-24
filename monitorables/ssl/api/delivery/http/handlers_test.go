package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ssl/api"
	"github.com/monitoror/monitoror/monitorables/ssl/api/mocks"
	"github.com/monitoror/monitoror/monitorables/ssl/api/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/info", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	ctx.QueryParams().Set("domain", "monitoror.example.com")
	ctx.QueryParams().Set("port", "443")
	ctx.QueryParams().Set("warnDays", "10")

	return
}

func missingParam(t *testing.T, param string) {
	ctx, _ := initEcho()
	ctx.QueryParams().Del(param)
	mockUsecase := new(mocks.Usecase)
	handler := NewSSLDelivery(mockUsecase)
	err := handler.GetSSL(ctx)
	assert.Error(t, err)
	assert.IsType(t, &coreModels.MonitororError{}, err)
}

func TestDelivery_SSLHandler_Success(t *testing.T) {
	ctx, res := initEcho()

	tile := coreModels.NewTile(api.SSLTileType)
	tile.Label = "monitoror.example.com:443"
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("SSL", &models.SSLParams{Domain: "monitoror.example.com", Port: 443, WarnDays: 10}).Return(tile, nil)
	handler := NewSSLDelivery(mockUsecase)

	jsonTile, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	if assert.NoError(t, handler.GetSSL(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(jsonTile), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "SSL", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_SSLHandler_QueryParamsError_MissingDomain(t *testing.T) {
	missingParam(t, "domain")
}

func TestDelivery_SSLHandler_QueryParamsError_MissingWarnDays(t *testing.T) {
	missingParam(t, "warnDays")
}

func TestDelivery_SSLHandler_Error(t *testing.T) {
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("SSL", Anything).Return(nil, errors.New("ssl error"))
	handler := NewSSLDelivery(mockUsecase)

	assert.Error(t, handler.GetSSL(ctx))
	mockUsecase.AssertNumberOfCalls(t, "SSL", 1)
	mockUsecase.AssertExpectations(t)
}
