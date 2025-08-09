//go:build mediacore_grpc

package mediacore

import (
	"context"
	"fmt"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	mediacorev1 "github.com/hywgb/pjSipKF/proto/gen/go/mediacore/v1"
)

type grpcClient struct {
	cli mediacorev1.MediaCoreClient
}

func NewStubClient() Client {
	panic("build with -tags mediacore_grpc and use NewGRPCClient")
}

func normalizeUnixTarget(p string) string {
	if strings.HasPrefix(p, "unix:") {
		return p
	}
	// ensure absolute path and unix:/// prefix
	path := p
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "unix:///" + strings.TrimPrefix(path, "/")
}

func NewGRPCClientUDS(udsPath string) (Client, error) {
	target := normalizeUnixTarget(udsPath)
	cc, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &grpcClient{cli: mediacorev1.NewMediaCoreClient(cc)}, nil
}

// For tests without filesystem UDS
func NewGRPCClientBuf(conn *bufconn.Listener) (Client, error) {
	if conn == nil {
		return nil, fmt.Errorf("nil bufconn")
	}
	dialer := func(context.Context, string) (net.Conn, error) { return conn.Dial() }
	cc, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &grpcClient{cli: mediacorev1.NewMediaCoreClient(cc)}, nil
}

func (g *grpcClient) CreateSession(ctx context.Context, sdpOffer string, metadata map[string]string) (string, string, error) {
	resp, err := g.cli.CreateSession(ctx, &mediacorev1.CreateSessionRequest{SdpOffer: sdpOffer, Metadata: metadata})
	if err != nil {
		return "", "", err
	}
	return resp.SessionId, resp.SdpAnswer, nil
}

func (g *grpcClient) TerminateSession(ctx context.Context, sessionID string) error {
	_, err := g.cli.TerminateSession(ctx, &mediacorev1.TerminateSessionRequest{SessionId: sessionID})
	return err
}