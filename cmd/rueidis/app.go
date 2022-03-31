package main

import (
	"github.com/rs/zerolog/log"
	"github.com/ujunglangit-id/redis-timeseries-bench/pkg/model/config"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal().Msgf("failed load config %#v", err)
		return
	}

}
