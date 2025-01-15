package masterdata

import (
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/message/event"
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/pubsub"
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

func (message *ZoneEvent) Discriminator() pubsub.Discriminator {
	return DiscriminatorZone
}
