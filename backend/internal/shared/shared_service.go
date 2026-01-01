package shared

import (
	"backend-lastfm/internal/auth"
	blend "backend-lastfm/internal/blending"
	"context"
	"fmt"

	"github.com/golang/glog"
)

type SharedService struct {
	authService  *auth.AuthService
	blendService *blend.BlendService
}

func (s *SharedService) DeleteAllUserData(context context.Context, userid string) error {
	//Have to delete all blends before deleting from auth related places in db
	err := s.blendService.DeleteUserBlends(context, userid)
	if err != nil {
		glog.Error("COULD NOT DELETE USER BLENDS")
		return fmt.Errorf(": %w", err)
	}
	err = s.authService.DeleteUser(context, userid)
	if err != nil {
		glog.Error("COULD NOT DELETE USER AUTH LEVEL")
		return fmt.Errorf("Could not delete auth side user data: %w", err)
	}

	return nil
}

func NewSharedService(auth *auth.AuthService, blend *blend.BlendService) *SharedService {
	return &SharedService{
		authService:  auth,
		blendService: blend,
	}
}
