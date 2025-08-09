//go:build !mediacore_grpc

package mediacore

import (
	"context"
	"fmt"
)

type stubClient struct{}

func NewStubClient() Client { return &stubClient{} }

func (s *stubClient) CreateSession(_ context.Context, sdpOffer string, metadata map[string]string) (string, string, error) {
	_ = metadata
	return "sess-0000001", fmt.Sprintf("v=0\n; answer for: %s", sdpOffer), nil
}

func (s *stubClient) TerminateSession(_ context.Context, sessionID string) error {
	_ = sessionID
	return nil
}