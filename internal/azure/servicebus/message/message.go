package message

import (
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/scaleforce/synchronization-for-go/internal/message/envelope"
	envelopemessage "github.com/scaleforce/synchronization-for-go/internal/message/envelope"
	"github.com/scaleforce/synchronization-for-go/pkg/azure/servicebus"
	"github.com/scaleforce/synchronization-for-go/pkg/message/event/hr"
	"github.com/scaleforce/synchronization-for-go/pkg/message/event/masterdata"
	"github.com/scaleforce/synchronization-for-go/pkg/message/event/partner"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

func NewMarshalMessageFunc() servicebus.MarshalMessageFunc {
	return func(message pubsub.Message) (*azservicebus.Message, error) {
		body, err := json.Marshal(message)

		if err != nil {
			return nil, err
		}

		serviceBusMessage := &azservicebus.Message{
			Body: body,
		}

		return serviceBusMessage, nil
	}
}

func CreateMessage(discriminator pubsub.Discriminator) pubsub.Message {
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

type CreateMessageFunc func(discriminator pubsub.Discriminator) pubsub.Message

func NewUnmarshalMessageFunc(createMessageFunc CreateMessageFunc) servicebus.UnmarshalMessageFunc {
	return func(serviceBusReceivedMessage *azservicebus.ReceivedMessage) (pubsub.Message, error) {
		partialMessage := &struct {
			Type string `json:"Type"`
		}{}

		if err := json.Unmarshal(serviceBusReceivedMessage.Body, &partialMessage); err != nil {
			return nil, err
		}

		discriminator := pubsub.Discriminator(partialMessage.Type)

		message := createMessageFunc(discriminator)

		if err := json.Unmarshal(serviceBusReceivedMessage.Body, &message); err != nil {
			return nil, err
		}

		return message, nil
	}
}

func NewMarshalEnvelopeFunc(marshalMessageFunc servicebus.MarshalMessageFunc) servicebus.MarshalMessageFunc {
	return func(message pubsub.Message) (*azservicebus.Message, error) {
		envelope, ok := message.(*envelopemessage.Envelope)

		if !ok {
			return nil, envelopemessage.ErrInvalidEnvelope
		}

		serviceBusMessage, err := marshalMessageFunc(envelope.Message)

		if err != nil {
			return nil, err
		}

		serviceBusMessage.ApplicationProperties = envelope.ApplicationProperties
		serviceBusMessage.SessionID = envelope.SessionID
		serviceBusMessage.MessageID = envelope.MessageID

		return serviceBusMessage, nil
	}
}

func NewUnmarshalReceivedEnvelopeFunc(unmarshalMessageFunc servicebus.UnmarshalMessageFunc) servicebus.UnmarshalMessageFunc {
	return func(serviceBusReceivedMessage *azservicebus.ReceivedMessage) (pubsub.Message, error) {
		message, err := unmarshalMessageFunc(serviceBusReceivedMessage)

		if err != nil {
			return nil, err
		}

		receivedEnvelope := envelope.NewReceivedEnvelope(message)

		receivedEnvelope.ApplicationProperties = serviceBusReceivedMessage.ApplicationProperties
		receivedEnvelope.EnqueuedSequenceNumber = serviceBusReceivedMessage.EnqueuedSequenceNumber
		receivedEnvelope.EnqueuedTime = serviceBusReceivedMessage.EnqueuedTime

		return receivedEnvelope, nil
	}
}
