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

	if err := publisher.sender.SendMessage(ctx, serviceBusMessage, nil); err != nil {
		return err
	}

	return nil
}

type UnmarshalDiscriminatorFunc func(serviceBusReceivedMessage *azservicebus.ReceivedMessage, discriminator *pubsub.Discriminator) error

type UnmarshalMessageFunc func(serviceBusReceivedMessage *azservicebus.ReceivedMessage, message pubsub.Message) error

type SubscriberOptions struct {
	Interval      time.Duration
	MessagesLimit int
}

type Subscriber struct {
	receiver                   *azservicebus.Receiver
	dispatcher                 *pubsub.Dispatcher
	unmarshalDiscriminatorFunc UnmarshalDiscriminatorFunc
	unmarshalMessageFunc       UnmarshalMessageFunc
	logger                     *slog.Logger
	options                    *SubscriberOptions
}

func NewSubscriber(receiver *azservicebus.Receiver, dispatcher *pubsub.Dispatcher, unmarshalDiscriminatorFunc UnmarshalDiscriminatorFunc, unmarshalMessageFunc UnmarshalMessageFunc, logger *slog.Logger, options *SubscriberOptions) *Subscriber {
	return &Subscriber{
		receiver:                   receiver,
		dispatcher:                 dispatcher,
		unmarshalDiscriminatorFunc: unmarshalDiscriminatorFunc,
		unmarshalMessageFunc:       unmarshalMessageFunc,
		logger:                     logger,
		options:                    options,
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
			serviceBusReceivedMessages, err := subscriber.receiver.ReceiveMessages(ctx, messagesLimit, nil)

			if err != nil {
				return err
			}

			for _, serviceBusReceivedMessage := range serviceBusReceivedMessages {
				discriminator := pubsub.DiscriminatorEmpty

				if err := subscriber.unmarshalDiscriminatorFunc(serviceBusReceivedMessage, &discriminator); err != nil {
					deadLetterOptions := &azservicebus.DeadLetterOptions{
						ErrorDescription: to.Ptr(err.Error()),
						Reason:           to.Ptr("UnmarshalDiscriminatorError"),
					}

					if err := subscriber.receiver.DeadLetterMessage(ctx, serviceBusReceivedMessage, deadLetterOptions); err != nil {
						var serviceBusErr *azservicebus.Error

						if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
							subscriber.logger.Warn("message lock lost", "discriminator", discriminator)

							continue
						}

						return err
					}

					subscriber.logger.Error("dead letter message", "discriminator", discriminator, "error", err)
				}

				if handler, ok := subscriber.dispatcher.Dispatch(discriminator); ok {
					message := handler.Create()

					if err := subscriber.unmarshalMessageFunc(serviceBusReceivedMessage, message); err != nil {
						deadLetterOptions := &azservicebus.DeadLetterOptions{
							ErrorDescription: to.Ptr(err.Error()),
							Reason:           to.Ptr("UnmarshalMessageError"),
						}

						if err := subscriber.receiver.DeadLetterMessage(ctx, serviceBusReceivedMessage, deadLetterOptions); err != nil {
							var serviceBusErr *azservicebus.Error

							if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
								subscriber.logger.Warn("message lock lost", "discriminator", discriminator)

								continue
							}

							return err
						}

						subscriber.logger.Error("dead letter message", "discriminator", discriminator, "error", err)
					}

					if err := handler.Handle(message); err != nil {
						if err := subscriber.receiver.AbandonMessage(ctx, serviceBusReceivedMessage, nil); err != nil {
							var serviceBusErr *azservicebus.Error

							if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
								subscriber.logger.Warn("message lock lost", "discriminator", discriminator)

								continue
							}

							return err
						}

						subscriber.logger.Error("abandon message", "discriminator", discriminator, "error", err)
					}
				} else {
					subscriber.logger.Info("message handler not found", "discriminator", discriminator)
				}

				if err := subscriber.receiver.CompleteMessage(ctx, serviceBusReceivedMessage, nil); err != nil {
					var serviceBusErr *azservicebus.Error

					if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
						subscriber.logger.Warn("message lock lost", "discriminator", discriminator)

						continue
					}

					return err
				}

				subscriber.logger.Info("complete message", "discriminator", discriminator)
			}
		}
	}
}
