package partner

import (
	"github.com/scaleforce/synchronization-for-go/pkg/message"
	"github.com/scaleforce/synchronization-for-go/pkg/message/event"
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

const (
	DiscriminatorPartnerGroup pubsub.Discriminator = "Partner_PartnerGroup"
	DiscriminatorPartner      pubsub.Discriminator = "Partner_Partner"
)

const (
	ScopeNone int = iota
	ScopeTenantGroup
	ScopePartnerGroup
)

type PartnerGroupData struct {
	Code               string  `json:"Code"`
	Name               string  `json:"Name"`
	TenantName         string  `json:"TenantName"`
	CreditLimit        float64 `json:"CreditLimit"`
	GroupOwnerCode     string  `json:"GroupOwnerCode"`
	GroupOwnerUserId   string  `json:"GroupOwnerUserId"`
	GroupOwnerUserName string  `json:"GroupOwnerUserName"`
	GroupOwnerScope    int     `json:"GroupOwnerScope"`
	CreatedByUserId    string  `json:"CreatedByUserId"`
	CreatedByUserName  string  `json:"CreatedByUserName"`
	ModifiedByUserId   string  `json:"ModifiedByUserId"`
	ModifiedByUserName string  `json:"ModifiedByUserName"`
}

type PartnerGroupEvent struct {
	event.TenantGroupEvent
	Data *PartnerGroupData `json:"Data"`
}

func NewPartnerGroupEvent(version, operation, timestamp, tenantGroupName string, data *PartnerGroupData) *PartnerGroupEvent {
	return &PartnerGroupEvent{
		TenantGroupEvent: event.TenantGroupEvent{
			Event: event.Event{
				Message: message.Message{
					Type: string(DiscriminatorPartnerGroup),
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

func (message *PartnerGroupEvent) Discriminator() pubsub.Discriminator {
	return DiscriminatorPartnerGroup
}

// Flags for Infrastructures
const InfrastructuresNone int = 0
const (
	InfrastructuresLAN int = 1 << iota
	InfrastructuresExcitelFiber
	InfrastructuresERPFiber
)

// Flags for Services
const ServicesNone int = 0
const (
	ServicesInternet int = 1 << iota
	ServicesCableTV
)

const (
	StatusNone int = iota
	StatusActive
	StatusInactive
	StatusTerminated
)

const (
	ContactTypeNone int = iota
	ContactTypeGeneral
	ContactTypeTechnical
	ContactTypeOwner
	ContactTypeBilling
	ContactTypeShipping
)

// Flags for PaymentTypes
const PaymentTypesNone int = 0
const (
	PaymentTypesCash int = 1 << iota
	PaymentTypesDigital
)

type EmployeeRelationModel struct {
	UserId       string `json:"UserId"`
	UserName     string `json:"UserName"`
	Code         string `json:"Code"`
	Scope        int    `json:"Scope"`
	PositionCode string `json:"PositionCode"`
	PositionName string `json:"PositionName"`
}

type PartnerTenantGroupEmployeeModel struct {
	EmployeeRelationModel
	ActingPositionCode string `json:"ActingPositionCode"`
	ActingPositionName string `json:"ActingPositionName"`
}

type ContactModel struct {
	ContactType    int     `json:"ContactType"`
	Name           string  `json:"Name"`
	PrimaryPhone   string  `json:"PrimaryPhone"`
	SecondaryPhone string  `json:"SecondaryPhone"`
	LandlinePhone  string  `json:"LandlinePhone"`
	Email          string  `json:"Email"`
	CityCode       string  `json:"CityCode"`
	CityName       string  `json:"CityName"`
	StateName      string  `json:"StateName"`
	SubareaName    string  `json:"SubareaName"`
	PostalCode     string  `json:"PostalCode"`
	Address        string  `json:"Address"`
	Latitude       float64 `json:"Latitude"`
	Longitude      float64 `json:"Longitude"`
}

type PartnerData struct {
	Code                      string                             `json:"Code"`
	Name                      string                             `json:"Name"`
	TenantName                string                             `json:"TenantName"`
	PartnerGroupCode          string                             `json:"PartnerGroupCode"`
	PartnerGroupName          string                             `json:"PartnerGroupName"`
	CreatedTime               string                             `json:"CreatedTime"`
	Nickname                  string                             `json:"Nickname"`
	Infrastructures           int                                `json:"Infrastructures"`
	Services                  int                                `json:"Services"`
	Status                    int                                `json:"Status"`
	Exclusive                 bool                               `json:"Exclusive"`
	ActiveForSales            bool                               `json:"ActiveForSales"`
	MigrationStartDate        string                             `json:"MigrationStartDate"`
	MigrationEndDate          string                             `json:"MigrationEndDate"`
	PlanBookCode              string                             `json:"PlanBookCode"`
	PlanBookName              string                             `json:"PlanBookName"`
	GSTIN                     string                             `json:"GSTIN"`
	PAN                       string                             `json:"PAN"`
	DefaultNASPort            string                             `json:"DefaultNASPort"`
	LegalEntity               string                             `json:"LegalEntity"`
	CityCode                  string                             `json:"CityCode"`
	CityName                  string                             `json:"CityName"`
	StateName                 string                             `json:"StateName"`
	CircleCode                string                             `json:"CircleCode"`
	CircleName                string                             `json:"CircleName"`
	ZoneCode                  string                             `json:"ZoneCode"`
	ZoneName                  string                             `json:"ZoneName"`
	Latitude                  float64                            `json:"Latitude"`
	Longitude                 float64                            `json:"Longitude"`
	UnlockedDevices           bool                               `json:"UnlockedDevices"`
	Homepass                  uint64                             `json:"Homepass"`
	BeneficiaryName           string                             `json:"BeneficiaryName"`
	BeneficiaryPhone          string                             `json:"BeneficiaryPhone"`
	BeneficiaryVPA            string                             `json:"BeneficiaryVPA"`
	BankAccountNumber         string                             `json:"BankAccountNumber"`
	BankIFSCCode              string                             `json:"BankIFSCCode"`
	ContactReferenceId        string                             `json:"ContactReferenceId"`
	TenantGroupEmployees      []*PartnerTenantGroupEmployeeModel `json:"TenantGroupEmployees"`
	Contacts                  []*ContactModel                    `json:"Contacts"`
	CreatedByUserId           string                             `json:"CreatedByUserId"`
	CreatedByUserName         string                             `json:"CreatedByUserName"`
	ModifiedByUserId          string                             `json:"ModifiedByUserId"`
	ModifiedByUserName        string                             `json:"ModifiedByUserName"`
	SubscriberPaymentTypes    int                                `json:"SubscriberPaymentTypes"`
	NewSubscriberPaymentTypes int                                `json:"NewSubscriberPaymentTypes"`
}

type PartnerEvent struct {
	event.TenantGroupEvent
	Data *PartnerData `json:"Data"`
}

func NewPartnerEvent(version, operation, timestamp, tenantGroupName string, data *PartnerData) *PartnerEvent {
	return &PartnerEvent{
		TenantGroupEvent: event.TenantGroupEvent{
			Event: event.Event{
				Message: message.Message{
					Type: string(DiscriminatorPartner),
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

func (message *PartnerEvent) Discriminator() pubsub.Discriminator {
	return DiscriminatorPartner
}
