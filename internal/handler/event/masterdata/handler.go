package masterdata

import (
	"encoding/json"
	"errors"
	"log"

	event "github.com/scaleforce/synchronization-sdk-for-go/pkg/message/event/masterdata"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/pubsub"
)

var ErrInvalidDiscriminator = errors.New("invalid discriminator")

type CityEventHandler struct{}

func (handler *CityEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorCity
}

func (handler *CityEventHandler) Create() pubsub.Message {
	return &event.CityEvent{}
}

func (handler *CityEventHandler) Handle(message pubsub.Message) error {
	cityEvent, ok := message.(*event.CityEvent)

	if !ok {
		return ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(cityEvent, "", "  ")

	if err != nil {
		return err
	}

	log.Println(string(data))

	return nil
}

type CircleEventHandler struct{}

func (handler *CircleEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorCircle
}

func (handler *CircleEventHandler) Create() pubsub.Message {
	return &event.CircleEvent{}
}

func (handler *CircleEventHandler) Handle(message pubsub.Message) error {
	circleEvent, ok := message.(*event.CircleEvent)

	if !ok {
		return ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(circleEvent, "", "  ")

	if err != nil {
		return err
	}

	log.Println(string(data))

	return nil
}

type ZoneEventHandler struct{}

func (handler *ZoneEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorZone
}

func (handler *ZoneEventHandler) Create() pubsub.Message {
	return &event.ZoneEvent{}
}

func (handler *ZoneEventHandler) Handle(message pubsub.Message) error {
	zoneEvent, ok := message.(*event.ZoneEvent)

	if !ok {
		return ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(zoneEvent, "", "  ")

	if err != nil {
		return err
	}

	log.Println(string(data))

	return nil
}
