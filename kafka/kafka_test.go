package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
)

type handler struct {
}

func (h handler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

func TestGroupReceiver(t *testing.T) {
	ctx := context.Background()
	GroupReceiver(ctx, []string{"172.30.225.215:9092"}, "group_ps", []string{"topic_plm_mom_to_ps", "topic_personnel_master_data"}, &handler{})
}
