package servicebus

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/pubsub"
	"github.com/spf13/viper"
)

type PublisherOptions struct{}

type Publisher struct {
	sender  *azservicebus.Sender
	options *PublisherOptions
}

func NewPublisher(sender *azservicebus.Sender, options *PublisherOptions) *Publisher {
	return &Publisher{
		sender:  sender,
		options: options,
	}
}

func (publisher *Publisher) Publish(ctx context.Context, message pubsub.Message) error {
	serviceBusMessage := &azservicebus.Message{
		Body: []byte{},
	}

	if err := publisher.sender.SendMessage(ctx, serviceBusMessage, nil); err != nil {
		return err
	}

	return nil
}

type SubscriberOptions struct {
	Interval time.Duration
}

type Subscriber struct {
	receiver   *azservicebus.Receiver
	dispatcher *pubsub.Dispatcher
	options    *SubscriberOptions
}

func NewSubscriber(receiver *azservicebus.Receiver, dispatcher *pubsub.Dispatcher, options *SubscriberOptions) *Subscriber {
	return &Subscriber{
		receiver:   receiver,
		dispatcher: dispatcher,
		options:    options,
	}
}

func (subscriber *Subscriber) Run(ctx context.Context) error {
	interval := 1 * time.Minute

	if subscriber.options != nil && subscriber.options.Interval > 0 {
		interval = subscriber.options.Interval
	}

	tick := time.Tick(interval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick:
			serviceBusReceivedMessages, err := subscriber.receiver.ReceiveMessages(ctx, viper.GetInt("AZURE_SERVICEBUS_MESSAGES_LIMIT"), nil)

			if err != nil {
				return err
			}

			for _, serviceBusReceivedMessage := range serviceBusReceivedMessages {
				discriminator, err := pubsub.UnmarshalDiscriminator(serviceBusReceivedMessage.Body)

				if err != nil {
					deadLetterOptions := &azservicebus.DeadLetterOptions{
						ErrorDescription: to.Ptr(err.Error()),
						Reason:           to.Ptr("UnmarshalDiscriminatorError"),
					}

					if err := subscriber.receiver.DeadLetterMessage(ctx, serviceBusReceivedMessage, deadLetterOptions); err != nil {
						var serviceBusErr *azservicebus.Error

						if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
							continue
						}

						return err
					}
				}

				if handler, ok := subscriber.dispatcher.Dispatch(discriminator); ok {
					message := handler.Create()

					if err := pubsub.UnmarshalMessage(serviceBusReceivedMessage.Body, message); err != nil {
						deadLetterOptions := &azservicebus.DeadLetterOptions{
							ErrorDescription: to.Ptr(err.Error()),
							Reason:           to.Ptr("UnmarshalMessageError"),
						}

						if err := subscriber.receiver.DeadLetterMessage(ctx, serviceBusReceivedMessage, deadLetterOptions); err != nil {
							var serviceBusErr *azservicebus.Error

							if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
								continue
							}

							return err
						}
					}

					if err := handler.Handle(message); err != nil {
						if err := subscriber.receiver.AbandonMessage(ctx, serviceBusReceivedMessage, nil); err != nil {
							var serviceBusErr *azservicebus.Error

							if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
								continue
							}

							return err
						}
					}
				}

				if err := subscriber.receiver.CompleteMessage(ctx, serviceBusReceivedMessage, nil); err != nil {
					var serviceBusErr *azservicebus.Error

					if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
						continue
					}

					return err
				}
			}
		}
	}
}
