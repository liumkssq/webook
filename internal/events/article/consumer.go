package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/liumkssq/webook/internal/repository"
	"github.com/liumkssq/webook/pkg/logger"
	"github.com/liumkssq/webook/pkg/saramax"
	"time"
)

type KafkaConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewKafkaConsumer(l logger.LoggerV1, repo repository.InteractiveRepository, client sarama.Client) *KafkaConsumer {
	return &KafkaConsumer{l: l, repo: repo, client: client}
}

func (r *KafkaConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", r.client)
	if err != nil {
		return err
	}
	go func() {
		cg.Consume(context.Background(),
			[]string{"article_read"},
			saramax.NewHandler[ReadEvent](r.l, r.Consume))
	}()

	return err
}

func (r *KafkaConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.repo.IncrReadCnt(ctx, "article", t.Aid)
}
