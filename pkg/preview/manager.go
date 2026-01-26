package preview

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"goyavision/pkg/ffmpeg"
)

type Provider string

const (
	ProviderMediaMTX Provider = "mediamtx"
	ProviderFFmpeg   Provider = "ffmpeg"
)

type Manager struct {
	provider    Provider
	mediamtxBin string
	ffmpegPool  *ffmpeg.Pool
	previewPool *Pool
	hlsBase     string
	hlsPath     string
	tasks       map[string]*PreviewTask
	tasksMu     sync.RWMutex
}

func NewManager(provider string, mediamtxBin string, ffmpegPool *ffmpeg.Pool, previewPool *Pool, hlsBase, hlsPath string) *Manager {
	p := Provider(provider)
	if p != ProviderMediaMTX && p != ProviderFFmpeg {
		p = ProviderFFmpeg
	}

	return &Manager{
		provider:    p,
		mediamtxBin: mediamtxBin,
		ffmpegPool:  ffmpegPool,
		previewPool: previewPool,
		hlsBase:     hlsBase,
		hlsPath:     hlsPath,
		tasks:       make(map[string]*PreviewTask),
	}
}

type PreviewTask struct {
	StreamID string
	RTSPURL  string
	HLSURL   string
	cmd      *exec.Cmd
	ctx      context.Context
	cancel   context.CancelFunc
	release  func()
}

func (m *Manager) StartPreview(ctx context.Context, streamID, rtspURL string) (*PreviewTask, error) {
	m.tasksMu.Lock()
	if _, exists := m.tasks[streamID]; exists {
		m.tasksMu.Unlock()
		return nil, fmt.Errorf("preview already running for stream %s", streamID)
	}
	m.tasksMu.Unlock()

	release, err := m.previewPool.AcquirePreviewSlot(ctx)
	if err != nil {
		return nil, err
	}

	var task *PreviewTask

	switch m.provider {
	case ProviderMediaMTX:
		task, err = m.startMediaMTX(ctx, streamID, rtspURL)
	case ProviderFFmpeg:
		task, err = m.startFFmpeg(ctx, streamID, rtspURL)
	default:
		release()
		return nil, fmt.Errorf("unsupported provider: %s", m.provider)
	}

	if err != nil {
		release()
		return nil, err
	}

	task.release = release

	m.tasksMu.Lock()
	m.tasks[streamID] = task
	m.tasksMu.Unlock()

	return task, nil
}

func (m *Manager) StopPreview(streamID string) error {
	m.tasksMu.Lock()
	defer m.tasksMu.Unlock()

	task, exists := m.tasks[streamID]
	if !exists {
		return fmt.Errorf("preview not running for stream %s", streamID)
	}

	if err := task.Stop(); err != nil {
		return err
	}

	delete(m.tasks, streamID)
	return nil
}

func (m *Manager) GetPreview(streamID string) (*PreviewTask, bool) {
	m.tasksMu.RLock()
	defer m.tasksMu.RUnlock()
	task, exists := m.tasks[streamID]
	return task, exists
}

func (m *Manager) startMediaMTX(ctx context.Context, streamID, rtspURL string) (*PreviewTask, error) {
	hlsURL := fmt.Sprintf("%s/%s", m.hlsBase, streamID)

	taskCtx, cancel := context.WithCancel(ctx)
	cmd := exec.CommandContext(taskCtx, m.mediamtxBin)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("start mediamtx: %w", err)
	}

	task := &PreviewTask{
		StreamID: streamID,
		RTSPURL:  rtspURL,
		HLSURL:   hlsURL,
		cmd:      cmd,
		ctx:      taskCtx,
		cancel:   cancel,
	}

	go func() {
		cmd.Wait()
		m.tasksMu.Lock()
		delete(m.tasks, streamID)
		m.tasksMu.Unlock()
	}()

	return task, nil
}

func (m *Manager) startFFmpeg(ctx context.Context, streamID, rtspURL string) (*PreviewTask, error) {
	ffmpegRelease, err := m.ffmpegPool.AcquireFrameSlot(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire ffmpeg slot: %w", err)
	}

	hlsPath := filepath.Join(m.hlsPath, streamID)
	if err := os.MkdirAll(hlsPath, 0755); err != nil {
		ffmpegRelease()
		return nil, fmt.Errorf("create hls directory: %w", err)
	}

	hlsURL := fmt.Sprintf("%s/%s/index.m3u8", m.hlsBase, streamID)

	args := []string{
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-f", "hls",
		"-hls_time", "2",
		"-hls_list_size", "3",
		"-hls_flags", "delete_segments",
		"-hls_segment_filename", filepath.Join(hlsPath, "segment_%03d.ts"),
		filepath.Join(hlsPath, "index.m3u8"),
	}

	taskCtx, cancel := context.WithCancel(ctx)
	cmd := exec.CommandContext(taskCtx, m.ffmpegPool.GetFFmpegBin(), args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		cancel()
		ffmpegRelease()
		return nil, fmt.Errorf("start ffmpeg: %w", err)
	}

	task := &PreviewTask{
		StreamID: streamID,
		RTSPURL:  rtspURL,
		HLSURL:   hlsURL,
		cmd:      cmd,
		ctx:      taskCtx,
		cancel:   cancel,
	}

	go func() {
		cmd.Wait()
		ffmpegRelease()
		if task.release != nil {
			task.release()
		}
		m.tasksMu.Lock()
		delete(m.tasks, streamID)
		m.tasksMu.Unlock()
	}()

	return task, nil
}

func (t *PreviewTask) Stop() error {
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

func (t *PreviewTask) IsRunning() bool {
	if t.cmd == nil || t.cmd.Process == nil {
		return false
	}
	if t.cmd.ProcessState != nil {
		return !t.cmd.ProcessState.Exited()
	}
	return t.ctx.Err() == nil
}
