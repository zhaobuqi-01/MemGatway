package flow_counter

import (
	"context"
	"fmt"
	"gateway/pkg/database/redis"
	"sync/atomic"
	"time"
)

type RedisFlowCount interface {
	Increase()
	GetDayData(t time.Time) (int64, error)
	GetHourData(t time.Time) (int64, error)
	GetDayKey(t time.Time) string
	GetHourKey(t time.Time) string
}

type redisFlowCounter struct {
	CounterName string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	QPD         int64
}

func NewRedisFlowCountService(CounterName string, interval time.Duration) *redisFlowCounter {

	var ctx = context.Background()

	reqCounter := &redisFlowCounter{
		CounterName: CounterName,
		Interval:    interval,
		QPS:         0,
		Unix:        0,
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
					return pipe.IncrBy(ctx, dayKey, tickerCount).Err()
				},
				func(pipe redis.Pipeliner) error {
					return pipe.Expire(ctx, dayKey, 86400*2*time.Second).Err()
				},
				func(pipe redis.Pipeliner) error {
					return pipe.IncrBy(ctx, hourKey, tickerCount).Err()
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
			tickerCount = totalCount - reqCounter.QPD
			if nowUnix > reqCounter.Unix {
				reqCounter.QPD = totalCount
				reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
				reqCounter.Unix = time.Now().Unix()
			}
		}
	}()
	return reqCounter
}

func (o *redisFlowCounter) GetDayKey(t time.Time) string {
	dayStr := t.Format("20060102")
	return fmt.Sprintf("%s_%s_%s", "flow_day_count", dayStr, o.CounterName)
}

func (o *redisFlowCounter) GetHourKey(t time.Time) string {
	hourStr := t.Format("2006010215")
	return fmt.Sprintf("%s_%s_%s", "flow_hour_count", hourStr, o.CounterName)
}

func (o *redisFlowCounter) GetHourData(t time.Time) (int64, error) {
	return redis.GetInt64(o.GetHourKey(t))
}

func (o *redisFlowCounter) GetDayData(t time.Time) (int64, error) {
	return redis.GetInt64(o.GetDayKey(t))
}

// 原子增加
func (o *redisFlowCounter) Increase() {
	atomic.AddInt64(&o.TickerCount, 1)
}
