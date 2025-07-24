package dns

import (
    "testing"

    "github.com/monitoror/monitoror/internal/pkg/monitorable/test"
    "github.com/stretchr/testify/assert"
)

func TestNewMonitorable(t *testing.T) {
    store, helper := test.InitMockAndStore()

    monitorable := NewMonitorable(store)
    assert.NotNil(t, monitorable)
    assert.NotNil(t, monitorable.GetDisplayName())

    for _, variant := range monitorable.GetVariantsNames() {
        if valid, _ := monitorable.Validate(variant); valid {
            monitorable.Enable(variant)
        }
    }

    helper.RouterAssertNumberOfCalls(t, 1, 1)
    helper.TileSettingsManagerAssertNumberOfCalls(t, 1, 0, 1, 0)
}
