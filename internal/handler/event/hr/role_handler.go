package hr

import (
	"encoding/json"
	"fmt"
	"log/slog"

	envelopemessage "github.com/scaleforce/synchronization-for-go/internal/message/envelope"
	hrevent "github.com/scaleforce/synchronization-for-go/pkg/message/event/hr"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

type RoleEventHandler struct {
	logger *slog.Logger
}

func NewRoleEventHandler(logger *slog.Logger) *RoleEventHandler {
	return &RoleEventHandler{
		logger: logger,
	}
}

func (handler *RoleEventHandler) Discriminator() pubsub.Discriminator {
	return hrevent.DiscriminatorRole
}

func (handler *RoleEventHandler) Handle(message pubsub.Message) error {
	receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

	if !ok {
		return envelopemessage.ErrInvalidReceivedEnvelope
	}

	roleEvent, ok := receivedEnvelope.Message.(*hrevent.RoleEvent)

	if !ok {
		return pubsub.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(roleEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
