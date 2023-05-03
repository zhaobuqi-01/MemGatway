package flow_control

import (
	"context"
	"fmt"
	"gateway/internal/pkg"
	"gateway/pkg/database/redis"
	"sync/atomic"
	"time"
)

type RedisFlowCountService struct {
	AppID       string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}

func NewRedisFlowCountService(appID string, interval time.Duration) *RedisFlowCountService {

	var ctx = context.Background()

	reqCounter := &RedisFlowCountService{
		AppID:    appID,
		Interval: interval,
		QPS:      0,
		Unix:     0,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount) //获取数据
			atomic.StoreInt64(&reqCounter.TickerCount, 0)            //重置数据

			currentTime := time.Now()
			dayKey := reqCounter.GetDayKey(currentTime)
			hourKey := reqCounter.GetHourKey(currentTime)
			if err := redis.Pipeline(
				func(pipe redis.Pipeliner) error {
					return pipe.Incr(ctx, dayKey).Err()
				},
				func(pipe redis.Pipeliner) error {
					return pipe.Expire(ctx, dayKey, 86400*2*time.Second).Err()
				},
				func(pipe redis.Pipeliner) error {
					return pipe.Incr(ctx, hourKey).Err()
				},
				func(pipe redis.Pipeliner) error {
					return pipe.Expire(ctx, hourKey, 86400*2*time.Second).Err()
				},
			); err != nil {
				fmt.Println("Pipeline error:", err)
				continue
			}

			totalCount, err := reqCounter.GetDayData(currentTime)
			if err != nil {
				fmt.Println("reqCounter.GetDayData err", err)
				continue
			}
			nowUnix := time.Now().Unix()
			if reqCounter.Unix == 0 {
				reqCounter.Unix = time.Now().Unix()
				continue
			}
			tickerCount = totalCount - reqCounter.TotalCount
			if nowUnix > reqCounter.Unix {
				reqCounter.TotalCount = totalCount
				reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
				reqCounter.Unix = time.Now().Unix()
			}
		}
	}()
	return reqCounter
}

func (o *RedisFlowCountService) GetDayKey(t time.Time) string {
	dayStr := t.Format("20060102")
	return fmt.Sprintf("%s_%s_%s", pkg.RedisFlowDayKey, dayStr, o.AppID)
}

func (o *RedisFlowCountService) GetHourKey(t time.Time) string {
	hourStr := t.Format("2006010215")
	return fmt.Sprintf("%s_%s_%s", pkg.RedisFlowHourKey, hourStr, o.AppID)
}

func (o *RedisFlowCountService) GetHourData(t time.Time) (int64, error) {
	return redis.GetInt64(o.GetHourKey(t))
}

func (o *RedisFlowCountService) GetDayData(t time.Time) (int64, error) {
	return redis.GetInt64(o.GetDayKey(t))
}

// 原子增加
func (o *RedisFlowCountService) Increase() {
	atomic.AddInt64(&o.TickerCount, 1)
}
