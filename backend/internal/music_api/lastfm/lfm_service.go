package musicapi

type LastFMServiceConfig struct {
	ApiKey    ApiKey
	LastFMURL LastFMURL
}

type LastFMAPIService struct {
	serviceConfig LastFMServiceConfig
}

func NewLastFMService(serviceConfig LastFMServiceConfig) *LastFMAPIService {
	return &LastFMAPIService{
		serviceConfig,
	}
}
