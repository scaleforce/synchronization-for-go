package servicebus

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

type MarshalMessageFunc func(message pubsub.Message) (*azservicebus.Message, error)

type PublisherOptions struct{}

type Publisher struct {
	sender             *azservicebus.Sender
	marshalMessageFunc MarshalMessageFunc
	logger             *slog.Logger
	options            *PublisherOptions
}

func NewPublisher(sender *azservicebus.Sender, marshalMessageFunc MarshalMessageFunc, logger *slog.Logger, options *PublisherOptions) *Publisher {
	return &Publisher{
		sender:             sender,
		marshalMessageFunc: marshalMessageFunc,
		logger:             logger,
		options:            options,
	}
}

func (publisher *Publisher) Publish(ctx context.Context, message pubsub.Message) error {
	serviceBusMessage, err := publisher.marshalMessageFunc(message)

	if err != nil {
		return err
	}

	publisher.logger.Debug("send message")

	if err := publisher.sender.SendMessage(ctx, serviceBusMessage, nil); err != nil {
		return err
	}

	return nil
}

type UnmarshalMessageFunc func(serviceBusReceivedMessage *azservicebus.ReceivedMessage) (pubsub.Message, error)

type SubscriberOptions struct {
	Interval      time.Duration
	MessagesLimit int
}

type Subscriber struct {
	receiver             *azservicebus.Receiver
	dispatcher           *pubsub.Dispatcher
	unmarshalMessageFunc UnmarshalMessageFunc
	logger               *slog.Logger
	options              *SubscriberOptions
}

func NewSubscriber(receiver *azservicebus.Receiver, dispatcher *pubsub.Dispatcher, unmarshalMessageFunc UnmarshalMessageFunc, logger *slog.Logger, options *SubscriberOptions) *Subscriber {
	return &Subscriber{
		receiver:             receiver,
		dispatcher:           dispatcher,
		unmarshalMessageFunc: unmarshalMessageFunc,
		logger:               logger,
		options:              options,
	}
}

func (subscriber *Subscriber) Run(ctx context.Context) error {
	interval := 1 * time.Minute
	messagesLimit := 1

	if subscriber.options != nil {
		if subscriber.options.Interval > 0 {
			interval = subscriber.options.Interval
		}

		if subscriber.options.MessagesLimit > 0 {
			messagesLimit = subscriber.options.MessagesLimit
		}
	}

	tick := time.Tick(interval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick:
			subscriber.logger.Debug("receive messages")

			serviceBusReceivedMessages, err := subscriber.receiver.ReceiveMessages(ctx, messagesLimit, nil)

			if err != nil {
				return err
			}

			for _, serviceBusReceivedMessage := range serviceBusReceivedMessages {
				subscriber.logger.Debug("unmarshal message")

				message, err := subscriber.unmarshalMessageFunc(serviceBusReceivedMessage)

				if err != nil {
					subscriber.logger.Error("dead letter message", "error", err)

					deadLetterOptions := &azservicebus.DeadLetterOptions{
						ErrorDescription: to.Ptr(err.Error()),
						Reason:           to.Ptr("UnmarshalMessageError"),
					}

					if err := subscriber.receiver.DeadLetterMessage(ctx, serviceBusReceivedMessage, deadLetterOptions); err != nil {
						var serviceBusErr *azservicebus.Error

						if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
							subscriber.logger.Warn("message lock lost when attempting to dead letter the message")

							continue
						}

						return err
					}
				}

				discriminator := message.Discriminator()

				if handler, ok := subscriber.dispatcher.Dispatch(discriminator); ok {
					subscriber.logger.Debug("handle message")

					if err := handler.Handle(message); err != nil {
						subscriber.logger.Error("abandon message", "error", err)

						if err := subscriber.receiver.AbandonMessage(ctx, serviceBusReceivedMessage, nil); err != nil {
							var serviceBusErr *azservicebus.Error

							if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
								subscriber.logger.Warn("message lock lost when attempting to abandon the message")

								continue
							}

							return err
						}
					}
				} else {
					subscriber.logger.Info("message handler not found", "discriminator", discriminator)
				}

				subscriber.logger.Debug("complete message")

				if err := subscriber.receiver.CompleteMessage(ctx, serviceBusReceivedMessage, nil); err != nil {
					var serviceBusErr *azservicebus.Error

					if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
						subscriber.logger.Warn("message lock lost when attempting to complete the message")

						continue
					}

					return err
				}
			}
		}
	}
}
