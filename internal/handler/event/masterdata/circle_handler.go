package masterdata

import (
	"encoding/json"
	"fmt"
	"log/slog"

	envelopemessage "github.com/scaleforce/synchronization-for-go/internal/message/envelope"
	masterdataevent "github.com/scaleforce/synchronization-for-go/pkg/message/event/masterdata"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

type CircleEventHandler struct {
	logger *slog.Logger
}

func NewCircleEventHandler(logger *slog.Logger) *CircleEventHandler {
	return &CircleEventHandler{
		logger: logger,
	}
}

func (handler *CircleEventHandler) Discriminator() pubsub.Discriminator {
	return masterdataevent.DiscriminatorCircle
}

func (handler *CircleEventHandler) Handle(message pubsub.Message) error {
	receivedEnvelope, ok := message.(*envelopemessage.ReceivedEnvelope)

	if !ok {
		return envelopemessage.ErrInvalidReceivedEnvelope
	}

	circleEvent, ok := receivedEnvelope.Message.(*masterdataevent.CircleEvent)

	if !ok {
		return pubsub.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(circleEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
