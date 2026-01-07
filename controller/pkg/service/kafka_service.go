package service

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type ReaderWriterService struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewReaderWriterService(address, port, topic, group string) *ReaderWriterService {
	brokerAddress := fmt.Sprintf("%v:%v", address, port)
	err := CreateTopic(brokerAddress, topic, 0)
	if err != nil {
		return nil
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: group, // Consumer group для распределения
		//MinBytes: 10e3,  // Минимальный размер батча
		//MaxBytes: 10e6,  // Максимальный размер батча
	})

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &ReaderWriterService{reader: reader, writer: writer}
}

// CreateTopic создаёт топик, если он не существует
func CreateTopic(brokerAddress string, topicName string, partitions int) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokerAddress, topicName, partitions)
	if err != nil {
		log.Printf("Failed to dial leader: %v", err)
		return err
	}
	defer conn.Close()
	return nil
}

func (rw *ReaderWriterService) ReadMessage() (*kafka.Message, error) {
	msg, err := rw.reader.ReadMessage(context.Background())
	if err != nil {
		log.Printf("Failed to read message: %v", err)
		return nil, err
	}
	return &msg, nil
}

func (rw *ReaderWriterService) WriteMessage(message kafka.Message) error {
	err := rw.writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("Failed to write message: %v", err)
		return err
	}
	return err
}
