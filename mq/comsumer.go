package mq

import "fmt"

var done chan bool

// 开始监听队列,获取消息
func StartConsumer(qName, cName string, callback func(msg []byte) bool) {
	// 通过channel.Consume获取消息信道
	msgs, err := channel.Consume(qName, cName, true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	done = make(chan bool)
	// 循环的从信道里面获取消息,没有的话就阻塞
	go func() {
		for msg := range msgs {
			processSuc := callback(msg.Body)
			if !processSuc {
				// TODO: 写到另一个队列,用于异常重试
			}
		}
	}()

	// 调用callback方法来处理新的消息
	<-done
}
