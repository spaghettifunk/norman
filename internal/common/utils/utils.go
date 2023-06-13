package utils

import (
	"fmt"
	"net"
	"strings"
	"unicode"
)

func RPCAddr(bindAddr string, rpcPort int) (string, error) {
	host, _, err := net.SplitHostPort(bindAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, rpcPort), nil
}

func RemovePunctuations(input string) string {
	var builder strings.Builder
	for _, char := range input {
		if !unicode.IsPunct(char) {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}
