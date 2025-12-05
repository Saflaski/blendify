package blend

import musicapi "backend-lastfm/internal/music_api/lastfm"

type blendId string
type blendLinkValue string

type userid string

type blendCategory string
type blendTimeDuration string

var durationMap = map[blendTimeDuration]musicapi.Period{
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

var durationRange = []blendTimeDuration{
	BlendTimeDurationOneMonth,
	BlendTimeDurationThreeMonth,
	BlendTimeDurationYear,
}

var categoryRange = []blendCategory{
	BlendCategoryAlbum,
	BlendCategoryArtist,
	BlendCategoryTrack,
}
