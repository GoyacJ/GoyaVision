package eventbus

import (
	"context"
	"sync"

	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"
	"goyavision/pkg/logger"
)

// LocalEventBus 本地内存事件总线实现
type LocalEventBus struct {
	mu         sync.RWMutex
	handlers   map[string][]handlerWrapper
	bufferSize int
}

// handlerWrapper 处理器包装器，用于比较和删除
type handlerWrapper struct {
	id      int
	handler port.EventHandler
}

var handlerIDCounter = 0
var handlerIDMutex sync.Mutex

// NewLocalEventBus 创建本地事件总线实例
func NewLocalEventBus(bufferSize int) *LocalEventBus {
	if bufferSize <= 0 {
		bufferSize = 100
	}

	return &LocalEventBus{
		handlers:   make(map[string][]handlerWrapper),
		bufferSize: bufferSize,
	}
}

// Publish 发布事件
func (bus *LocalEventBus) Publish(ctx context.Context, event port.Event) error {
	if event == nil {
		return apperr.InvalidInput("event is required")
	}

	eventType := event.EventType()
	if eventType == "" {
		return apperr.InvalidInput("event type is required")
	}

	bus.mu.RLock()
	handlers := bus.handlers[eventType]
	if len(handlers) == 0 {
		bus.mu.RUnlock()
		logger.Debug("no handlers registered for event type", "event_type", eventType)
		return nil
	}

	handlersCopy := make([]handlerWrapper, len(handlers))
	copy(handlersCopy, handlers)
	bus.mu.RUnlock()

	for _, hw := range handlersCopy {
		go func(handler port.EventHandler) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("event handler panic", "panic", r, "event_type", eventType)
				}
			}()

			if err := handler(ctx, event); err != nil {
				logger.Error("event handler error", "error", err, "event_type", eventType)
			}
		}(hw.handler)
	}

	return nil
}

// Subscribe 订阅事件
func (bus *LocalEventBus) Subscribe(eventType string, handler port.EventHandler) {
	if eventType == "" {
		logger.Warn("attempted to subscribe with empty event type")
		return
	}
	if handler == nil {
		logger.Warn("attempted to subscribe with nil handler")
		return
	}

	handlerIDMutex.Lock()
	handlerIDCounter++
	id := handlerIDCounter
	handlerIDMutex.Unlock()

	hw := handlerWrapper{
		id:      id,
		handler: handler,
	}

	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.handlers[eventType] = append(bus.handlers[eventType], hw)
	logger.Debug("subscribed handler to event type", "handler_id", id, "event_type", eventType)
}

// Unsubscribe 取消订阅
func (bus *LocalEventBus) Unsubscribe(eventType string, handler port.EventHandler) {
	if eventType == "" {
		logger.Warn("attempted to unsubscribe with empty event type")
		return
	}
	if handler == nil {
		logger.Warn("attempted to unsubscribe with nil handler")
		return
	}

	bus.mu.Lock()
	defer bus.mu.Unlock()

	handlers, exists := bus.handlers[eventType]
	if !exists {
		logger.Debug("no handlers found for event type", "event_type", eventType)
		return
	}

	for i, hw := range handlers {
		if isSameHandler(hw.handler, handler) {
			bus.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			logger.Debug("unsubscribed handler from event type", "handler_id", hw.id, "event_type", eventType)

			if len(bus.handlers[eventType]) == 0 {
				delete(bus.handlers, eventType)
			}
			return
		}
	}

	logger.Debug("handler not found for event type", "event_type", eventType)
}

// GetSubscriberCount 获取指定事件类型的订阅者数量（用于测试和监控）
func (bus *LocalEventBus) GetSubscriberCount(eventType string) int {
	bus.mu.RLock()
	defer bus.mu.RUnlock()
	return len(bus.handlers[eventType])
}

// Clear 清空所有订阅（用于测试和重置）
func (bus *LocalEventBus) Clear() {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.handlers = make(map[string][]handlerWrapper)
	logger.Debug("event bus cleared")
}

// isSameHandler 比较两个 handler 是否相同
// 注意：Go 中函数不能直接比较，这里使用简单的实现
// 实际使用中，建议使用具名函数或结构体方法作为 handler，便于管理
func isSameHandler(h1, h2 port.EventHandler) bool {
	return &h1 == &h2
}
