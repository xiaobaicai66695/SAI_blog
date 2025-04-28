package common

import (
	"github.com/IBM/sarama"
)

func InitProducer() (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
