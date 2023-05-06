# prometheus 和Grafana

```sh
docker run -d --name=prometheus --network prometheus-net -p 9090:9090 -v "$(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml" prom/prometheus  
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

go_gc_duration_seconds

go_gc_duration_seconds_count

go_gc_duration_seconds_sum

go_goroutines

go_info

go_memstats_alloc_bytes

go_memstats_alloc_bytes_total

go_memstats_buck_hash_sys_bytes

go_memstats_frees_total

go_memstats_gc_sys_bytes

go_memstats_heap_alloc_bytes

go_memstats_heap_idle_bytes

go_memstats_heap_inuse_bytes

go_memstats_heap_objects

go_memstats_heap_released_bytes

go_memstats_heap_sys_bytes

go_memstats_last_gc_time_seconds