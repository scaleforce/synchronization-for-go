package hr

import (
	"encoding/json"
	"fmt"
	"log/slog"

	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/hr"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsuberr"
)

type PositionEventHandler struct {
	logger *slog.Logger
}

func NewPositionEventHandler(logger *slog.Logger) *PositionEventHandler {
	return &PositionEventHandler{
		logger: logger,
	}
}

func (handler *PositionEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorPosition
}

func (handler *PositionEventHandler) Create() pubsub.Message {
	return &event.PositionEvent{}
}

func (handler *PositionEventHandler) Handle(message pubsub.Message) error {
	positionEvent, ok := message.(*event.PositionEvent)

	if !ok {
		return pubsuberr.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(positionEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
