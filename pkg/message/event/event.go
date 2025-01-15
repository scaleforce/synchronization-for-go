package event

import (
	"github.com/scaleforce/synchronization-sdk-for-go/pkg/message"
)

const (
	Version1 string = "1"
)

const (
	OperationAdd      string = "Add"
	OperationAddOrSet string = "AddOrSet"
	OperationRemove   string = "Remove"
)

type Event struct {
	message.Message
	Version   string `json:"Version"`
	Operation string `json:"Operation"`
	Timestamp string `json:"Timestamp"`
}

const (
	TenantGroupNameExcitel string = "excitel"
)

type TenantGroupEvent struct {
	Event
	TenantGroupName string `json:"TenantGroupName"`
}

const (
	TenantNameDelhi     string = "delhi"
	TenantNameHyderabad string = "hyderabad"
	TenantNameMumbai    string = "mumbai"
)

type TenantEvent struct {
	TenantGroupEvent
	TenantName string `json:"TenantName"`
}
