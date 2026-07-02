package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizePortSlotUsesPortPrefixBounds(t *testing.T) {
	group := PortGroup{ID: "group-id", PortPrefix: 1031}

	slot, err := utilNormalizePortSlot(group, PortSlotPayload{Port: 10310, Name: "api"})
	require.NoError(t, err)
	require.Equal(t, 10310, slot.Port)

	slot, err = utilNormalizePortSlot(group, PortSlotPayload{Port: 10319, Name: "worker"})
	require.NoError(t, err)
	require.Equal(t, 10319, slot.Port)

	_, err = utilNormalizePortSlot(group, PortSlotPayload{Port: 10320, Name: "overflow"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "slot port must be inside the port group")
}
