package masterdata

import (
	"encoding/json"
	"fmt"
	"log/slog"

	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/masterdata"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsuberr"
)

type CityEventHandler struct {
	logger *slog.Logger
}

func NewCityEventHandler(logger *slog.Logger) *CityEventHandler {
	return &CityEventHandler{
		logger: logger,
	}
}

func (handler *CityEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorCity
}

func (handler *CityEventHandler) Create() pubsub.Message {
	return &event.CityEvent{}
}

func (handler *CityEventHandler) Handle(message pubsub.Message) error {
	cityEvent, ok := message.(*event.CityEvent)

	if !ok {
		return pubsuberr.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(cityEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
