package server_test

import (
	"testing"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/260605-miniport/internal/appcmd/server"
	"github.com/lwmacct/260605-miniport/internal/config"
)

func TestServerCommandCoversConfigFlags(t *testing.T) {
	cfgm.AssertCommandFlagCoverage(
		t,
		server.Command,
		config.DefaultConfig(),
		[]string{"server"},
		cfgm.IgnoreConfigKeys("server.database.pgsql.password"),
	)
}
