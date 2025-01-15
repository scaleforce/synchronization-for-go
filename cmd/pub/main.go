package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/azure/servicebus"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/message/event"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/message/event/xnms"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/pubsub"
	"github.com/spf13/viper"
)

const (
	TenantGroupNameExcitel string = "excitel"
)

const (
	TenantNameDelhi     string = "delhi"
	TenantNameHyderabad string = "hyderabad"
	TenantNameMumbai    string = "mumbai"
)

var (
	credential *azidentity.DefaultAzureCredential
	client     *azservicebus.Client

	dispatcher *pubsub.Dispatcher
)

func init() {
	viper.AddConfigPath(".")
	// viper.SetConfigFile(".env")
	// viper.SetConfigName("config")
	// viper.SetConfigType("env") // "env", "json", "yaml"
	viper.SetEnvPrefix("sync")
	viper.AutomaticEnv()

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
}

func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancelCtx()

	sender, err := client.NewSender(viper.GetString("AZURE_SERVICEBUS_TOPIC"), nil)

	if err != nil {
		log.Panic(err)
	}

	defer sender.Close(ctx)

	publisher := servicebus.NewPublisher(sender, nil)

	tick := time.Tick(10 * time.Second)

	for done := false; !done; {
		select {
		case <-ctx.Done():
			done = true
		case <-tick:
			message := xnms.NewDeviceEvent(event.Version1, event.OperationAddOrSet, time.Now().UTC().Format(time.RFC3339), TenantGroupNameExcitel,
				&xnms.DeviceData{
					Code:         "123",
					SerialNumber: "123",
					TenantName:   TenantNameDelhi,
				},
			)

			if err := publisher.Publish(ctx, message); err != nil {
				log.Panic(err)
			}
		}
	}
}
