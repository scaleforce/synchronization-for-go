package partner

import (
	"encoding/json"
	"fmt"
	"log/slog"

	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/partner"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsuberr"
)

type PartnerGroupEventHandler struct {
	logger *slog.Logger
}

func NewPartnerGroupEventHandler(logger *slog.Logger) *PartnerGroupEventHandler {
	return &PartnerGroupEventHandler{
		logger: logger,
	}
}

func (handler *PartnerGroupEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorPartnerGroup
}

func (handler *PartnerGroupEventHandler) Create() pubsub.Message {
	return &event.PartnerGroupEvent{}
}

func (handler *PartnerGroupEventHandler) Handle(message pubsub.Message) error {
	partnerGroupEvent, ok := message.(*event.PartnerGroupEvent)

	if !ok {
		return pubsuberr.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(partnerGroupEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
