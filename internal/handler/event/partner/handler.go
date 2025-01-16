package partner

import (
	"encoding/json"
	"errors"
	"log"

	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/partner"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

var ErrInvalidDiscriminator = errors.New("invalid discriminator")

type PartnerGroupEventHandler struct{}

func (handler *PartnerGroupEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorPartnerGroup
}

func (handler *PartnerGroupEventHandler) Create() pubsub.Message {
	return &event.PartnerGroupEvent{}
}

func (handler *PartnerGroupEventHandler) Handle(message pubsub.Message) error {
	partnerGroupEvent, ok := message.(*event.PartnerGroupEvent)

	if !ok {
		return ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(partnerGroupEvent, "", "  ")

	if err != nil {
		return err
	}

	log.Println(string(data))

	return nil
}

type PartnerEventHandler struct{}

func (handler *PartnerEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorPartner
}

func (handler *PartnerEventHandler) Create() pubsub.Message {
	return &event.PartnerEvent{}
}

func (handler *PartnerEventHandler) Handle(message pubsub.Message) error {
	partnerEvent, ok := message.(*event.PartnerEvent)

	if !ok {
		return ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(partnerEvent, "", "  ")

	if err != nil {
		return err
	}

	log.Println(string(data))

	return nil
}
