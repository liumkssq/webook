package ioc

import (
	"github.com/IBM/sarama"
	events2 "github.com/liumkssq/webook/interactive/events"
	"github.com/liumkssq/webook/interactive/repository/dao"
	"github.com/liumkssq/webook/internal/events"
	"github.com/liumkssq/webook/pkg/migrator/events/fixer"
	"github.com/spf13/viper"
)

func InitSaramaClient() sarama.Client {
	type Config struct {
		Addr []string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, scfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSaramaSyncProducer(client sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return p
}

func InitConsumers(c1 *events2.InteractiveReadEventConsumer, fixConsumer *fixer.Consumer[dao.Interactive]) []events.Consumer {
	return []events.Consumer{c1, fixConsumer}
}
