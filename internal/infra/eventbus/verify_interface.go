package eventbus

import "goyavision/internal/app/port"

// Compile-time interface verification
var _ port.EventBus = (*LocalEventBus)(nil)
