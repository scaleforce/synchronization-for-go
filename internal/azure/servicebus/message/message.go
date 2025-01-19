package message

import (
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/scaleforce/synchronization-for-go/pkg/azure/servicebus"
	"github.com/scaleforce/synchronization-for-go/pkg/message/envelope"
	"github.com/scaleforce/synchronization-for-go/pkg/message/event/hr"
	"github.com/scaleforce/synchronization-for-go/pkg/message/event/masterdata"
	"github.com/scaleforce/synchronization-for-go/pkg/message/event/partner"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

func MarshalJSONMessage(message pubsub.Message) (*azservicebus.Message, error) {
	body, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	serviceBusMessage := &azservicebus.Message{
		Body: body,
	}

	return serviceBusMessage, nil
}

func createMessage(discriminator pubsub.Discriminator) pubsub.Message {
	var message pubsub.Message

	switch discriminator {
	case masterdata.DiscriminatorCity:
		message = &masterdata.CityEvent{}
	case masterdata.DiscriminatorCircle:
		message = &masterdata.CircleEvent{}
	case masterdata.DiscriminatorZone:
		message = &masterdata.ZoneEvent{}
	case partner.DiscriminatorPartnerGroup:
		message = &partner.PartnerGroupEvent{}
	case partner.DiscriminatorPartner:
		message = &partner.PartnerEvent{}
	case hr.DiscriminatorEmployee:
		message = &hr.EmployeeEvent{}
	case hr.DiscriminatorPosition:
		message = &hr.PositionEvent{}
	case hr.DiscriminatorRole:
		message = &hr.RoleEvent{}
	}

	return message
}

func UnmarshalJSONMessage(serviceBusReceivedMessage *azservicebus.ReceivedMessage) (pubsub.Message, error) {
	partialMessage := &struct {
		Type string `json:"Type"`
	}{}

	if err := json.Unmarshal(serviceBusReceivedMessage.Body, &partialMessage); err != nil {
		return nil, err
	}

	discriminator := pubsub.Discriminator(partialMessage.Type)

	message := createMessage(discriminator)

	if err := json.Unmarshal(serviceBusReceivedMessage.Body, &message); err != nil {
		return nil, err
	}

	return message, nil
}

func NewMarshalEnvelopeMessageFunc(marshalMessageFunc servicebus.MarshalMessageFunc) servicebus.MarshalMessageFunc {
	return func(message pubsub.Message) (*azservicebus.Message, error) {
		envelopeMessage, ok := message.(*envelope.Envelope)

		if !ok {
			return nil, envelope.ErrInvalidEnvelope
		}

		return marshalMessageFunc(envelopeMessage.Message)
	}
}

func NewUnmarshalEnvelopeMessageFunc(unmarshalMessageFunc servicebus.UnmarshalMessageFunc) servicebus.UnmarshalMessageFunc {
	return func(serviceBusReceivedMessage *azservicebus.ReceivedMessage) (pubsub.Message, error) {
		message, err := unmarshalMessageFunc(serviceBusReceivedMessage)

		if err != nil {
			return nil, err
		}

		envelopeMessage := envelope.NewEnvelope(message)

		envelopeMessage.SequenceNumber = serviceBusReceivedMessage.SequenceNumber

		return envelopeMessage, nil
	}
}
