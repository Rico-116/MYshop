package util

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var MQConn *amqp.Connection
var MQChannel *amqp.Channel

const (
	OrderDelayExchange        = "order.delay.exchange"
	OrderCloseExchange        = "order.close.exchange"
	OrderCloseDelayQueue      = "q.order.close.delay"
	OrderCloseExecuteQueue    = "q.order.close.execute"
	OrderCloseDelayRoutingKey = "order.close.delay"
	OrderCloseRoutingKey      = "order.close"
)

func InitRabbitMQ() error {
	var err error
	dsn := "amqp://guest:guest@192.168.0.147:5672/"
	MQConn, err = amqp.Dial(dsn)
	if err != nil {
		return fmt.Errorf("MQ链接失败: %w", err)
	}
	MQChannel, err = MQConn.Channel()
	if err != nil {
		return fmt.Errorf("MQ通道打开失败: %w", err)
	}
	if err = MQChannel.Confirm(false); err != nil {
		return fmt.Errorf("rabbitmq confirm mode failed: %w", err)
	}
	//延迟交换机
	err = MQChannel.ExchangeDeclare(
		OrderDelayExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		//logger.Log.Error("声明延迟交换机失败", zap.Error(err))
		return fmt.Errorf("声明延迟交换机失败: %w", err)
	}
	err = MQChannel.ExchangeDeclare(
		OrderCloseExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		//logger.Log.Error("声明执行关单交换机失败", zap.Error(err))
		return fmt.Errorf("声明执行关单交换机失败: %w", err)
	}
	delayArgs := amqp.Table{
		"x-message-ttl":             int32(60 * 30 * 1000), //30min
		"x-dead-letter-exchange":    OrderCloseExchange,
		"x-dead-letter-routing-key": OrderCloseRoutingKey,
	}
	_, err = MQChannel.QueueDeclare(
		OrderCloseDelayQueue,
		true,
		false,
		false,
		false,
		delayArgs,
	)
	if err != nil {
		//logger.Log.Error("声明延迟队列失败", zap.Error(err))
		return fmt.Errorf("声明延迟队列失败: %w", err)
	}
	_, err = MQChannel.QueueDeclare(
		OrderCloseExecuteQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		//logger.Log.Error("声明执行关单队列失败", zap.Error(err))
		return fmt.Errorf("声明执行关单队列失败: %w", err)
	}
	err = MQChannel.QueueBind(
		OrderCloseDelayQueue,
		OrderCloseDelayRoutingKey,
		OrderDelayExchange,
		false,
		nil,
	)
	if err != nil {
		//logger.Log.Error("绑定延迟队列失败", zap.Error(err))
		return fmt.Errorf("绑定延迟队列失败: %w", err)
	}
	err = MQChannel.QueueBind(
		OrderCloseExecuteQueue,
		OrderCloseRoutingKey, //转到延迟队列使用
		OrderCloseExchange,
		false,
		nil,
	)
	if err != nil {
		//logger.Log.Error("绑定执行关单队列失败", zap.Error(err))
		return fmt.Errorf("绑定执行关单队列失败: %w", err)
	}
	log.Println("RabbitMQ初始化成功")
	return nil
}
func CloseRabbitMQ() {
	if MQChannel != nil {
		_ = MQChannel.Close()
	}
	if MQConn != nil {
		_ = MQConn.Close()
	}
}
func PublishWithConfirm(exchange, routingKey string, body []byte) error {
	confirms := MQChannel.NotifyPublish(make(chan amqp.Confirmation, 1))
	err := MQChannel.Publish(exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Timestamp:   time.Now(),
			Body:        body,
		},
	)
	if err != nil {
		return err
	}
	confirmed := <-confirms
	if !confirmed.Ack {
		return fmt.Errorf("message publish not confirmed")
	}
	return nil
}
