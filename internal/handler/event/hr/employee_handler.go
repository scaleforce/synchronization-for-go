package hr

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/scaleforce/synchronization-for-go/pkg/message/envelope"
	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/hr"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

type EmployeeEventHandler struct {
	logger *slog.Logger
}

func NewEmployeeEventHandler(logger *slog.Logger) *EmployeeEventHandler {
	return &EmployeeEventHandler{
		logger: logger,
	}
}

func (handler *EmployeeEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorEmployee
}

func (handler *EmployeeEventHandler) Handle(message pubsub.Message) error {
	envelopeMessage, ok := message.(*envelope.Envelope)

	if !ok {
		return envelope.ErrInvalidEnvelope
	}

	employeeEvent, ok := envelopeMessage.Message.(*event.EmployeeEvent)

	if !ok {
		return pubsub.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(employeeEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
