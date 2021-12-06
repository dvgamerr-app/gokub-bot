package bitkub

var (
	UrlAPI = "https://api.bitkub.com"
)

type Config struct {
	ApiKey    string
	SecretKey string
}

func (cfg *Config) Init() {
	if cfg.ApiKey == "" {
		panic("'ApiKey' is empty.")
	}
	if cfg.SecretKey == "" {
		panic("'SecretKey' is empty.")
	}
}
