package auth

import (
	"context"

	"github.com/STLeee/mediation-platform/backend/core/model"
)

// BaseAuth interface for authentication
type BaseAuth interface {
	AuthenticateByToken(ctx context.Context, token string) (uid string, err error)
	GetUseInfo(ctx context.Context, uid string) (*model.UserInfo, error)
}
