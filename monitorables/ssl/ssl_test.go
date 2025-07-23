package ssl

import (
	"os"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitorable(t *testing.T) {
	store, mockHelper := test.InitMockAndStore()

	_ = os.Setenv("MO_MONITORABLE_SSL_VARIANT0_TIMEOUT", "-1000")

	monitorable := NewMonitorable(store)
	assert.NotNil(t, monitorable)

	assert.NotNil(t, monitorable.GetDisplayName())

	if assert.Len(t, monitorable.GetVariantsNames(), 2) {
		_, errors := monitorable.Validate("variant0")
		assert.NotEmpty(t, errors)
	}

	for _, variantName := range monitorable.GetVariantsNames() {
		if valid, _ := monitorable.Validate(variantName); valid {
			monitorable.Enable(variantName)
		}
	}

	mockHelper.RouterAssertNumberOfCalls(t, 1, 1)
	mockHelper.TileSettingsManagerAssertNumberOfCalls(t, 1, 0, 1, 0)
}
