package message

type Message struct {
	Type string `json:"Type"`
}

// Azure Service Bus infrastructure information can be mapped from and to the envelope message.
// This information can be valuable when leveraging the full capabilities of Azure Service Bus to address advanced scenarios such as message ordering, deduplication, partitioning, sessions, request-reply patterns, and more.
type EnvelopeMessage struct {
	SequenceNumber *int64
}
