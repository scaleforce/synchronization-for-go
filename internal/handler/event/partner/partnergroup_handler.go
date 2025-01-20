package partner

import (
	"encoding/json"
	"fmt"
	"log/slog"

	envelopemessage "github.com/scaleforce/synchronization-for-go/internal/message/envelope"
	partnerevent "github.com/scaleforce/synchronization-for-go/pkg/message/event/partner"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
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
	return partnerevent.DiscriminatorPartnerGroup
}

func (handler *PartnerGroupEventHandler) Handle(message pubsub.Message) error {
	receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

	if !ok {
		return envelopemessage.ErrInvalidReceivedEnvelope
	}

	partnerGroupEvent, ok := receivedEnvelope.Message.(*partnerevent.PartnerGroupEvent)

	if !ok {
		return pubsub.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(partnerGroupEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
