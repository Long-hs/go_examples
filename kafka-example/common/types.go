package common

const (
	LogPrefixSync    = "[同步生产者] "
	LogPrefixAsync   = "[异步生产者] "
	LogPrefixService = "[生产者服务] "

	Broker     = "localhost:9092"
	SyncTopic  = "kafka-example-sync"
	AsyncTopic = "kafka-example-async"
)

// IsRetryableError 判断错误是否可重试
func IsRetryableError(err error) bool {
	return err != nil
}
