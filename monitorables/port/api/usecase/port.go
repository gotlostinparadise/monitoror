//go:build !faker

package usecase

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api"
	"github.com/monitoror/monitoror/monitorables/port/api/models"
)

type (
	portUsecase struct {
		repository api.Repository
	}
)

func NewPortUsecase(repository api.Repository) api.Usecase {
	return &portUsecase{repository}
}

func (pu *portUsecase) Port(params *models.PortParams) (tile *coreModels.Tile, err error) {
	tile = coreModels.NewTile(api.PortTileType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	var payload []byte
	if params.Payload != "" {
		if matched, _ := regexp.MatchString(`^0x[0-9a-fA-F]+$`, params.Payload); matched {
			payload, err = hex.DecodeString(params.Payload[2:])
		} else {
			s := params.Payload
			if u, e := strconv.Unquote("\"" + s + "\""); e == nil {
				s = u
			}
			payload = []byte(s)
		}
		if err != nil {
			tile.Status = coreModels.FailedStatus
			err = nil
			return
		}
	}

	responding, banner, duration, err := pu.repository.OpenSocket(params.Hostname, params.Port, string(params.GetType()), payload)
	if err != nil {
		tile.Status = coreModels.FailedStatus
		err = nil
		return
	}

	tile.Status = coreModels.SuccessStatus
	tile.Message = "no response"
	if responding {
		tile.Message = "responding"
	}

	if params.Display != "" && banner != "" {
		if re, e := regexp.Compile(params.Display); e == nil {
			if m := re.FindStringSubmatch(banner); m != nil {
				if len(m) > 1 {
					tile.Message = m[1]
				} else {
					tile.Message = m[0]
				}
			}
		}
	}

	tile.WithMetrics(coreModels.MillisecondUnit)
	tile.Metrics.Values = []string{fmt.Sprintf("%d", duration.Milliseconds())}

	return
}
