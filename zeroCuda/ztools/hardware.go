package ztools

import (
	"errors"
	"net"
)

func GetFirstMacAddr() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	if len(interfaces) == 0 {
		return "", errors.New("no physical address")
	}
	addr := interfaces[0].HardwareAddr
	return addr.String(), nil
}
