package pubsub

import (
	"context"
	"encoding/json"
)

const (
	DiscriminatorEmpty Discriminator = ""
)

type Discriminator string

type Message interface {
	Discriminator() Discriminator
}

type Publisher interface {
	Publish(message Message) error
}

type Handler interface {
	Discriminator() Discriminator
	Create() Message
	Handle(message Message) error
}

type Dispatcher struct {
	handlers map[Discriminator]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: map[Discriminator]Handler{},
	}
}

func (dispatcher *Dispatcher) Register(handler Handler) {
	discriminator := handler.Discriminator()

	dispatcher.handlers[discriminator] = handler
}

func (dispatcher *Dispatcher) Unregister(discriminator Discriminator) {
	delete(dispatcher.handlers, discriminator)
}

func (dispatcher *Dispatcher) Dispatch(discriminator Discriminator) (Handler, bool) {
	handler, ok := dispatcher.handlers[discriminator]

	return handler, ok
}

func UnmarshalDiscriminator(data []byte) (Discriminator, error) {
	message := &struct {
		Discriminator string `json:"discriminator"`
	}{}

	if err := json.Unmarshal(data, &message); err != nil {
		return DiscriminatorEmpty, err
	}

	return Discriminator(message.Discriminator), nil
}

func UnmarshalMessage(data []byte, message Message) error {
	return json.Unmarshal(data, message)
}

type Subscriber interface {
	Run(ctx context.Context) error
}
