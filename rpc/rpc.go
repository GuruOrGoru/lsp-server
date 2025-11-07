package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type BaseMessage struct {
	Method string `json:"method"`
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %v\r\n\r\n%s", len(content), content)
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("separator not found")
	}
	lines := strings.Split(string(header), "\r\n")
	var contentLength int
	var err error
	for _, line := range lines {
		if strings.HasPrefix(line, "Content-Length: ") {
			contentLengthStr := strings.TrimSpace(line[len("Content-Length: "):])
			contentLength, err = strconv.Atoi(contentLengthStr)
			if err != nil {
				return "", nil, err
			}
			break
		}
	}
	if contentLength == 0 {
		return "", nil, errors.New("Content-Length header not found")
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, err
	}

	return baseMessage.Method, content[:contentLength], nil
}

func Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}
	lines := strings.Split(string(header), "\r\n")
	var contentLength int
	for _, line := range lines {
		if strings.HasPrefix(line, "Content-Length: ") {
			contentLengthStr := strings.TrimSpace(line[len("Content-Length: "):])
			contentLength, err = strconv.Atoi(contentLengthStr)
			if err != nil {
				return 0, nil, err
			}
			break
		}
	}
	if contentLength == 0 {
		return 0, nil, errors.New("Content-Length header not found")
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	totalLengthOfIncomingMessage := len(header) + 4 + contentLength

	return totalLengthOfIncomingMessage, data[:totalLengthOfIncomingMessage], nil
}
