package utils

import (
	"fmt"
	"net"
)

func RPCAddr(bindAddr string, rpcPort int) (string, error) {
	host, _, err := net.SplitHostPort(bindAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, rpcPort), nil
}
