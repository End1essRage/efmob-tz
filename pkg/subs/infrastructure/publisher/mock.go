package publisher

import (
	"context"

	"github.com/end1essrage/efmob-tz/pkg/common/logger"
)

type MockPublisher struct {
}

func NewMockPublisher() *MockPublisher {
	return &MockPublisher{}
}

func (*MockPublisher) Publish(ctx context.Context, topic string, payload []byte) error {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "MockPublisher",
		Func: "Publish",
		Ctx:  ctx,
	})

	log.Infof("topic=%s payload=%s", topic, string(payload))
	return nil
}
