package hr

import (
	"encoding/json"
	"fmt"
	"log/slog"

	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/hr"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsuberr"
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

func (handler *EmployeeEventHandler) Create() pubsub.Message {
	return &event.EmployeeEvent{}
}

func (handler *EmployeeEventHandler) Handle(message pubsub.Message) error {
	employeeEvent, ok := message.(*event.EmployeeEvent)

	if !ok {
		return pubsuberr.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(employeeEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
