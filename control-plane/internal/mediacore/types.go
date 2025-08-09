package mediacore

import "context"

type Client interface {
	CreateSession(ctx context.Context, sdpOffer string, metadata map[string]string) (sessionID string, sdpAnswer string, err error)
	TerminateSession(ctx context.Context, sessionID string) error
}