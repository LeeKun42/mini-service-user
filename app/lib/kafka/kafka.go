package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"net"
	"strconv"
)

func getBrokers() []string {
	brokers := viper.GetStringSlice("kafka.brokers")
	if len(brokers) <= 0 {
		panic("请配置kafka服务器节点")
	}
	return brokers
}

func CreateTopic(topic string, numPartitions int, replicationFactor int) error {

	conn, _ := kafka.Dial("tcp", getBrokers()[0])
	controller, _ := conn.Controller()
	var controllerConn *kafka.Conn
	controllerConn, _ = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	defer conn.Close()
	defer controllerConn.Close()
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
		},
	}
	err := controllerConn.CreateTopics(topicConfigs...)
	return err
}

func getWriter(topic string) *kafka.Writer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(getBrokers()...),
		Topic:    topic,
		Balancer: &kafka.Hash{}, //分区均衡算法，根据key使用hash
	}
	return writer
}

func PushMessage(topic string, key string, val string) {
	getWriter(topic).WriteMessages(context.Background(), kafka.Message{Key: []byte(key), Value: []byte(val)})
}

type Message struct {
	Key   string
	Value string
}
type readCallBack func(message Message)

func getReader(topic string, groupId string) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: getBrokers(),
		GroupID: groupId,
		Topic:   topic,
	})
	return reader
}

func ReadMessage(topic string, groupId string, callback readCallBack) {
	reader := getReader(topic, groupId)
	defer reader.Close()
	fmt.Println("start ReadMessage ... ")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("kafka read msg err：", err.Error())
		}
		callback(Message{Key: string(m.Key), Value: string(m.Value)})
	}
}
