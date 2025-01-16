package xnms

import (
	"github.com/scaleforce/synchronization-for-go/pkg/message"
	"github.com/scaleforce/synchronization-for-go/pkg/message/event"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

// Fictitious XNMS_Device event for demonstration purposes.

const (
	DiscriminatorDevice pubsub.Discriminator = "XNMS_Device"
)

const (
	StatusNone int = iota
	StatusOffline
	StatusOnline
)

type DeviceData struct {
	Code         string `json:"Code"`
	SerialNumber string `json:"SerialNumber"`
	TenantName   string `json:"TenantName"`
	Status       int    `json:"Status"`
}

type DeviceEvent struct {
	event.TenantGroupEvent
	Data *DeviceData `json:"Data"`
}

func NewDeviceEvent(version, operation, timestamp, tenantGroupName string, data *DeviceData) *DeviceEvent {
	return &DeviceEvent{
		TenantGroupEvent: event.TenantGroupEvent{
			Event: event.Event{
				Message: message.Message{
					Type: string(DiscriminatorDevice),
				},
				Version:   version,
				Operation: operation,
				Timestamp: timestamp,
			},
			TenantGroupName: tenantGroupName,
		},
		Data: data,
	}
}

func (message *DeviceEvent) Discriminator() pubsub.Discriminator {
	return DiscriminatorDevice
}
