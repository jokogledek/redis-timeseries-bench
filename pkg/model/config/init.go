package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

func InitConfig() (cfg *Config, err error) {
	log.Infof("config load [OK]")
	cfg = &Config{}
	err = cfg.readFile()
	if err != nil {
		return
	}

	return
}

func (cfg *Config) readFile() (err error) {
	var (
		f *os.File
	)

	cfg.Environment = os.Getenv(enum.APP_ENV)
	if cfg.Environment != "staging" && cfg.Environment != "production" {
		cfg.Environment = "development"
	}
	log.Infof("[cfg.Environment ADALAH] %s", cfg.Environment)

	path := []string{
		"/config/api-config",
		"files/config/api-config",
		"./files/config/api-config",
		"../files/config/api-config",
		"../../files/config/api-config",
	}

	for _, val := range path {
		f, err = os.Open(fmt.Sprintf(`%s/%s/config.main.yml`, val, cfg.Environment))
		if err == nil {
			log.Infof("[config][init] load config file from %s", fmt.Sprintf(`%s/%s/config.main.yml`, val, cfg.Environment))
			decoder := yaml.NewDecoder(f)
			err = decoder.Decode(cfg)
			break
		}
	}

	if err != nil {
		return
	}
	cfg.setLocalAddress()
	log.Infof("[config][ReadConfig] Config load success, running on \"%s\".", cfg.Environment)
	return
}
