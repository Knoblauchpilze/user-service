package server

import (
	"net"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUnit_DeterminePort_Tcp(t *testing.T) {
	address := &net.TCPAddr{
		Port: 3251,
	}

	port, err := determinePort(address)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(uint16(3251), port)
}

func TestUnit_DeterminePort_Tcp_WhenPortTooBig_ExpectsTruncated(t *testing.T) {
	address := &net.TCPAddr{
		Port: 99999,
	}

	port, err := determinePort(address)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(uint16(34463), port)
}

func TestUnit_DeterminePort_UnsupportedProtocol(t *testing.T) {
	address := &net.UDPAddr{
		Port: 3251,
	}

	_, err := determinePort(address)

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, UnsupportedProtocol), "Actual err: %v", err)
}
