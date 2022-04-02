package rueidis

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/rueian/rueidis"
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
	RedisClient    rueidis.Client
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
	b.RedisClient, err = rueidis.NewClient(rueidis.ClientOption{
		BlockingPoolSize:  5000,
		CacheSizeEachConn: 256 * (1 << 20),
		InitAddress:       []string{redisHost},
	})
	if err != nil {
		return
	}

	//prepare timeseries key
	b.Key = b.cfg.Redis.Key + ":rueidis_01"
	log.Info().Msgf("delete redis key  %s if exists", b.Key)
	//delete key if exist
	err = b.RedisClient.Do(context.Background(), b.RedisClient.B().Del().Key(b.Key).Build()).Error()
	if err != nil {
		log.Info().Msgf("failed to delete timeseries key %s, %#v", b.Key, err)
		err = nil
	}

	err = b.RedisClient.Do(context.Background(),
		b.RedisClient.B().TsCreate().
			Key(b.Key).
			Labels().Labels("volume", "0").
			Build()).Error()
	if err != nil {
		log.Fatal().Msgf("failed to create timeseries key %#v", err)
		return
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

func (b *Bench) RedisInsert() {
	if err := b.RedisClient.Dedicated(func(client rueidis.DedicatedClient) error {
		cmds := make(rueidis.Commands, 0, len(b.PayloadData)/2)
		for i := 0; i < len(b.PayloadData); i += 2 {
			cmd := b.RedisClient.B().TsMadd().KeyTimestampValue()
			cmd = cmd.KeyTimestampValue(b.Key, time.Now().UnixNano(), b.PayloadData[i+0].Volume)
			cmd = cmd.KeyTimestampValue(b.Key, time.Now().UnixNano(), b.PayloadData[i+1].Volume)
			cmds = append(cmds, cmd.Build())
		}
		for _, resp := range client.DoMulti(context.Background(), cmds...) {
			if err := resp.Error(); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.Error().Msgf("error add timeseries %#v", err)
	}
}

func (b *Bench) InsertToRedis() {
	var (
		wg  = sync.WaitGroup{}
		mtx = sync.RWMutex{}
	)

	counter := 0
	poolData := b.RedisClient.B().TsMadd().KeyTimestampValue()
	for k, v := range b.PayloadData {
		counter++
		poolData.KeyTimestampValue(b.Key, time.Now().UnixNano(), v.Volume)
		if counter > b.cfg.Redis.MaxQueue || k == len(b.PayloadData)-1 {
			counter = 0
			wg.Add(1)
			buffPool := poolData
			go func() {
				defer wg.Done()
				mtx.Lock()
				err := b.RedisClient.Do(context.Background(), buffPool.Build()).Error()
				if err != nil {
					log.Error().Msgf("error add timeseries %#v", err)
				}
				mtx.Unlock()
			}()
			poolData = b.RedisClient.B().TsMadd().KeyTimestampValue()
		}
	}
	wg.Wait()
}

func (b *Bench) FetchFromRedis() {

}
