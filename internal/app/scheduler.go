package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"goyavision/config"
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/domain"
	"goyavision/internal/port"
	"goyavision/pkg/ffmpeg"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Schedule struct {
	Start      string `json:"start"`
	End        string `json:"end"`
	DaysOfWeek []int  `json:"days_of_week"`
}

type Scheduler struct {
	scheduler gocron.Scheduler
	repo      port.Repository
	inference port.Inference
	manager   *ffmpeg.Manager
	mtxCli    *mediamtx.Client
	mtxCfg    config.MediaMTX
	basePath  string
	jobs      map[uuid.UUID]gocron.Job
	jobsMu    sync.RWMutex
}

func NewScheduler(
	repo port.Repository,
	inference port.Inference,
	manager *ffmpeg.Manager,
	mtxCli *mediamtx.Client,
	mtxCfg config.MediaMTX,
	basePath string,
) (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("create scheduler: %w", err)
	}

	return &Scheduler{
		scheduler: s,
		repo:      repo,
		inference: inference,
		manager:   manager,
		mtxCli:    mtxCli,
		mtxCfg:    mtxCfg,
		basePath:  basePath,
		jobs:      make(map[uuid.UUID]gocron.Job),
	}, nil
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.scheduler.Start()
	return s.loadAndScheduleBindings(ctx)
}

func (s *Scheduler) Stop() error {
	s.scheduler.Shutdown()
	return nil
}

func (s *Scheduler) loadAndScheduleBindings(ctx context.Context) error {
	streams, err := s.repo.ListStreams(ctx, nil)
	if err != nil {
		return fmt.Errorf("list streams: %w", err)
	}

	for _, stream := range streams {
		if !stream.Enabled {
			continue
		}

		bindings, err := s.repo.ListAlgorithmBindingsByStream(ctx, stream.ID)
		if err != nil {
			continue
		}

		for _, binding := range bindings {
			if binding.Enabled {
				s.ScheduleBinding(ctx, binding, stream)
			}
		}
	}

	return nil
}

func (s *Scheduler) ScheduleBinding(ctx context.Context, binding *domain.AlgorithmBinding, stream *domain.Stream) error {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	if _, exists := s.jobs[binding.ID]; exists {
		return nil
	}

	var schedule *Schedule
	if binding.Schedule != nil && len(binding.Schedule) > 0 {
		if err := json.Unmarshal(binding.Schedule, &schedule); err != nil {
			return fmt.Errorf("parse schedule: %w", err)
		}
	}

	job, err := s.createJob(binding, stream, schedule)
	if err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	s.jobs[binding.ID] = job
	return nil
}

func (s *Scheduler) UnscheduleBinding(bindingID uuid.UUID) error {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	job, exists := s.jobs[bindingID]
	if !exists {
		return nil
	}

	if err := s.scheduler.RemoveJob(job.ID()); err != nil {
		return fmt.Errorf("remove job: %w", err)
	}

	delete(s.jobs, bindingID)
	return nil
}

func (s *Scheduler) createJob(binding *domain.AlgorithmBinding, stream *domain.Stream, schedule *Schedule) (gocron.Job, error) {
	firstRun := time.Now().Add(time.Duration(binding.InitialDelaySec) * time.Second)

	var job gocron.Job
	var err error

	if schedule != nil {
		job, err = s.createScheduledJob(binding, stream, schedule, firstRun)
	} else {
		job, err = s.createIntervalJob(binding, stream, firstRun)
	}

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s *Scheduler) createIntervalJob(binding *domain.AlgorithmBinding, stream *domain.Stream, firstRun time.Time) (gocron.Job, error) {
	duration := time.Duration(binding.IntervalSec) * time.Second

	job, err := s.scheduler.NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(s.runInference, binding.ID, stream.ID),
		gocron.WithStartAt(gocron.WithStartDateTime(firstRun)),
	)
	if err != nil {
		return nil, fmt.Errorf("create interval job: %w", err)
	}

	return job, nil
}

