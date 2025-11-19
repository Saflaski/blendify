package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"

	"github.com/golang/glog"
)

type BlendService struct {
	BlendStore     *RedisStateStore
	LastFMExternal *musicapi.LastFMAPIExternal
}

type blendCategory string
type blendTimeDuration string

const (
	BlendCategoryArtist blendCategory = "artist"
	BlendCategoryTrack  blendCategory = "track"
	BlendCategoryAlbum  blendCategory = "album"
)

const (
	BlendTimeDurationOneMonth   blendTimeDuration = "1month"
	BlendTimeDurationThreeMonth blendTimeDuration = "3month"
	// BlendTimeDurationSixMonth   blendTimeDuration = "6month"
	BlendTimeDurationYear blendTimeDuration = "12month"
	// BlendTimeDurationAllTime    blendTimeDuration = "alltime"
)

func NewBlendService(blendStore RedisStateStore, lfmAdapter musicapi.LastFMAPIExternal) *BlendService {
	return &BlendService{&blendStore, &lfmAdapter}
}

func (s *BlendService) GetBlend(user string, category blendCategory, timeDuration blendTimeDuration) (int, error) {
	//Implement logic to calculate blend percentage based on user data, category, and time duration
	//For now, return mock data

	glog.Info("Calculating blend for user: ", user, " category: ", category, " timeDuration: ", timeDuration)
	userListenHistory, err := s.BlendStore.GetUserListenHistory(user)
	if err != nil {
		return 0, err
	}
	_ = userListenHistory // Placeholder to avoid unused variable error

	//Mock Changing Data
	if category == BlendCategoryArtist && timeDuration == BlendTimeDurationOneMonth {
		return 75, nil
	} else if category == BlendCategoryArtist && timeDuration == BlendTimeDurationThreeMonth {
		return 60, nil
	} else if category == BlendCategoryTrack && timeDuration == BlendTimeDurationOneMonth {
		return 85, nil
	}

	return 42, nil
}
