# prometheus 和Grafana

```sh
docker run -d --name prometheus   --user 473:473   --network prometheus-net   -p 9090:9090   --mount type=bind,source=/home/zhaobuqi/docker_files/prometheus.yml,target=/etc/prometheus/prometheus.yml   --mount type=bind,source=/home/zhaobuqi/docker_files/prometheus-data,target=/prometheus   prom/prometheus
```

```go
docker run -d --name=grafana -p 3000:3000 grafana/grafana
```

http://localhost:3000

```sh
docker network create prometheus-net
```

http://prometheus:9090

```shell
docker run -d --name grafana --network prometheus-net -p 3000:3000   -v /home/zhaobuqi/docker_files/grafana-data:/var/lib/grafana   -e "GF_AUTH_ANONYMOUS_ENABLED=true"     -e "GF_AUTH_ANONYMOUS_ORG_ROLE=Viewer"   -e "GF_SERVER_ALLOW_EMBEDDING=true"   -e "GF_SERVER_ENABLE_CORS=true"   --user 472:472  grafana/grafana 
```

1. `requestsTotal`（请求总数）：您可以使用**折线图**或**柱状图**来显示请求总数随时间的变化。这将帮助您了解请求量的趋势。
2. `responseTimeHistogram`（响应时间直方图）：**热图**或**直方图**是展示响应时间分布的理想选择。您还可以使用**折线图**显示响应时间的百分位数随时间的变化，以了解性能变化。
3. `errorRate`（错误率）：错误率可以使用**折线图**或**柱状图**来展示。这将帮助您了解错误发生的频率和趋势。
4. `memoryUsage`（内存使用情况）：**折线图**或**面积图**是展示内存使用情况随时间变化的理想选择。这将帮助您了解内存使用情况的趋势，以及是否存在潜在的内存泄漏。
5. `cpuUsage`（CPU使用率）：与内存使用情况类似，**折线图**或**面积图**是展示CPU使用率随时间变化的理想选择。
6. `limiterCount`（限制器事件计数）：您可以使用**折线图**或**柱状图**来显示限制器事件计数随时间的变化。这将帮助您了解限制器事件发生的频率和趋势。
7. `circuitBreakerCount`（熔断器事件计数）：**折线图**或**柱状图**是展示熔断器事件计数随时间变化的理想选择。您还可以为不同的熔断器状态（如 open、closed 或 half_open）使用不同的线条或柱子颜色。