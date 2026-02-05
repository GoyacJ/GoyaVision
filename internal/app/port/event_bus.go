package port

import "context"

// EventBus 事件总线接口
//
// 职责：
//  1. 发布领域事件
//  2. 订阅领域事件
//  3. 解耦业务逻辑
//
// 实现：
//  - infra/eventbus/local.go (本地内存实现)
//  - 未来可扩展：Redis Pub/Sub, Kafka, RabbitMQ
//
// 使用场景：
//  - 媒体源创建后，触发资产索引
//  - 任务完成后，发送通知
//  - 工作流状态变更，记录审计日志
type EventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event Event) error

	// Subscribe 订阅事件（按事件类型）
	// handler 在新 goroutine 中异步执行
	Subscribe(eventType string, handler EventHandler)

	// Unsubscribe 取消订阅
	Unsubscribe(eventType string, handler EventHandler)
}

// Event 领域事件接口
type Event interface {
	// EventType 事件类型（用于路由）
	EventType() string

	// OccurredAt 事件发生时间
	OccurredAt() int64
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event Event) error
