package domain

import "github.com/Netflix/go-env"

type ServiceConfig struct {
	InfluxConfig
}

type InfluxConfig struct {
	InfluxHost  string `env:"INFLUX_HOST,default=http://127.0.0.1:8086"`
	InfluxToken string `env:"INFLUX_TOKEN,default=my-token"`
	InfluxOrg   string `env:"INFLUX_ORG,default=my-org"`
}

func InitConfig() (*ServiceConfig, error) {
	cfg := &ServiceConfig{}
	_, err := env.UnmarshalFromEnviron(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
