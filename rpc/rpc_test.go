package rpc_test

import (
	"testing"

	"github.com/guruorgoru/lsp-server/rpc"
)

type EncodingExample struct {
	Testing bool
}

func TestEncodeMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  any
		want string
	}{
		{
			"Test1",
			EncodingExample{Testing: true},
			"Content-Length: 16\r\n\r\n{\"Testing\":true}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rpc.EncodeMessage(tt.msg)
			if got != tt.want {
				t.Errorf("EncodeMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeMessage(t *testing.T) {
	tests := []struct {
		name          string
		msg           []byte
		method        string
		contentLength int
		wantErr       bool
	}{
		{
			"Test1",
			[]byte("Content-Length: 23\r\n\r\n{\"method\":\"completion\"}"),
			"completion",
			23,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, content, gotErr := rpc.DecodeMessage(tt.msg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DecodeMessage() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DecodeMessage() succeeded unexpectedly")
			}
			if method != tt.method {
				t.Errorf("DecodeMessage() = %v, want %v", method, tt.method)
			}

			contentLength := len(content)
			if contentLength != tt.contentLength {
				t.Errorf("DecodeMessage() = %v, want %v", content, tt.contentLength)
			}
		})
	}
}
