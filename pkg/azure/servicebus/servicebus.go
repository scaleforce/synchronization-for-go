package servicebus

import (
	"context"
	"errors"
	"hash/fnv"
	"log/slog"
	"strings"
	"sync"
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

type UnmarshalMessageFunc func(serviceBusReceivedMessage *azservicebus.ReceivedMessage) (pubsub.Message, error)
type GetPartitionNameFunc func(message pubsub.Message) (string, error)

type SubscriberOptions struct {
	Interval        time.Duration
	MessagesLimit   int
	PartitionsCount int
	PartitionsDrain bool
}

type partitionMessage struct {
	message                   pubsub.Message
	serviceBusReceivedMessage *azservicebus.ReceivedMessage
}

type Subscriber struct {
	receiver             *azservicebus.Receiver
	dispatcher           *pubsub.Dispatcher
	unmarshalMessageFunc UnmarshalMessageFunc
	getPartitionNameFunc GetPartitionNameFunc
	logger               *slog.Logger
	options              *SubscriberOptions
}

func NewSubscriber(receiver *azservicebus.Receiver, dispatcher *pubsub.Dispatcher, unmarshalMessageFunc UnmarshalMessageFunc, getPartitionNameFunc GetPartitionNameFunc, logger *slog.Logger, options *SubscriberOptions) *Subscriber {
	return &Subscriber{receiver: receiver, dispatcher: dispatcher, unmarshalMessageFunc: unmarshalMessageFunc, getPartitionNameFunc: getPartitionNameFunc, logger: logger, options: options}
}

type RunError struct {
	ProducerErr  error
	ConsumerErrs []error
}

func (runErr *RunError) Error() string {
	errMsgs := make([]string, 0, 1+len(runErr.ConsumerErrs))

	if runErr.ProducerErr != nil {
		errMsgs = append(errMsgs, runErr.ProducerErr.Error())
	}

	for _, consumerErr := range runErr.ConsumerErrs {
		errMsgs = append(errMsgs, consumerErr.Error())
	}

	errMsg := strings.Join(errMsgs, "\n")

	return errMsg
}

func (subscriber *Subscriber) Run(ctx context.Context) error {
	partitionsCount := 1
	partitionsDrain := false

	if subscriber.options != nil {
		if subscriber.options.PartitionsCount > 0 {
			partitionsCount = subscriber.options.PartitionsCount
		}

		partitionsDrain = subscriber.options.PartitionsDrain
	}

	partitions := make([]chan *partitionMessage, 0, partitionsCount)

	for range partitionsCount {
		partition := make(chan *partitionMessage, 10)

		partitions = append(partitions, partition)
	}

	var consumerCtx context.Context

	if partitionsDrain {
		consumerCtx = context.Background()
	} else {
		consumerCtx = ctx
	}

	consumerErrs := make(chan error, len(partitions))

	consumerGroup := sync.WaitGroup{}

	consumerGroup.Add(len(partitions))

	for _, partition := range partitions {
		go func() {
			defer consumerGroup.Done()

			consumerErrs <- subscriber.consume(consumerCtx, partition)
		}()
	}

	producerErr := subscriber.produce(ctx, partitions)

	if partitionsDrain {
		for _, partition := range partitions {
			close(partition)
		}
	}

	consumerGroup.Wait()

	close(consumerErrs)

	errs := make([]error, 0, len(partitions))

	for consumerErr := range consumerErrs {
		if consumerErr != nil {
			errs = append(errs, consumerErr)
		}
	}

	if producerErr != nil || len(errs) != 0 {
		return &RunError{
			ProducerErr:  producerErr,
			ConsumerErrs: errs,
		}
	}

	return nil
}

var (
// errPartitionLimit = errors.New("partition limit")
)

func (subscriber *Subscriber) enqueue(ctx context.Context, partitions []chan *partitionMessage, partitionMessage *partitionMessage) error {
	if len(partitions) == 0 {
		return nil
	}

	partitionName, err := subscriber.getPartitionNameFunc(partitionMessage.message)

	if err != nil {
		return err
	}

	hash32 := fnv.New32a()

	hash32.Write([]byte(partitionName))

	partitionHash := hash32.Sum32()

	partitionIndex := int(partitionHash % uint32(len(partitions)))

	select {
	case <-ctx.Done():
		return ctx.Err()
	case partitions[partitionIndex] <- partitionMessage:
		// default:
		// 	return errPartitionLimit
	}

	return nil
}

func (subscriber *Subscriber) produce(ctx context.Context, partitions []chan *partitionMessage) error {
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

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.Tick(interval):
			serviceBusReceivedMessages, err := subscriber.receiver.ReceiveMessages(ctx, messagesLimit, nil)

			if err != nil {
				return err
			}

			for _, serviceBusReceivedMessage := range serviceBusReceivedMessages {
				message, err := subscriber.unmarshalMessageFunc(serviceBusReceivedMessage)

				if err != nil {
					deadLetterOptions := &azservicebus.DeadLetterOptions{
						ErrorDescription: to.Ptr(err.Error()),
						Reason:           to.Ptr("UnmarshalMessageError"),
					}

					if err := subscriber.receiver.DeadLetterMessage(ctx, serviceBusReceivedMessage, deadLetterOptions); err != nil {
						var serviceBusErr *azservicebus.Error

						if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
							subscriber.logger.Warn("message lock was lost while trying to dead letter the message")

							continue
						}

						return err
					}

					subscriber.logger.Error("message was dead lettered", "error", err)

					continue
				}

				partitionMessage := &partitionMessage{
					message:                   message,
					serviceBusReceivedMessage: serviceBusReceivedMessage,
				}

				if err := subscriber.enqueue(ctx, partitions, partitionMessage); err != nil {
					if err := subscriber.receiver.AbandonMessage(ctx, serviceBusReceivedMessage, nil); err != nil {
						var serviceBusErr *azservicebus.Error

						if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
							subscriber.logger.Warn("message lock was lost while trying to abandon the message")

							continue
						}

						return err
					}

					subscriber.logger.Error("message was abandoned", "error", err)

					continue
				}
			}
		}
	}
}

func (subscriber *Subscriber) consume(ctx context.Context, partition <-chan *partitionMessage) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case partitionMessage, ok := <-partition:
			if !ok {
				return nil
			}

			discriminator := partitionMessage.message.Discriminator()

			if handler, ok := subscriber.dispatcher.Dispatch(discriminator); ok {
				if err := handler.Handle(partitionMessage.message); err != nil {
					if err := subscriber.receiver.AbandonMessage(ctx, partitionMessage.serviceBusReceivedMessage, nil); err != nil {
						var serviceBusErr *azservicebus.Error

						if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
							subscriber.logger.Warn("message lock was lost while trying to abandon the message")

							continue
						}

						return err
					}

					subscriber.logger.Error("message was abandoned", "error", err)

					continue
				}
			} else {
				subscriber.logger.Info("message handler was not found", "discriminator", discriminator)
			}

			if err := subscriber.receiver.CompleteMessage(ctx, partitionMessage.serviceBusReceivedMessage, nil); err != nil {
				var serviceBusErr *azservicebus.Error

				if errors.As(err, &serviceBusErr) && serviceBusErr.Code == azservicebus.CodeLockLost {
					subscriber.logger.Warn("message lock was lost while trying to complete the message")

					continue
				}

				return err
			}
		}
	}
}
