package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	JwtSecretKey string `toml:"jwt_secret_key"`
	MongoURI     string `toml:"mongo_uri"`
	RedisURI     string `toml:"redis_uri"`
}

func Load(path string) (*Config, error) {
	var config Config

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
