package server

import (
	"net"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
)

func determinePort(address net.Addr) (uint16, error) {
	if tcpAddress, ok := address.(*net.TCPAddr); ok {
		return uint16(tcpAddress.Port), nil
	}

	return 0, errors.NewCode(UnsupportedProtocol)
}
