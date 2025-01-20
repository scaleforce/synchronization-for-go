package hr

import (
	"encoding/json"
	"fmt"
	"log/slog"

	envelopemessage "github.com/scaleforce/synchronization-for-go/internal/message/envelope"
	hrevent "github.com/scaleforce/synchronization-for-go/pkg/message/event/hr"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
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
	return hrevent.DiscriminatorPosition
}

func (handler *PositionEventHandler) Handle(message pubsub.Message) error {
	receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

	if !ok {
		return envelopemessage.ErrInvalidReceivedEnvelope
	}

	positionEvent, ok := receivedEnvelope.Message.(*hrevent.PositionEvent)

	if !ok {
		return pubsub.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(positionEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
