package goredis

import (
	"encoding/csv"
	"fmt"
	redis "github.com/RedisTimeSeries/redistimeseries-go"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
	"github.com/ujunglangit-id/redis-timeseries-bench/pkg/model/config"
	"github.com/ujunglangit-id/redis-timeseries-bench/pkg/model/types"
	"os"
	"strconv"
	"sync"
	"time"
)

type Bench struct {
	cfg            *config.Config
	PayloadData    []*types.Stocks
	RedisClient    *redis.Client
	Key            string
	PipelineBuffer chan int
}

func NewBench(cfg *config.Config) *Bench {
	b := &Bench{
		cfg: cfg,
	}

	err := b.InitRedis()
	if err != nil {
		log.Fatal().Msgf("failed to init redis, %#v", err)
	}

	b.PipelineBuffer = make(chan int, cfg.Redis.Buffer)
	return b
}

func (b *Bench) InitRedis() (err error) {
	redisHost := fmt.Sprintf("%s:%d", b.cfg.Redis.Host, b.cfg.Redis.Port)
	pool := &redigo.Pool{
		MaxActive: 0,
		MaxIdle:   10,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", redisHost)
		}}

	b.RedisClient = redis.NewClientFromPool(pool, "timeseries")
	//prepare timeseries key
	b.Key = b.cfg.Redis.Key + ":goseries_01"
	log.Info().Msgf("redis key : %s", b.Key)
	_, exists := b.RedisClient.Info(b.Key)
	if exists == nil {
		err = b.RedisClient.DeleteSerie(b.Key)
		if err != nil {
			log.Fatal().Msgf("failed to delete timeseries key %#v", err)
			return err
		}
	}

	opt := redis.DefaultCreateOptions
	opt.Labels["volume"] = "0"
	err = b.RedisClient.CreateKeyWithOptions(b.Key, opt)
	if err != nil {
		log.Fatal().Msgf("failed to create timeseries key %#v", err)
		return err
	}
	return
}

func (b *Bench) LoadCsv() {
	// open file
	f, err := os.Open(b.cfg.Files.Csv)
	if err != nil {
		log.Error().Msgf("failed to open csv %#v", err)
		return
	}
	defer f.Close()

	//read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Error().Msgf("failed to read csv %#v", err)
		return
	}

	for i, line := range data {
		if i > 0 {
			vol, _ := strconv.ParseFloat(line[6], 64)
			b.PayloadData = append(b.PayloadData, &types.Stocks{
				Ticker: line[0],
				Day:    line[1],
				Open:   line[2],
				High:   line[3],
				Low:    line[4],
				Close:  line[5],
				Volume: vol,
			})
		}
	}
	log.Info().Msgf("data length : %d", len(b.PayloadData))
}

func (b *Bench) InsertToRedis() {
	var (
		wg  = sync.WaitGroup{}
		mtx = sync.RWMutex{}
	)

	counter := 0
	datapoints := []redis.Sample{}

	for k, v := range b.PayloadData {
		counter++
		datapoints = append(datapoints, redis.Sample{
			Key: b.Key,
			DataPoint: redis.DataPoint{
				Timestamp: time.Now().UnixNano(),
				Value:     v.Volume,
			},
		})

		if counter > b.cfg.Redis.MaxQueue || k == len(b.PayloadData)-1 {
			counter = 0
			wg.Add(1)
			//buffPool := poolData
			go func(sample []redis.Sample) {
				defer wg.Done()
				mtx.Lock()
				_, err := b.RedisClient.MultiAdd(sample...)
				if err != nil {
					log.Error().Msgf("error add timeseries %#v", err)
				}
				mtx.Unlock()
			}(datapoints)
			datapoints = []redis.Sample{}
		}
	}
	wg.Wait()
}
