package hr

import (
	"encoding/json"
	"errors"
	"log"

	event "github.com/scaleforce/synchronization-sdk-for-go/pkg/message/event/hr"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/pubsub"
)

var ErrInvalidDiscriminator = errors.New("invalid discriminator")

type EmployeeEventHandler struct{}

func (handler *EmployeeEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorEmployee
}

func (handler *EmployeeEventHandler) Create() pubsub.Message {
	return &event.EmployeeEvent{}
}

func (handler *EmployeeEventHandler) Handle(message pubsub.Message) error {
	employeeEvent, ok := message.(*event.EmployeeEvent)

	if !ok {
		return ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(employeeEvent, "", "  ")

	if err != nil {
		return err
	}

	log.Println(string(data))

	return nil
}

type PositionEventHandler struct{}

func (handler *PositionEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorPosition
}

func (handler *PositionEventHandler) Create() pubsub.Message {
	return &event.PositionEvent{}
}

func (handler *PositionEventHandler) Handle(message pubsub.Message) error {
	positionEvent, ok := message.(*event.PositionEvent)

	if !ok {
		return ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(positionEvent, "", "  ")

	if err != nil {
		return err
	}

	log.Println(string(data))

	return nil
}

type RoleEventHandler struct{}

func (handler *RoleEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorRole
}

func (handler *RoleEventHandler) Create() pubsub.Message {
	return &event.RoleEvent{}
}

func (handler *RoleEventHandler) Handle(message pubsub.Message) error {
	roleEvent, ok := message.(*event.RoleEvent)

	if !ok {
		return ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(roleEvent, "", "  ")

	if err != nil {
		return err
	}

	log.Println(string(data))

	return nil
}
