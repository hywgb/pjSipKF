//go:build !mediacore_grpc

package mediacore

import (
	"github.com/hywgb/pjSipKF/control-plane/internal/config"
)

func NewClientFromConfig(cfg config.Config) (Client, error) {
	_ = cfg
	return NewStubClient(), nil
}