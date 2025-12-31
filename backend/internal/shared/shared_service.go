package shared

import (
	"backend-lastfm/internal/auth"
	blend "backend-lastfm/internal/blending"
	"context"
)

type SharedService struct {
	authService  *auth.AuthService
	blendService *blend.BlendService
}

func (s *SharedService) DeleteAllUserData(context context.Context, userid string) error {
	return nil
}

func NewSharedService(auth *auth.AuthService, blend *blend.BlendService) *SharedService {
	return &SharedService{
		authService:  auth,
		blendService: blend,
	}
}
