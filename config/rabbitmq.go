package config

const (
	// 是否开启文件异步转移
	AsyncTransferEnable = true
	// rabbit mq服务的入口url
	RabbitURL = "amqp://guest:guest@192.168.123.91:5672"
	// 用与文件transfer的交换机
	TransExchangeName = "uploadserver.trans"
	// oss转移队列名称
	TransOSSQueueName = "uploadserver.trans.oss"
	// 转移失败后,写入另一个队列名
	TransOSSErrQueueName = "uploadserver.trans.oss.err"
	// routing key
	TransOSSRoutingKey = "oss"
)
