package preview

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNoPreviewSlot = errors.New("no available preview slot")
)

type Pool struct {
	maxPreview int
	mu         sync.Mutex
	count      int
}

func NewPool(maxPreview int) *Pool {
	if maxPreview <= 0 {
		maxPreview = 10
	}
	return &Pool{
		maxPreview: maxPreview,
	}
}

func (p *Pool) AcquirePreviewSlot(ctx context.Context) (release func(), err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.count >= p.maxPreview {
		return nil, ErrNoPreviewSlot
	}

	p.count++
	released := false

	release = func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		if !released && p.count > 0 {
			p.count--
			released = true
		}
	}

	go func() {
		<-ctx.Done()
		release()
	}()

	return release, nil
}

func (p *Pool) GetPreviewCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.count
}
