package masterdata

import (
	"encoding/json"
	"fmt"

	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/masterdata"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsuberr"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsublog"
)

type CircleEventHandler struct {
	logger pubsublog.Logger
}

func NewCircleEventHandler(logger pubsublog.Logger) *CircleEventHandler {
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
