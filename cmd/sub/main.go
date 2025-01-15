package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/scaleforce/synchronization-sdk-for-go/internal/handler/event/hr"
	"github.com/scaleforce/synchronization-sdk-for-go/internal/handler/event/masterdata"
	"github.com/scaleforce/synchronization-sdk-for-go/internal/handler/event/partner"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/azure/servicebus"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/pubsub"
	"github.com/spf13/viper"
)

var (
	credential *azidentity.DefaultAzureCredential
	client     *azservicebus.Client

	dispatcher *pubsub.Dispatcher
)

type partialMessage struct {
	Type string `json:"Type"`
}

func (message *partialMessage) Discriminator() pubsub.Discriminator {
	return pubsub.Discriminator(message.Type)
}

func init() {
	viper.AddConfigPath(".")
	// viper.SetConfigFile(".env")
	// viper.SetConfigName("config")
	// viper.SetConfigType("env") // "env", "json", "yaml"
	viper.SetEnvPrefix("sync")
	viper.AutomaticEnv()

	viper.SetDefault("AZURE_SERVICEBUS_INTERVAL", 10*time.Second)
	viper.SetDefault("AZURE_SERVICEBUS_MESSAGES_LIMIT", 1)

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	var err error

	credential, err = azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Panic(err)
	}

	_ = credential
	// client, err = azservicebus.NewClient(viper.GetString("AZURE_SERVICEBUS_NAMESPACE"), credential, nil)
	client, err = azservicebus.NewClientFromConnectionString(viper.GetString("AZURE_SERVICEBUS_CONNECTION_STRING"), nil)

	if err != nil {
		log.Panic(err)
	}

	dispatcher = pubsub.NewDispatcher()

	dispatcher.Register(&hr.PositionEventHandler{})
	dispatcher.Register(&hr.RoleEventHandler{})
	dispatcher.Register(&hr.EmployeeEventHandler{})

	dispatcher.Register(&masterdata.CityEventHandler{})
	dispatcher.Register(&masterdata.CircleEventHandler{})
	dispatcher.Register(&masterdata.ZoneEventHandler{})

	dispatcher.Register(&partner.PartnerGroupEventHandler{})
	dispatcher.Register(&partner.PartnerEventHandler{})
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

	subscriberOptions := &servicebus.SubscriberOptions{
		Interval: viper.GetDuration("AZURE_SERVICEBUS_INTERVAL"),
	}

	subscriber := servicebus.NewSubscriber(receiver, dispatcher, subscriberOptions)

	if err := subscriber.Run(ctx); err != nil {
		log.Panic(err)
	}
}
