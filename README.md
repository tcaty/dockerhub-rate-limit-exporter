# Docker Hub rate limit exporter

This tool created for monitoring Dokcer Hub rate limits according to [this](https://www.docker.com/blog/checking-your-current-docker-pull-rate-limits-and-status/) tutorial. 

## Usage

There are two commands:
* `fetch` - fetch and print current rate limits to console
* `scrape` - scrape rate limits with some interval and export them in Prometheus format

### fetch

If you want to understand your rate limits immediatly, you should use command below:
```
docker run --rm -it --entrypoint /exporter tcaty/dockerhub-rate-limit-exporter fetch
```
Output:
```
Mode: anonymous
Repository: ratelimitpreview/test
Host: xxx.xxx.xxx.xxx
RateLimit [Total]: 100
RateLimit [Remaining]: 100
```

### scrape

If you want to have access to rate limits usage history, you should run scraper command:
```
docker run \
  --rm --publish "8080:8080" \
  -it --entrypoint /exporter tcaty/dockerhub-rate-limit-exporter scrape
```
Check metrics with command:
```
curl -s 0.0.0.0:8080/metrics | head -n 6
```
Output:
```
# HELP dockerhub_rate_limit_remaining RateLimit-Remaining
# TYPE dockerhub_rate_limit_remaining gauge
dockerhub_rate_limit_remaining{host="xxx.xxx.xxx.xxx",mode="anonymous",username=""} 100
# HELP dockerhub_rate_limit_total RateLimit-Limit
# TYPE dockerhub_rate_limit_total gauge
dockerhub_rate_limit_total{host="xxx.xxx.xxx.xxx",mode="anonymous",username=""} 100
```

## Configuration

You can configure exporter via command line argmunets. Use command below to display available flags:
```
docker run --rm -it --entrypoint /app dockerhub-rate-limit-exporter --help
```
Also you can set flag value by using environment variable. For example, flag `--username` corresponds to `$USERNAME` env variable.

## Examples

There are two examples in `examples/` folder:
* `simple` - simple example with one exporter instance in anonymous mode.
* `auth-n-anon` - example with several exporter instances in anonymous and authenticated modes.

Both of them provide provisioning grafana dashboard for rate limit metrics. It looks like this:

![Без имени](https://github.com/tcaty/dockerhub-rate-limit-exporter/assets/79706809/ffbb3050-33ae-45a4-bceb-c3eba45c84c4)

To run any of them move to their folder and run command:
```
docker-compose up -d
```
Then open in browser:
* http://0.0.0.0:3000 - grafana
* http://0.0.0.0:9090 - prometheus
