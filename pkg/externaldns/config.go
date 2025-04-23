package externaldns

import "github.com/Netflix/go-env"

type Config struct {
	Gravity struct {
		URL   string `env:"GRAVITY_URL"`
		Token string `env:"GRAVITY_TOKEN"`
	}
	DomainFilter []string `env:"DOMAIN_FILTER"`
	Listen       struct {
		// https://kubernetes-sigs.github.io/external-dns/v0.14.2/tutorials/webhook-provider/
		API     string `env:"LISTEN_API,default=localhost:8888"`
		Metrics string `env:"LISTEN_METRICS,default=0.0.0.0:8080"`
	}
}

var globalConfig *Config

func Get() *Config {
	if globalConfig != nil {
		return globalConfig
	}
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		panic(err)
	}
	globalConfig = &cfg
	return &cfg
}
