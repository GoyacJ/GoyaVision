package ffmpeg

import "context"

type Pool struct{}

func NewPool(ffmpegBin string, maxRecord, maxFrame int) *Pool {
	return &Pool{}
}

func (p *Pool) AcquireRecordSlot(ctx context.Context) (release func(), err error) {
	return func() {}, nil
}

func (p *Pool) AcquireFrameSlot(ctx context.Context) (release func(), err error) {
	return func() {}, nil
}
