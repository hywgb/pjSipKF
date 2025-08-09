package mediacorev1

import (
	"context"
)

type CreateSessionRequest struct{
	SdpOffer string
	Metadata map[string]string
}

type CreateSessionResponse struct{
	SessionId string
	SdpAnswer string
}

type UpdateSessionRequest struct{
	SessionId string
	SdpOffer string
}

type UpdateSessionResponse struct{
	SdpAnswer string
}

type TerminateSessionRequest struct{
	SessionId string
}

type TerminateSessionResponse struct{ Ok bool }

type MediaCoreClient interface{
	CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...any) (*CreateSessionResponse, error)
	UpdateSession(ctx context.Context, in *UpdateSessionRequest, opts ...any) (*UpdateSessionResponse, error)
	TerminateSession(ctx context.Context, in *TerminateSessionRequest, opts ...any) (*TerminateSessionResponse, error)
}

type UnimplementedMediaCoreClient struct{}

func (UnimplementedMediaCoreClient) CreateSession(context.Context, *CreateSessionRequest, ...any) (*CreateSessionResponse, error){ return nil, nil }
func (UnimplementedMediaCoreClient) UpdateSession(context.Context, *UpdateSessionRequest, ...any) (*UpdateSessionResponse, error){ return nil, nil }
func (UnimplementedMediaCoreClient) TerminateSession(context.Context, *TerminateSessionRequest, ...any) (*TerminateSessionResponse, error){ return nil, nil }

func NewMediaCoreClient(_ any) MediaCoreClient { return UnimplementedMediaCoreClient{} }