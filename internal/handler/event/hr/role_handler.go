package hr

import (
	"encoding/json"
	"fmt"
	"log/slog"

	event "github.com/scaleforce/synchronization-for-go/pkg/message/event/hr"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsuberr"
)

type RoleEventHandler struct {
	logger *slog.Logger
}

func NewRoleEventHandler(logger *slog.Logger) *RoleEventHandler {
	return &RoleEventHandler{
		logger: logger,
	}
}

func (handler *RoleEventHandler) Discriminator() pubsub.Discriminator {
	return event.DiscriminatorRole
}

func (handler *RoleEventHandler) Create() pubsub.Message {
	return &event.RoleEvent{}
}

func (handler *RoleEventHandler) Handle(message pubsub.Message) error {
	roleEvent, ok := message.(*event.RoleEvent)

	if !ok {
		return pubsuberr.ErrInvalidDiscriminator
	}

	data, err := json.MarshalIndent(roleEvent, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
