package partner

import (
	"encoding/json"
	"fmt"
	"log/slog"

	envelopemessage "github.com/scaleforce/synchronization-for-go/internal/message/envelope"
	partnerevent "github.com/scaleforce/synchronization-for-go/pkg/message/event/partner"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
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
	return partnerevent.DiscriminatorPartner
}

func (handler *PartnerEventHandler) Handle(message pubsub.Message) error {
	receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

	if !ok {
		return envelopemessage.ErrInvalidReceivedEnvelope
	}

	partnerEvent, ok := receivedEnvelope.Message.(*partnerevent.PartnerEvent)

	if !ok {
		return pubsub.ErrInvalidDiscriminator
	}

	// Replace with your own synchronization logic
	data, err := json.MarshalIndent(partnerEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))
	// Replace with your own synchronization logic

	return nil
}
