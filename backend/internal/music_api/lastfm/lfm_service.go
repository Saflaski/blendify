package musicapi

type LastFMServiceConfig struct {
	ApiKey    string
	LastFMURL string
}

type LastFMAPIService struct {
	serviceConfig LastFMServiceConfig
}

func NewLastFMService(serviceConfig LastFMServiceConfig) *LastFMAPIService {
	return &LastFMAPIService{
		serviceConfig,
	}
}
