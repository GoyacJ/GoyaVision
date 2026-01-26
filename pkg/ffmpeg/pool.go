package ffmpeg

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNoRecordSlot = errors.New("no available record slot")
	ErrNoFrameSlot  = errors.New("no available frame slot")
)

type Pool struct {
	ffmpegBin string
	maxRecord int
	maxFrame  int

	recordMu    sync.Mutex
	recordCount int

	frameMu    sync.Mutex
	frameCount int
}

func NewPool(ffmpegBin string, maxRecord, maxFrame int) *Pool {
	if maxRecord <= 0 {
		maxRecord = 16
	}
	if maxFrame <= 0 {
		maxFrame = 16
	}
	return &Pool{
		ffmpegBin: ffmpegBin,
		maxRecord: maxRecord,
		maxFrame:  maxFrame,
	}
}

func (p *Pool) AcquireRecordSlot(ctx context.Context) (release func(), err error) {
	p.recordMu.Lock()
	defer p.recordMu.Unlock()

	if p.recordCount >= p.maxRecord {
		return nil, ErrNoRecordSlot
	}

	p.recordCount++
	released := false

	release = func() {
		p.recordMu.Lock()
		defer p.recordMu.Unlock()
		if !released && p.recordCount > 0 {
			p.recordCount--
			released = true
		}
	}

	go func() {
		<-ctx.Done()
		release()
	}()

	return release, nil
}

func (p *Pool) AcquireFrameSlot(ctx context.Context) (release func(), err error) {
	p.frameMu.Lock()
	defer p.frameMu.Unlock()

	if p.frameCount >= p.maxFrame {
		return nil, ErrNoFrameSlot
	}

	p.frameCount++
	released := false

	release = func() {
		p.frameMu.Lock()
		defer p.frameMu.Unlock()
		if !released && p.frameCount > 0 {
			p.frameCount--
			released = true
		}
	}

	go func() {
		<-ctx.Done()
		release()
	}()

	return release, nil
}

func (p *Pool) GetRecordCount() int {
	p.recordMu.Lock()
	defer p.recordMu.Unlock()
	return p.recordCount
}

func (p *Pool) GetFrameCount() int {
	p.frameMu.Lock()
	defer p.frameMu.Unlock()
	return p.frameCount
}

func (p *Pool) GetFFmpegBin() string {
	return p.ffmpegBin
}
