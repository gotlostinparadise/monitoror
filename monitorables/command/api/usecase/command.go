package usecase

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/command/api"
	"github.com/monitoror/monitoror/monitorables/command/api/models"
)

type commandUsecase struct {
	repository api.Repository
}

func NewCommandUsecase(r api.Repository) api.Usecase { return &commandUsecase{r} }

func (cu *commandUsecase) CommandStatus(params *models.CommandParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.CommandTileType)
	tile.Label = params.Command

	output, exitCode, duration, err := cu.repository.Exec(params.Command)
	if err != nil {
		tile.Status = coreModels.FailedStatus
		tile.Message = err.Error()
		return tile, nil
	}

	if exitCode <= params.GetExitCodeMax() {
		tile.Status = coreModels.SuccessStatus
	} else {
		tile.Status = coreModels.FailedStatus
	}

	cleanedOutput := strings.TrimSpace(output)

	if params.Display != "" {
		if re, err := regexp.Compile(params.Display); err == nil {
			if m := re.FindStringSubmatch(cleanedOutput); m != nil {
				if len(m) > 1 {
					tile.Message = m[1]
				} else {
					tile.Message = m[0]
				}
			}
		}
	}

	if params.Metrics != "" {
		switch params.Metrics {
		case "duration":
			tile.WithMetrics(coreModels.MillisecondUnit)
			tile.Metrics.Values = []string{fmt.Sprintf("%d", duration.Milliseconds())}
		case "exitCode":
			tile.WithMetrics(coreModels.NumberUnit)
			tile.Metrics.Values = []string{strconv.Itoa(exitCode)}
		default:
			if re, err := regexp.Compile(params.Metrics); err == nil {
				if m := re.FindStringSubmatch(cleanedOutput); m != nil {
					value := m[0]
					if len(m) > 1 {
						value = m[1]
					}
					unit := coreModels.RawUnit
					if _, err := strconv.ParseFloat(value, 64); err == nil {
						unit = coreModels.NumberUnit
					}
					tile.WithMetrics(unit)
					tile.Metrics.Values = []string{value}
				}
			}
		}
	}

	return tile, nil
}
