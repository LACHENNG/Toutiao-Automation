package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

// 创建测试函数
func TestCanPublish_NormalCase(t *testing.T) {
	// 创建 Gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTime := NewMockTimeProvider(ctrl)
	mockRand := NewMockRandProvider(ctrl)
	mockLogger := NewMockLogger(ctrl)

	// 创建 PublisherController 控制器，假设每小时最多发布2次，最小间隔为1秒
	pc := NewPublisherController(2, time.Second, time.Minute, mockTime, mockRand, mockLogger)

	// 模拟当前时间
	mockTime.EXPECT().Now().Return(time.Now()).Times(1)

	// 模拟随机数返回 0.5
	mockRand.EXPECT().Float64().Return(0.5).Times(1)

	// 模拟日志输出
	mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).Times(1)

	// 执行 CanPublish
	ctx := context.Background()
	canPublish, err := pc.CanPublish(ctx)

	// 验证
	assert.NoError(t, err)
	assert.True(t, canPublish)

	// 验证模拟行为

}

func TestCanPublish_MaxPerHourReached(t *testing.T) {
	// 创建 Gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTime := NewMockTimeProvider(ctrl)
	mockRand := NewMockRandProvider(ctrl)
	mockLogger := NewMockLogger(ctrl)

	// 创建 PublisherController 控制器，假设每小时最多发布2次
	pc := NewPublisherController(2, time.Second, time.Minute, mockTime, mockRand, mockLogger)

	// 模拟当前时间
	mockTime.EXPECT().Now().Return(time.Now()).Once()

	// 模拟已经有2次发布记录
	pc.publishTimes = append(pc.publishTimes, time.Now(), time.Now())

	// 验证日志输出
	mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).Once()

	// 执行 CanPublish
	ctx := context.Background()
	canPublish, err := pc.CanPublish(ctx)

	// 验证
	assert.NoError(t, err)
	assert.False(t, canPublish)

	// 验证模拟行为
	ctrl.VerifyAll()
}

func TestCanPublish_MinIntervalNotMet(t *testing.T) {
	// 创建 Gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTime := NewMockTimeProvider(ctrl)
	mockRand := NewMockRandProvider(ctrl)
	mockLogger := NewMockLogger(ctrl)

	// 创建 PublisherController 控制器，设置最小间隔为2秒
	pc := NewPublisherController(10, 2*time.Second, time.Minute, mockTime, mockRand, mockLogger)

	// 模拟当前时间
	mockTime.EXPECT().Now().Return(time.Now()).Once()

	// 模拟发布记录，最后一次发布时间为1秒前
	pc.publishTimes = append(pc.publishTimes, time.Now().Add(-1*time.Second))

	// 验证日志输出
	mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).Once()

	// 执行 CanPublish
	ctx := context.Background()
	canPublish, err := pc.CanPublish(ctx)

	// 验证
	assert.NoError(t, err)
	assert.False(t, canPublish)

	// 验证模拟行为
	ctrl.VerifyAll()
}

func TestCanPublish_TimeProviderError(t *testing.T) {
	// 创建 Gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTime := NewMockTimeProvider(ctrl)
	mockRand := NewMockRandProvider(ctrl)
	mockLogger := NewMockLogger(ctrl)

	// 创建 PublisherController
	pc := NewPublisherController(10, time.Second, time.Minute, mockTime, mockRand, mockLogger)

	// 模拟时间提供者返回错误
	mockTime.EXPECT().Now().Return(time.Time{}, fmt.Errorf("time provider error")).Once()

	// 执行 CanPublish
	ctx := context.Background()
	canPublish, err := pc.CanPublish(ctx)

	// 验证
	assert.Error(t, err)
	assert.False(t, canPublish)

	// 验证模拟行为
	ctrl.VerifyAll()
}
