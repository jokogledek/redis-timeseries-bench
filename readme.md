### Overview
read & write performance benchmark for redis timeseries between rueidis & redistimeseries-go

### Quick start
https://redis.io/docs/stack/timeseries/quickstart/
```shell
docker run -p 0.0.0.0:16379:6379 -d --rm --name=rds_timeseries redislabs/redistimeseries
```

### Client libraries
| Project                                      | Language | License  | Author                             | Stars                                                 | Bulk Insert |
|----------------------------------------------|----------|----------|------------------------------------|-------------------------------------------------------|-------------|
| [redistimeseries-go][redistimeseries-go-url] | Go       | Apache-2 | [Redis][redistimeseries-go-author] | [![redistimeseries-go-stars]][redistimeseries-go-url] | 0 ms        |
| [rueidis][rueidis-url]                       | Go       | Apache-2 | [Rueian][rueidis-author]           | [![rueidis-stars]][rueidis-url]                       | 416.76 ms   |

[redistimeseries-go-url]: https://github.com/RedisTimeSeries/redistimeseries-go/
[redistimeseries-go-author]: https://redis.com
[redistimeseries-go-stars]: https://img.shields.io/github/stars/RedisTimeSeries/redistimeseries-go.svg?style=social&amp;label=Star&amp;maxAge=2592000

[rueidis-url]: https://github.com/rueian/rueidis
[rueidis-author]: https://github.com/rueian
[rueidis-stars]: https://img.shields.io/github/stars/rueian/rueidis.svg?style=social&amp;label=Star&amp;maxAge=2592000

### Test Data
* 681.540 rows
* 100.000 rows / chunk
* test on 8 core / 16GB machine