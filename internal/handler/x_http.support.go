package handler

import (
	"github.com/lwmacct/260630-go-hsr-shared/pkg/identity"

	"github.com/lwmacct/260605-miniport/internal/service"
)

type Config struct {
	Identity identity.SessionResolver
}

type Services struct {
	Portsvc *service.PortsvcService
}
