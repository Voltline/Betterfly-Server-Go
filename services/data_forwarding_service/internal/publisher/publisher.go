package publisher

import (
	"context"
	"data_forwarding_service/config"
	"data_forwarding_service/internal/logger"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
)

var log = logger.NewLogger()

var RocketMQProducer rocketmq.Producer

// InitRocketMQProducer 初始化 RocketMQ 生产者
func InitRocketMQProducer() error {
	var err error
	topic := os.Getenv("HOSTNAME")
	nsServer := os.Getenv("NAMESERVER")
	if topic == "" {
		nsServer = config.DefaultNsServer
		topic = "message-topic"
	}
	log.Info.Printf("当前nsServer: %s, topic: %s\n", nsServer, topic)
	RocketMQProducer, err = rocketmq.NewProducer(
		producer.WithGroupName("message-group"),
		producer.WithNameServer([]string{nsServer}),
	)
	if err != nil {
		return fmt.Errorf("创建RocketMQ生产者错误: %v", err)
	}

	err = RocketMQProducer.Start()
	if err != nil {
		return fmt.Errorf("启动RocketMQ生产者错误: %v", err)
	}
	return nil
}

// PublishMessage 发布消息
func PublishMessage(message string) error {
	topic := os.Getenv("HOSTNAME")
	log.Info.Printf("当前Pod Topic为: %s", topic)
	if topic == "" {
		topic = "message-topic"
	}
	msg := &primitive.Message{
		// 还没想好主题怎么设置
		Topic: topic,
		Body:  []byte(message),
	}

	sendResult, err := RocketMQProducer.SendSync(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("向消息队列发布消息错误: %v", err)
	}
	log.Info.Printf("消息发布成功: %v\n", sendResult)
	return nil
}
