package utility

import (
	"encoding/base64"
	"sync"
)

type IdGenerator interface {
	// Return base64 encoded string
	Generate() string
}

type mockIdGenerator struct {
	counter uint
	nodeId  string
	mu      sync.Mutex
}

func (mig *mockIdGenerator) Generate() string {
	mig.mu.Lock()
	mig.counter++
	mig.mu.Unlock()

	return base64.StdEncoding.EncodeToString([]byte(mig.nodeId + "-" + string(mig.counter)))
}

func NewMockIdGenerator(node string) IdGenerator {
	return &mockIdGenerator{
		nodeId:  node,
		counter: 0,
		mu:      sync.Mutex{},
	}
}
