package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Manager struct {
	pool     *Pool
	basePath string
}

func NewManager(pool *Pool, basePath string) *Manager {
	return &Manager{
		pool:     pool,
		basePath: basePath,
	}
}

type RecordTask struct {
	StreamID   string
	RTSPURL    string
	OutputDir  string
	SegmentSec int
	cmd        *exec.Cmd
	ctx        context.Context
	cancel     context.CancelFunc
	release    func()
}

type FrameTask struct {
	StreamID    string
	RTSPURL     string
	OutputPath  string
	IntervalSec int
	cmd         *exec.Cmd
	ctx         context.Context
	cancel      context.CancelFunc
	release     func()
}

func (m *Manager) StartRecord(ctx context.Context, streamID, rtspURL string, segmentSec int) (*RecordTask, error) {
	release, err := m.pool.AcquireRecordSlot(ctx)
	if err != nil {
		return nil, err
	}

	outputDir := filepath.Join(m.basePath, streamID)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		release()
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	taskCtx, cancel := context.WithCancel(ctx)
	outputPattern := filepath.Join(outputDir, "segment_%03d.mp4")

	args := []string{
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-c", "copy",
		"-f", "segment",
		"-segment_time", fmt.Sprintf("%d", segmentSec),
		"-segment_format", "mp4",
		"-reset_timestamps", "1",
		"-strftime", "1",
		"-segment_atclocktime", "1",
		outputPattern,
	}

	cmd := exec.CommandContext(taskCtx, m.pool.ffmpegBin, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		cancel()
		release()
		return nil, fmt.Errorf("start ffmpeg: %w", err)
	}

	task := &RecordTask{
		StreamID:   streamID,
		RTSPURL:    rtspURL,
		OutputDir:  outputDir,
		SegmentSec: segmentSec,
		cmd:        cmd,
		ctx:        taskCtx,
		cancel:     cancel,
		release:    release,
	}

	go func() {
		cmd.Wait()
		release()
	}()

	return task, nil
}

func (t *RecordTask) Stop() error {
	if t.cancel != nil {
		t.cancel()
	}
	if t.cmd != nil && t.cmd.Process != nil {
		t.cmd.Process.Kill()
	}
	if t.release != nil {
		t.release()
	}
	return nil
}

func (m *Manager) ExtractFrame(ctx context.Context, streamID, rtspURL, outputPath string) error {
	release, err := m.pool.AcquireFrameSlot(ctx)
	if err != nil {
		return err
	}
	defer release()

	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	args := []string{
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-frames:v", "1",
		"-q:v", "2",
		"-y",
		outputPath,
	}

	cmd := exec.CommandContext(ctx, m.pool.ffmpegBin, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extract frame: %w", err)
	}

	return nil
}

func (m *Manager) StartFrameExtraction(ctx context.Context, streamID, rtspURL, outputDir string, intervalSec int) (*FrameTask, error) {
	release, err := m.pool.AcquireFrameSlot(ctx)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		release()
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	taskCtx, cancel := context.WithCancel(ctx)
	outputPattern := filepath.Join(outputDir, "frame_%Y%m%d_%H%M%S_%03d.jpg")

	args := []string{
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-vf", fmt.Sprintf("fps=1/%d", intervalSec),
		"-q:v", "2",
		"-strftime", "1",
		outputPattern,
	}

	cmd := exec.CommandContext(taskCtx, m.pool.ffmpegBin, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		cancel()
		release()
		return nil, fmt.Errorf("start ffmpeg: %w", err)
	}

	task := &FrameTask{
		StreamID:    streamID,
		RTSPURL:     rtspURL,
		OutputPath:  outputDir,
		IntervalSec: intervalSec,
		cmd:         cmd,
		ctx:         taskCtx,
		cancel:      cancel,
		release:     release,
	}

	go func() {
		cmd.Wait()
		release()
	}()

	return task, nil
}

func (t *FrameTask) Stop() error {
	if t.cancel != nil {
		t.cancel()
	}
	if t.cmd != nil && t.cmd.Process != nil {
		t.cmd.Process.Kill()
	}
	if t.release != nil {
		t.release()
	}
	return nil
}

func (t *RecordTask) IsRunning() bool {
	if t.cmd == nil || t.cmd.Process == nil {
		return false
	}
	if t.cmd.ProcessState != nil {
		return !t.cmd.ProcessState.Exited()
	}
	return t.ctx.Err() == nil
}

func (t *FrameTask) IsRunning() bool {
	if t.cmd == nil || t.cmd.Process == nil {
		return false
	}
	if t.cmd.ProcessState != nil {
		return !t.cmd.ProcessState.Exited()
	}
	return t.ctx.Err() == nil
}

func (m *Manager) GetPool() *Pool {
	return m.pool
}
