package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// PublisherController 用于定义发布控制器接口
type PublisherController interface {
	CanPublish(ctx context.Context) (bool, time.Duration)
}

// TimeProvider 用于获取当前时间的接口
type TimeProvider interface {
	Now() time.Time
}

// RandProvider 用于生成随机数的接口
type RandProvider interface {
	Float64() float64
}

// Logger 用于记录日志的接口
type Logger interface {
	Printf(format string, v ...interface{})
}

// publisherController 发布控制器
type publisherController struct {
	mu           sync.Mutex
	publishTimes []time.Time   // 发布记录队列
	maxPerHour   int           // 每小时最多发布次数
	minInterval  time.Duration // 最小发布间隔
	maxInterval  time.Duration // 最大发布间隔
	timeProvider TimeProvider  // 时间提供者
	randProvider RandProvider  // 随机数提供者
	logger       Logger        // 日志提供者
}

// NewPublisherController 创建发布控制器
func NewPublisherController(maxPerHour int, minInterval, maxInterval time.Duration, timeProvider TimeProvider, randProvider RandProvider, logger Logger) PublisherController {
	return &publisherController{
		maxPerHour:   maxPerHour,
		minInterval:  minInterval,
		maxInterval:  maxInterval,
		publishTimes: make([]time.Time, 0),
		timeProvider: timeProvider,
		randProvider: randProvider,
		logger:       logger,
	}
}

// CanPublish 检查是否可以发布，不阻塞，返回状态, 和等待时间
func (pc *publisherController) CanPublish(ctx context.Context) (bool, time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// 清理超过1小时的发布记录
	now := pc.timeProvider.Now()
	for len(pc.publishTimes) > 0 && now.Sub(pc.publishTimes[0]) > time.Hour {
		pc.publishTimes = pc.publishTimes[1:]
	}

	// 检查是否超过每小时最大发布次数
	if len(pc.publishTimes) >= pc.maxPerHour {
		// 计算需要等待的时间
		waitTime := pc.publishTimes[0].Add(time.Hour).Sub(now)
		pc.logger.Printf("已达到每小时发布上限，等待 %v 后重试\n", waitTime)
		return false, waitTime
	}

	// 检查是否满足最小发布间隔
	if len(pc.publishTimes) > 0 {
		lastPublish := pc.publishTimes[len(pc.publishTimes)-1]
		elapsed := now.Sub(lastPublish)
		if elapsed < pc.minInterval {
			waitTime := pc.minInterval - elapsed
			pc.logger.Printf("未达到最小发布间隔，等待 %v 后重试\n", waitTime)
			return false, waitTime
		}
	}

	// 随机化发布间隔
	if pc.maxInterval > pc.minInterval {
		randomInterval := pc.minInterval + time.Duration(float64(pc.maxInterval-pc.minInterval)*pc.randProvider.Float64())
		pc.logger.Printf("随机化发布间隔，等待 %v\n", randomInterval)
		// 不阻塞，只是记录间隔，实际测试时可以验证日志输出
	}

	// 记录本次发布时间
	pc.publishTimes = append(pc.publishTimes, pc.timeProvider.Now())
	return true, 0
}

// 默认的时间提供者实现
type defaultTimeProvider struct{}

func (d *defaultTimeProvider) Now() time.Time {
	return time.Now()
}

// 默认的随机数提供者实现
type defaultRandProvider struct{}

func (d *defaultRandProvider) Float64() float64 {
	return rand.Float64()
}

// 默认的日志提供者实现
type defaultLogger struct{}

func (d *defaultLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
