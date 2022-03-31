package rueidis

import (
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"github.com/ujunglangit-id/redis-timeseries-bench/pkg/model/config"
	"os"
)

type Bench struct {
	cfg *config.Config
}

func NewBench(cfg *config.Config) *Bench {
	return &Bench{
		cfg: cfg,
	}
}

func (b *Bench) LoadCsvToRedis() {
	// open file
	f, err := os.Open(b.cfg.Files.Csv)
	if err != nil {
		log.Error().Msgf("failed Open csv %#v", err)
		return
	}
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Error().Msgf("failed to read csv %#v", err)
		return
	}

	for i, line := range data {
		if i > 0 { // omit header line
			var rec ShoppingRecord
			for j, field := range line {
				if j == 0 {
					rec.Vegetable = field
				} else if j == 1 {
					rec.Fruit = field
				}
			}
			shoppingList = append(shoppingList, rec)
		}
	}
}

func (b *Bench) FetchFromRedis() {

}
