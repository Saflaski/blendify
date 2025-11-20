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

var durationMap = map[blendTimeDuration]musicapi.Period{ //TODO Delete this shit
	BlendTimeDurationOneMonth:   musicapi.ONE_MONTH,
	BlendTimeDurationThreeMonth: musicapi.THREE_MONTHS,
	BlendTimeDurationYear:       musicapi.YEAR,
}

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

	//For now import from lastfm

	//Mock Changing Data

	return 42, nil
}

func (s *BlendService) getArists(userName string, timeDuration blendTimeDuration) (map[string]int, error) {
	topArtist, err := s.LastFMExternal.GetUserTopArtists(
		userName,
		durationMap[timeDuration],
		1,
		50,
	)
	glog.Info(topArtist, err)

	//TODO finish this function

	return make(map[string]int), nil
}
