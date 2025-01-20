package empty

import (
	"github.com/scaleforce/synchronization-for-go/pkg/pubsub"
)

type Empty struct{}

func (message *Empty) Discriminator() pubsub.Discriminator {
	return pubsub.DiscriminatorEmpty
}
