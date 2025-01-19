package envelope

import (
	"errors"

	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

var (
	ErrInvalidEnvelope = errors.New("invalid envelope")
)

// Azure Service Bus infrastructure information like the message sequence number can be mapped from and to this envelope message.
// This information can be valuable when leveraging the full capabilities of Azure Service Bus to address advanced scenarios
// such as message ordering, deduplication, partitioning, sessions, request-reply patterns, and more.
type Envelope struct {
	SequenceNumber *int64
	Message        pubsub.Message
}

func NewEnvelope(message pubsub.Message) *Envelope {
	return &Envelope{
		Message: message,
	}
}

// Delegate the discriminator to the message contained in the envelope.
func (message *Envelope) Discriminator() pubsub.Discriminator {
	return message.Message.Discriminator()
}
