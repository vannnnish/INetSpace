package mq

import (
	"github.com/streadway/amqp"
	"netspace/config"
	"time"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
)

func initChannel() bool {
	// 判断channel是否创建
	if channel != nil {
		return true
	}
	// 获取rabbitmq 连接
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		panic(err)
	}
	// 打开channel 用于消息发布接受
	channel, err = conn.Channel()
	if err != nil {
		panic(err)
	}
	return true
}

// 发布消息
func Publish(exchange, routingKey string, msg []byte) bool {
	// 检查channel 是否正常
	if !initChannel() {
		return false
	}
	// 调用channel 发送消息方法,
	err := channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		Headers:         nil,
		ContentType:     "text/plain",
		ContentEncoding: "",
		DeliveryMode:    0,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Time{},
		Type:            "",
		UserId:          "",
		AppId:           "",
		Body:            msg,
	})
	if err != nil {
		panic(err)
	}
	return true
}
