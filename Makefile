
build_rueidis:
	@go build -o ./cmd/rueidis/bench_rueidis ./cmd/rueidis

build_goredis:
	@go build -o ./cmd/goredis/goredis ./cmd/goredis

bench: build_goredis build_goredis
	@echo "benchmark with redistimeseries-go"
	@echo "=================================="
	@./cmd/goredis/goredis
	@echo "=================================="
	@echo "benchmark with rueidis"
	@echo "=================================="
	@./cmd/rueidis/bench_rueidis
	@echo "=================================="