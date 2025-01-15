package event

import (
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/message"
)

type Event struct {
	message.Message
	Version   string `json:"Version"`
	Timestamp string `json:"Timestamp"`
	Operation string `json:"Operation"`
}

type TenantGroupEvent struct {
	Event
	TenantGroupName string `json:"TenantGroupName"`
}

type TenantEvent struct {
	TenantGroupEvent
	TenantName string `json:"TenantName"`
}
