package partner

import (
	"encoding/json"
	"fmt"
	"log/slog"

	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/partner"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsuberr"
)

type PartnerEventHandler struct {
	logger *slog.Logger
}

func NewPartnerEventHandler(logger *slog.Logger) *PartnerEventHandler {
	return &PartnerEventHandler{
		logger: logger,
	}
}

func (handler *PartnerEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorPartner
}

func (handler *PartnerEventHandler) Create() pubsub.Message {
	return &event.PartnerEvent{}
}

func (handler *PartnerEventHandler) Handle(message pubsub.Message) error {
	partnerEvent, ok := message.(*event.PartnerEvent)

	if !ok {
		return pubsuberr.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(partnerEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
