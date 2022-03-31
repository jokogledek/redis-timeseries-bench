
build_rueidis:
	@go build -o ./cmd/rueidis/bench_rueidis ./cmd/rueidis

build_goredis:
	@go build -o ./cmd/goredis/goredis ./cmd/goredis