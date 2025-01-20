package masterdata

import (
	"encoding/json"
	"fmt"
	"log/slog"

	envelopemessage "github.com/scaleforce/synchronization-for-go/internal/message/envelope"
	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/masterdata"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

type ZoneEventHandler struct {
	logger *slog.Logger
}

func NewZoneEventHandler(logger *slog.Logger) *ZoneEventHandler {
	return &ZoneEventHandler{
		logger: logger,
	}
}

func (handler *ZoneEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorZone
}

func (handler *ZoneEventHandler) Handle(message pubsub.Message) error {
	receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

	if !ok {
		return envelopemessage.ErrInvalidReceivedEnvelope
	}

	zoneEvent, ok := receivedEnvelope.Message.(*event.ZoneEvent)

	if !ok {
		return pubsub.ErrInvalidDiscriminator
	}

	// Replace with your own synchronization logic
	data, err := json.MarshalIndent(zoneEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))
	// Replace with your own synchronization logic

	return nil
}
