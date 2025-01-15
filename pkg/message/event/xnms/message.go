package xnms

import (
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/message/event"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/pubsub"
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

func (message *DeviceEvent) Discriminator() pubsub.Discriminator {
	return DiscriminatorDevice
}
