package masterdata

import (
	"encoding/json"
	"fmt"
	"log/slog"

	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/masterdata"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsuberr"
)

type CircleEventHandler struct {
	logger *slog.Logger
}

func NewCircleEventHandler(logger *slog.Logger) *CircleEventHandler {
	return &CircleEventHandler{
		logger: logger,
	}
}

func (handler *CircleEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorCircle
}

func (handler *CircleEventHandler) Create() pubsub.Message {
	return &event.CircleEvent{}
}

func (handler *CircleEventHandler) Handle(message pubsub.Message) error {
	circleEvent, ok := message.(*event.CircleEvent)

	if !ok {
		return pubsuberr.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(circleEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
