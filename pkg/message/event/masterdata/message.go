package masterdata

import (
	"github.com/scaleforce/synchronization-for-go/pkg/message"
	"github.com/scaleforce/synchronization-for-go/pkg/message/event"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

const (
	DiscriminatorCity   pubsub.Discriminator = "MasterData_City"
	DiscriminatorCircle pubsub.Discriminator = "MasterData_Circle"
	DiscriminatorZone   pubsub.Discriminator = "MasterData_Zone"
)

type CityData struct {
	Code        string  `json:"Code"`
	Name        string  `json:"Name"`
	TenantName  string  `json:"TenantName"`
	StateName   string  `json:"StateName"`
	CoverageMap bool    `json:"CoverageMap"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
}

type CityEvent struct {
	event.TenantGroupEvent
	Data *CityData `json:"Data"`
}

func NewCityEvent(version, operation, timestamp, tenantGroupName string, data *CityData) *CityEvent {
	return &CityEvent{
		TenantGroupEvent: event.TenantGroupEvent{
			Event: event.Event{
				Message: message.Message{
					Type: string(DiscriminatorCity),
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

func (message *CityEvent) Discriminator() pubsub.Discriminator {
	return DiscriminatorCity
}

type CircleData struct {
	Code       string `json:"Code"`
	Name       string `json:"Name"`
	TenantName string `json:"TenantName"`
}

type CircleEvent struct {
	event.TenantGroupEvent
	Data *CircleData `json:"Data"`
}

func NewCircleEvent(version, operation, timestamp, tenantGroupName string, data *CircleData) *CircleEvent {
	return &CircleEvent{
		TenantGroupEvent: event.TenantGroupEvent{
			Event: event.Event{
				Message: message.Message{
					Type: string(DiscriminatorCircle),
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

func (message *CircleEvent) Discriminator() pubsub.Discriminator {
	return DiscriminatorCircle
}

type ZoneData struct {
	Code       string `json:"Code"`
	Name       string `json:"Name"`
	TenantName string `json:"TenantName"`
	CircleCode string `json:"CircleCode"`
	CircleName string `json:"CircleName"`
}

type ZoneEvent struct {
	event.TenantGroupEvent
	Data *ZoneData `json:"Data"`
}

func NewZoneEvent(version, operation, timestamp, tenantGroupName string, data *ZoneData) *ZoneEvent {
	return &ZoneEvent{
		TenantGroupEvent: event.TenantGroupEvent{
			Event: event.Event{
				Message: message.Message{
					Type: string(DiscriminatorZone),
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

func (message *ZoneEvent) Discriminator() pubsub.Discriminator {
	return DiscriminatorZone
}
