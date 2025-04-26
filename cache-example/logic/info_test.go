package logic

import (
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

// 公共并发测试函数
func runConcurrentTest(t *testing.T, testName string, url string, totalRequests int, timeout time.Duration) {
	var (
		wg           sync.WaitGroup
		successCount int
		failCount    int
		mutex        sync.Mutex
	)

	client := &http.Client{Timeout: timeout}
	startTime := time.Now()

	wg.Add(totalRequests)
	for i := 0; i < totalRequests; i++ {
		go func(reqID int) { // 传递reqID避免闭包变量捕获问题
			defer wg.Done()
			resp, err := client.Get(url)
			if err != nil {
				mutex.Lock()
				failCount++
				mutex.Unlock()
				t.Logf("[%s][%d] 请求失败: %v", testName, reqID, err)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Logf("[%s][%d] 关闭响应体失败: %v", testName, reqID, err)
				}
			}(resp.Body)

			statusCode := resp.StatusCode
			if statusCode == http.StatusOK {
				mutex.Lock()
				successCount++
				mutex.Unlock()
			} else {
				mutex.Lock()
				failCount++
				mutex.Unlock()
				t.Logf("[%s][%d] 请求返回非200状态码: %d", testName, reqID, statusCode)
			}
		}(i) // 传入当前循环变量i作为reqID
	}

	wg.Wait()
	duration := time.Since(startTime)

	t.Logf("[%s] 测试完成:", testName)
	t.Logf("[%s] 并发数: %d", testName, totalRequests)
	t.Logf("[%s] 成功请求数: %d", testName, successCount)
	t.Logf("[%s] 失败请求数: %d", testName, failCount)
	t.Logf("[%s] 测试耗时: %v", testName, duration)
	t.Logf("[%s] QPS: %.2f", testName, float64(totalRequests)/duration.Seconds())

	if successCount == 0 {
		t.Fatalf("[%s] 所有请求都失败了", testName)
	}
	successRate := float64(successCount) / float64(totalRequests)
	if successRate < 0.9 {
		t.Fatalf("[%s] 成功率低于90%%, 实际成功率: %.2f%%", testName, successRate*100)
	}
}

func TestHandlerCache1(t *testing.T) {
	url := "http://localhost:8080/cache1?id=1"
	runConcurrentTest(t, "Cache1", url, 200, 5*time.Second)
}

func TestHandlerMysql1(t *testing.T) {
	url := "http://localhost:8080/mysql1?id=1"
	runConcurrentTest(t, "Mysql1", url, 200, 5*time.Second)
}

func TestHandlerDoubleWrite(t *testing.T) {
	//写入数据库
	url := "http://localhost:8080/doubleWrite?id=1&name=test2"
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Logf("关闭响应体失败: %v", err)
		}
	}(resp.Body)

	statusCode := resp.StatusCode
	if statusCode == http.StatusOK {
		t.Logf("请求成功")
	} else {
		t.Fatalf("请求返回非200状态码: %d", statusCode)
	}
	//读取缓存
	url = "http://localhost:8080/cache1?id=1"
	resp, err = client.Get(url)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	statusCode = resp.StatusCode
	if statusCode == http.StatusOK {
		t.Logf("请求成功")
	} else {
		t.Fatalf("请求返回非200状态码: %d", statusCode)
	}
}
