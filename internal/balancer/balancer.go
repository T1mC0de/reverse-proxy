package balancer

import (
	"sync/atomic"
)

type Balancer interface {
	Pick(upstreams []string) *string
}

type RoundRobinBalancer struct {
	currentIndex atomic.Uint64
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{}
}

func (b *RoundRobinBalancer) Pick(upstreams []string) *string {
	if len(upstreams) == 0 {
		return nil
	}

	index := b.currentIndex.Add(1) - 1
	selectedIndex := int(index) % len(upstreams)
	return &upstreams[selectedIndex]
}