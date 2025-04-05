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
	Interval                 time.Duration
	MessagesLimit            int
	PartitionsCount          int
	ConsumersRunToCompletion bool
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
	partitions           []chan *partitionMessage
}

func NewSubscriber(receiver *azservicebus.Receiver, dispatcher *pubsub.Dispatcher, unmarshalMessageFunc UnmarshalMessageFunc, getPartitionNameFunc GetPartitionNameFunc, logger *slog.Logger, options *SubscriberOptions) *Subscriber {
	partitionsCount := 1

	if options != nil {
		if options.PartitionsCount > 0 {
			partitionsCount = options.PartitionsCount
		}
	}

	partitions := make([]chan *partitionMessage, partitionsCount)

	for partitionIndex := range partitions {
		partitions[partitionIndex] = make(chan *partitionMessage, 10)
	}

	return &Subscriber{
		receiver:             receiver,
		dispatcher:           dispatcher,
		unmarshalMessageFunc: unmarshalMessageFunc,
		getPartitionNameFunc: getPartitionNameFunc,
		logger:               logger,
		options:              options,
		partitions:           partitions,
	}
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
	consumersRunToCompletion := false

	if subscriber.options != nil {
		consumersRunToCompletion = subscriber.options.ConsumersRunToCompletion
	}

	var producerErr error

	producerGroup := sync.WaitGroup{}

	producerGroup.Add(1)

	go func() {
		defer producerGroup.Done()

		producerErr = subscriber.produce(ctx)
	}()

	var consumeCtx context.Context

	if consumersRunToCompletion {
		consumeCtx = context.Background()
	} else {
		consumeCtx = ctx
	}

	consumerErrs := []error{}

	consumersGroup := sync.WaitGroup{}

	consumersGroup.Add(len(subscriber.partitions))

	for _, partition := range subscriber.partitions {
		go func() {
			defer consumersGroup.Done()

			if err := subscriber.consume(consumeCtx, partition); err != nil {
				consumerErrs = append(consumerErrs, err)
			}
		}()
	}

	producerGroup.Wait()

	if consumersRunToCompletion {
		for _, partition := range subscriber.partitions {
			close(partition)
		}
	}

	consumersGroup.Wait()

	if producerErr != nil || len(consumerErrs) != 0 {
		return &RunError{
			ProducerErr:  producerErr,
			ConsumerErrs: consumerErrs,
		}
	}

	return nil
}

var (
// errPartitionLimit = errors.New("partition limit")
)

func (subscriber *Subscriber) enqueue(ctx context.Context, partitionMessage *partitionMessage) error {
	partitionName, err := subscriber.getPartitionNameFunc(partitionMessage.message)

	if err != nil {
		return err
	}

	hash32 := fnv.New32a()

	hash32.Write([]byte(partitionName))

	partitionHash := hash32.Sum32()

	partitionIndex := int(partitionHash % uint32(len(subscriber.partitions)))

	select {
	case <-ctx.Done():
		return ctx.Err()
	case subscriber.partitions[partitionIndex] <- partitionMessage:
		// default:
		// 	return errPartitionLimit
	}

	return nil
}

func (subscriber *Subscriber) produce(ctx context.Context) error {
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

				if err := subscriber.enqueue(ctx, partitionMessage); err != nil {
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
