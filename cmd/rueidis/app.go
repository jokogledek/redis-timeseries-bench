package main

import (
	"github.com/rs/zerolog/log"
	"github.com/ujunglangit-id/redis-timeseries-bench/pkg/bench/rueidis"
	"github.com/ujunglangit-id/redis-timeseries-bench/pkg/model/config"
	"time"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal().Msgf("failed load config %#v", err)
		return
	}

	start := time.Now()
	bench := rueidis.NewBench(cfg)
	bench.LoadCsvToRedis()
	log.Printf("load to redis finished in %s", time.Since(start))

	bench.FetchFromRedis()
	log.Printf("load from redis finished  in %s", time.Since(start))
}
