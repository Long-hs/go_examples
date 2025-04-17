package producer

const (
	logPrefixSync    = "[同步生产者] "
	logPrefixAsync   = "[异步生产者] "
	logPrefixService = "[生产者服务] "

	broker     = "localhost:9092"
	syncTopic  = "kafka-example-sync"
	asyncTopic = "kafka-example-async"
)

// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	return err != nil
}
