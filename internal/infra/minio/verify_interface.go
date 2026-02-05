package minio

import "goyavision/internal/app/port"

// Compile-time interface verification
var _ port.ObjectStorage = (*Client)(nil)
