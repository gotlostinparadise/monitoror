//go:generate mockery --name Usecase

package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/rtsp/api/models"
)

const (
	RTSPTileType coreModels.TileType = "RTSP"
)

type (
	Usecase interface {
		RTSP(params *models.RTSPParams) (*coreModels.Tile, error)
	}
)
