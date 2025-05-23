package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/scaleforce/synchronization-for-go/internal/azure/servicebus/util"
	"github.com/scaleforce/synchronization-for-go/internal/handler/event/hr"
	"github.com/scaleforce/synchronization-for-go/internal/handler/event/masterdata"
	"github.com/scaleforce/synchronization-for-go/internal/handler/event/partner"

	// "github.com/scaleforce/synchronization-for-go/pkg/azure/servicebus"
	partitionedservicebus "github.com/scaleforce/synchronization-for-go/pkg/azure/servicebus/partitioned"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
	"github.com/spf13/viper"
)

var (
	logger *slog.Logger

	credential *azidentity.DefaultAzureCredential
	client     *azservicebus.Client

	dispatcher *pubsub.Dispatcher
)

func init() {
	logger = slog.Default()
	// Use otelslog bridge to integrate with OpenTelemetry (https://pkg.go.dev/go.opentelemetry.io/otel/sdk/log)
	// logger := slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{AddSource: true}))
	// logger := slog.New(slog.NewJSONHandler(nil, &slog.HandlerOptions{AddSource: true}))

	viper.AddConfigPath(".")
	// viper.SetConfigFile(".env")
	// viper.SetConfigName("config")
	// viper.SetConfigType("env") // "env", "json", "yaml"
	viper.SetEnvPrefix("sync")
	viper.AutomaticEnv()

	viper.SetDefault("AZURE_SERVICEBUS_INTERVAL", 10*time.Second)
	viper.SetDefault("AZURE_SERVICEBUS_MESSAGES_LIMIT", 10)

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	var err error

	/*
		credential, err = azidentity.NewDefaultAzureCredential(nil)

		if err != nil {
			log.Panic(err)
		}

		client, err = azservicebus.NewClient(viper.GetString("AZURE_SERVICEBUS_NAMESPACE"), credential, nil)
	*/
	client, err = azservicebus.NewClientFromConnectionString(viper.GetString("AZURE_SERVICEBUS_CONNECTION_STRING"), nil)

	if err != nil {
		log.Panic(err)
	}

	dispatcher = pubsub.NewDispatcher()

	dispatcher.Register(hr.NewPositionEventHandler(logger))
	dispatcher.Register(hr.NewRoleEventHandler(logger))
	dispatcher.Register(hr.NewEmployeeEventHandler(logger))

	dispatcher.Register(masterdata.NewCityEventHandler(logger))
	dispatcher.Register(masterdata.NewCircleEventHandler(logger))
	dispatcher.Register(masterdata.NewZoneEventHandler(logger))

	dispatcher.Register(partner.NewPartnerGroupEventHandler(logger))
	dispatcher.Register(partner.NewPartnerEventHandler(logger))
}

func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancelCtx()

	receiverOptions := &azservicebus.ReceiverOptions{
		ReceiveMode: azservicebus.ReceiveModePeekLock,
	}

	receiver, err := client.NewReceiverForSubscription(viper.GetString("AZURE_SERVICEBUS_TOPIC"), viper.GetString("AZURE_SERVICEBUS_SUBSCRIPTION"), receiverOptions)

	if err != nil {
		log.Panic(err)
	}

	defer receiver.Close(ctx)

	subscriberOptions := &partitionedservicebus.SubscriberOptions{
		Interval:        viper.GetDuration("AZURE_SERVICEBUS_INTERVAL"),
		MessagesLimit:   viper.GetInt("AZURE_SERVICEBUS_MESSAGES_LIMIT"),
		PartitionsCount: viper.GetInt("AZURE_SERVICEBUS_PARTITIONS_COUNT"),
		PartitionsLimit: viper.GetInt("AZURE_SERVICEBUS_PARTITIONS_LIMIT"),
		PartitionsDrain: viper.GetBool("AZURE_SERVICEBUS_PARTITIONS_DRAIN"),
	}

	subscriber := partitionedservicebus.NewSubscriber(receiver, dispatcher, util.NewUnmarshalReceivedEnvelopeFunc(util.NewUnmarshalMessageFunc(util.CreateMessage)), util.GetPartitionName, logger, subscriberOptions)

	if err := subscriber.Run(ctx); err != nil {
		log.Panic(err)
	}
}
