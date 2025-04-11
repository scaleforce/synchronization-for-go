package util

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	emptymessage "github.com/scaleforce/synchronization-for-go/internal/message/empty"
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
	default:
		message = &emptymessage.Empty{}
	}

	return message
}

type CreateMessageFunc func(discriminator pubsub.Discriminator) pubsub.Message

func NewUnmarshalMessageFunc(createMessageFunc CreateMessageFunc) servicebus.UnmarshalMessageFunc {
	return func(serviceBusReceivedMessage *azservicebus.ReceivedMessage) (pubsub.Message, error) {
		// discriminator := pubsub.Discriminator(serviceBusReceivedMessage.ApplicationProperties["Type"].(string))

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

		receivedEnvelope := envelopemessage.NewReceivedEnvelope(message)

		receivedEnvelope.ApplicationProperties = serviceBusReceivedMessage.ApplicationProperties
		receivedEnvelope.EnqueuedSequenceNumber = serviceBusReceivedMessage.EnqueuedSequenceNumber
		receivedEnvelope.EnqueuedTime = serviceBusReceivedMessage.EnqueuedTime

		return receivedEnvelope, nil
	}
}

func GetPartitionName(message pubsub.Message) (string, error) {
	var partitionName string

	discriminator := message.Discriminator()

	switch discriminator {
	case masterdata.DiscriminatorCity:
		receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

		if !ok {
			return "", envelopemessage.ErrInvalidReceivedEnvelope
		}

		cityEvent, ok := receivedEnvelope.Message.(*masterdata.CityEvent)

		if !ok {
			return "", pubsub.ErrInvalidDiscriminator
		}

		partitionName = fmt.Sprintf("%s~%s~%s", cityEvent.TenantGroupName, cityEvent.Type, cityEvent.Data.Code)
	case masterdata.DiscriminatorCircle:
		receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

		if !ok {
			return "", envelopemessage.ErrInvalidReceivedEnvelope
		}

		circleEvent, ok := receivedEnvelope.Message.(*masterdata.CircleEvent)

		if !ok {
			return "", pubsub.ErrInvalidDiscriminator
		}

		partitionName = fmt.Sprintf("%s~%s~%s", circleEvent.TenantGroupName, circleEvent.Type, circleEvent.Data.Code)
	case masterdata.DiscriminatorZone:
		receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

		if !ok {
			return "", envelopemessage.ErrInvalidReceivedEnvelope
		}

		zoneEvent, ok := receivedEnvelope.Message.(*masterdata.ZoneEvent)

		if !ok {
			return "", pubsub.ErrInvalidDiscriminator
		}

		partitionName = fmt.Sprintf("%s~%s~%s", zoneEvent.TenantGroupName, zoneEvent.Type, zoneEvent.Data.Code)
	case partner.DiscriminatorPartnerGroup:
		receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

		if !ok {
			return "", envelopemessage.ErrInvalidReceivedEnvelope
		}

		partnerGroupEvent, ok := receivedEnvelope.Message.(*partner.PartnerGroupEvent)

		if !ok {
			return "", pubsub.ErrInvalidDiscriminator
		}

		partitionName = fmt.Sprintf("%s~%s~%s", partnerGroupEvent.TenantGroupName, partnerGroupEvent.Type, partnerGroupEvent.Data.Code)
	case partner.DiscriminatorPartner:
		receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

		if !ok {
			return "", envelopemessage.ErrInvalidReceivedEnvelope
		}

		partnerEvent, ok := receivedEnvelope.Message.(*partner.PartnerEvent)

		if !ok {
			return "", pubsub.ErrInvalidDiscriminator
		}

		partitionName = fmt.Sprintf("%s~%s~%s", partnerEvent.TenantGroupName, partnerEvent.Type, partnerEvent.Data.Code)
	case hr.DiscriminatorEmployee:
		receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

		if !ok {
			return "", envelopemessage.ErrInvalidReceivedEnvelope
		}

		employeeEvent, ok := receivedEnvelope.Message.(*hr.EmployeeEvent)

		if !ok {
			return "", pubsub.ErrInvalidDiscriminator
		}

		partitionName = fmt.Sprintf("%s~%s~%s", employeeEvent.TenantGroupName, employeeEvent.Type, employeeEvent.Data.Code)
	case hr.DiscriminatorPosition:
		receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

		if !ok {
			return "", envelopemessage.ErrInvalidReceivedEnvelope
		}

		positionEvent, ok := receivedEnvelope.Message.(*hr.PositionEvent)

		if !ok {
			return "", pubsub.ErrInvalidDiscriminator
		}

		partitionName = fmt.Sprintf("%s~%s~%s", positionEvent.TenantGroupName, positionEvent.Type, positionEvent.Data.Code)
	case hr.DiscriminatorRole:
		receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

		if !ok {
			return "", envelopemessage.ErrInvalidReceivedEnvelope
		}

		roleEvent, ok := receivedEnvelope.Message.(*hr.RoleEvent)

		if !ok {
			return "", pubsub.ErrInvalidDiscriminator
		}

		partitionName = fmt.Sprintf("%s~%s~%s", roleEvent.TenantGroupName, roleEvent.Type, roleEvent.Data.Code)
	default:
		partitionName = ""
	}

	return partitionName, nil
}
