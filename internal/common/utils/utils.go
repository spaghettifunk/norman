package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"unicode"

	"github.com/andybalholm/brotli"
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

func CompressBrotli(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := brotli.NewWriterLevel(&b, brotli.BestCompression)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func DecompressBrotli(data []byte) ([]byte, error) {
	br := bytes.NewReader(data)
	decompressor := brotli.NewReader(br)
	return ioutil.ReadAll(decompressor)
}
