package masterdata

import (
	"encoding/json"
	"fmt"
	"log/slog"

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
	zoneEvent, ok := message.(*event.ZoneEvent)

	if !ok {
		return pubsub.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(zoneEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
