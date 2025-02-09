package main

import (
	"context"
	"log"
	"strings"
	"time"

	cdp "github.com/chromedp/cdproto/cdp"
	target "github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

type DraftPublisher struct {
	publishController PublisherController
}

func createHeadlessDriver() (context.Context, context.CancelFunc) {
	// 设置 Chrome 的选项（无头模式）
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // 可选：是否无头模式
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("window-size", "1920,1080"),
		chromedp.Flag("disable-geolocation", true), // 禁用地理位置
	)

	// 创建一个新的 Chrome 上下文和取消函数
	allocator, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ := chromedp.NewContext(allocator)

	return ctx, cancel
}

func waitLogin(ctx context.Context) error {
	var username = "17355605329"
	var password = "Aa123456"

	// 打开登录页面
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://mp.toutiao.com/"),
		chromedp.WaitVisible(`li[aria-label="账密登录"]`),
		chromedp.Click(`li[aria-label="账密登录"]`),
		chromedp.WaitVisible(`input[name="normal-input"]`),
		chromedp.SendKeys(`input[name="normal-input"]`, username),
		chromedp.SendKeys(`input[name="button-input"]`, password),
		chromedp.Click(`span[class="web-login-confirm-info__checkbox"]`),
		chromedp.Click(`button[class="web-login-button"]`),
		// chromedp.WaitVisible(`div[class="mp-menu-wrapper f-min-scroll f-hover-scroll"]`),
	)

	if err != nil {
		return err
	}

	// 等待登陆成功后跳转
	err = WaitForURLContains(ctx, "profile_v4/index", 1500*time.Second)
	return err
}

func publishDrafts(ctx context.Context, controller PublisherController) error {
	var draftsXPath = "//div[contains(@class, 'draft-item')]"
	var editButtonXPath = ".//a[contains(@class, 'op-button') and text()='编辑']"
	var publishButtonXPath = "//button[contains(@class, 'publish-content')]"

	for {
		// 导航到草稿页面
		err := chromedp.Run(ctx,
			chromedp.Navigate("https://mp.toutiao.com/profile_v4/manage/draft"),
			chromedp.WaitVisible("//div[contains(@class, 'list')]"),
		)
		if err != nil {
			return err
		}

		var draftNodes []*cdp.Node
		err = chromedp.Run(ctx,
			chromedp.Nodes(draftsXPath, &draftNodes, chromedp.AtLeast(1)),
		)
		if err != nil || len(draftNodes) == 0 {
			log.Println("没有更多草稿了，等待 1 小时后再重试")
			time.Sleep(1 * time.Hour)
			continue
		}

		for _, draft := range draftNodes {
			_ = draft
			// 等待发布条件
			for {
				canPub, wait := controller.CanPublish(ctx)
				if !canPub {
					time.Sleep(wait)
					continue
				}
				break
			}

			log.Println("处理草稿...")

			// 监听新窗口
			newTargetChan := chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
				return true
			})

			err = chromedp.Run(ctx,
				chromedp.Click(editButtonXPath),
			)
			if err != nil {
				log.Println("点击编辑按钮失败:", err)
				continue
			}

			// 获取新窗口 ID
			newTarget := <-newTargetChan
			if newTarget == "" {
				log.Println("未能检测到新窗口")
				continue
			}

			// 切换到新标签页
			newTargetCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(newTarget))
			defer cancel()

			// 执行发布
			err = chromedp.Run(newTargetCtx,
				chromedp.WaitVisible(publishButtonXPath),
				chromedp.MouseClickXY(200, 200), // 收起发文助手
				chromedp.ScrollIntoView(publishButtonXPath),
				chromedp.WaitEnabled(publishButtonXPath),
				chromedp.Click(publishButtonXPath),
				chromedp.Sleep(3*time.Second),
			)
			if err != nil {
				log.Println("发布草稿失败:", err)
				continue
			}

			log.Println("发布成功")

			// 关闭新标签页
			chromedp.Cancel(newTargetCtx)

			// 刷新当前标签页
			err = chromedp.Run(ctx,
				chromedp.Reload(),
			)
			if err != nil {
				log.Println("刷新页面失败:", err)
			}

		}
	}
}

func main() {
	// 创建无头浏览器上下文
	ctx, cancel := createHeadlessDriver()
	defer cancel()

	// 创建发布频率控制器
	// 规则：每小时最多发布8条，间隔3～8分钟
	controller := NewPublisherController(8, 3*time.Minute, 15*time.Minute, &defaultTimeProvider{}, &defaultRandProvider{}, log.Default())

	// 执行登录操作
	if err := waitLogin(ctx); err != nil {
		log.Fatal("登录失败:", err)
	}

	// 发布草稿
	for {
		if err := publishDrafts(ctx, controller); err != nil {
			log.Default().Println("发布草稿失败:", err)
		}

		time.Sleep(1 * time.Hour)
	}

	// log.Println("脚本执行完成")
}

// WaitForURLContains 等待当前页面 URL 包含指定的字符串
func WaitForURLContains(ctx context.Context, substr string, timeout time.Duration) error {
	deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		var currentURL string
		err := chromedp.Run(deadlineCtx,
			chromedp.Location(&currentURL),
		)
		if err != nil {
			return err
		}

		if strings.Contains(currentURL, substr) {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}
}
