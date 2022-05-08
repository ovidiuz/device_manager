package domain

import (
	"time"

	"github.com/Netflix/go-env"
)

type ServiceConfig struct {
	ServicePort     int    `env:"SERVICE_PORT,default=8080"`
	CasbinTableName string `env:"CASBIN_TABLE_NAME,default=casbin_rule"`
	CasbinModelFile string `env:"CASBIN_MODEL_FILE,default=config/rbac_model.conf"`

	JwtTTL       time.Duration `env:"JWT_TTL,default=2h"`
	HTTPSEnabled bool          `env:"HTTPS_ENABLED,default=true"`
	PostgreSQLConfig
}

type PostgreSQLConfig struct {
	PSQLHost     string `env:"PSQL_HOST,default=localhost"`
	PSQLUsername string `env:"PSQL_USERNAME,default=admin"`
	PSQLPassword string `env:"PSQL_PASSWORD,default=admin"`
	PSQLDBName   string `env:"PSQL_DB_NAME,default=db1"`
}

func InitConfig() (*ServiceConfig, error) {
	cfg := &ServiceConfig{}
	_, err := env.UnmarshalFromEnviron(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