func (s *Scheduler) createScheduledJob(binding *domain.AlgorithmBinding, stream *domain.Stream, schedule *Schedule, firstRun time.Time) (gocron.Job, error) {
	startTime, err := parseTime(schedule.Start)
	if err != nil {
		return nil, fmt.Errorf("parse start time: %w", err)
	}

	endTime, err := parseTime(schedule.End)
	if err != nil {
		return nil, fmt.Errorf("parse end time: %w", err)
	}

	duration := time.Duration(binding.IntervalSec) * time.Second

	job, err := s.scheduler.NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(s.runInferenceWithSchedule, binding.ID, stream.ID, startTime, endTime, schedule.DaysOfWeek),
		gocron.WithStartAt(gocron.WithStartDateTime(firstRun)),
	)
	if err != nil {
		return nil, fmt.Errorf("create scheduled job: %w", err)
	}

	return job, nil
}

func (s *Scheduler) runInference(bindingID, streamID uuid.UUID) {
	ctx := context.Background()
	s.executeInference(ctx, bindingID, streamID)
}

func (s *Scheduler) runInferenceWithSchedule(bindingID, streamID uuid.UUID, startTime, endTime time.Duration, daysOfWeek []int) {
	now := time.Now()
	currentTime := time.Duration(now.Hour())*time.Hour + time.Duration(now.Minute())*time.Minute + time.Duration(now.Second())*time.Second
	currentDay := int(now.Weekday())

	if currentTime < startTime || currentTime > endTime {
		return
	}

	if len(daysOfWeek) > 0 {
		dayMatch := false
		for _, day := range daysOfWeek {
			if day == currentDay {
				dayMatch = true
				break
			}
		}
		if !dayMatch {
			return
		}
	}

	ctx := context.Background()
	s.executeInference(ctx, bindingID, streamID)
}

func (s *Scheduler) executeInference(ctx context.Context, bindingID, streamID uuid.UUID) {
	binding, err := s.repo.GetAlgorithmBinding(ctx, bindingID)
	if err != nil {
		return
	}

	if !binding.Enabled {
		s.UnscheduleBinding(bindingID)
		return
	}

	stream, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		return
	}

	if !stream.Enabled {
		return
	}

	algorithm, err := s.repo.GetAlgorithm(ctx, binding.AlgorithmID)
	if err != nil {
		return
	}

	rtspURL := s.getStreamRTSPURL(stream)

	outputDir := filepath.Join(s.basePath, "frames", streamID.String())
	outputPath := filepath.Join(outputDir, fmt.Sprintf("frame_%s.jpg", uuid.New().String()))

	if err := s.manager.ExtractFrame(ctx, streamID.String(), rtspURL, outputPath); err != nil {
		return
	}

	frameData, err := os.ReadFile(outputPath)
	if err != nil {
		return
	}

	base64Data := encodeBase64(frameData)
	requestBody := buildInferenceRequest(algorithm, base64Data)

	startTime := time.Now()
	resp, err := s.inference.Post(ctx, algorithm.Endpoint, port.InferenceRequest{Body: requestBody})
	latencyMs := int(time.Since(startTime).Milliseconds())

	if err != nil {
		return
	}

	result := &domain.InferenceResult{
		ID:                 uuid.New(),
		AlgorithmBindingID: bindingID,
		StreamID:           streamID,
		Ts:                 time.Now(),
		FrameRef:           outputPath,
		Output:             resp.Body,
		LatencyMs:          &latencyMs,
	}

	s.repo.CreateInferenceResult(ctx, result)
}

// getStreamRTSPURL 获取流的 MediaMTX RTSP URL
func (s *Scheduler) getStreamRTSPURL(stream *domain.Stream) string {
	pathName := s.pathName(stream)
	return fmt.Sprintf("%s/%s", s.mtxCfg.RTSPAddress, pathName)
}

// pathName 生成 MediaMTX 路径名
func (s *Scheduler) pathName(stream *domain.Stream) string {
	name := strings.ReplaceAll(stream.Name, " ", "_")
	name = strings.ToLower(name)
	return name
}

func parseTime(timeStr string) (time.Duration, error) {
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return 0, err
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute()*int(time.Minute)) + time.Duration(t.Second()*int(time.Second)), nil
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func buildInferenceRequest(algorithm *domain.Algorithm, base64Data string) []byte {
	request := map[string]interface{}{
		"image": base64Data,
	}
	if algorithm.InputSpec != nil {
		var inputSpec map[string]interface{}
		if err := json.Unmarshal(algorithm.InputSpec, &inputSpec); err == nil {
			for k, v := range inputSpec {
				request[k] = v
			}
		}
	}
	body, _ := json.Marshal(request)
	return body
}
