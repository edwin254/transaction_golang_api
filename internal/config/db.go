package config

type DBConfig struct {
	POSTGRES_HOST     string `env:"POSTGRES_HOST,required"`
	POSTGRES_USER     string `env:"POSTGRES_USER,required"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD,required"`
	POSTGRES_DB       string `env:"POSTGRES_DB,required"`
}
