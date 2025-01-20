package hr

import (
	"encoding/json"
	"fmt"
	"log/slog"

	envelopemessage "github.com/scaleforce/synchronization-for-go/internal/message/envelope"
	hrevent "github.com/scaleforce/synchronization-for-go/pkg/message/event/hr"
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
	return hrevent.DiscriminatorEmployee
}

func (handler *EmployeeEventHandler) Handle(message pubsub.Message) error {
	receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

	if !ok {
		return envelopemessage.ErrInvalidReceivedEnvelope
	}

	employeeEvent, ok := receivedEnvelope.Message.(*hrevent.EmployeeEvent)

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
