package auth

import "goyavision/internal/app/port"

// Compile-time interface verification
var _ port.TokenService = (*JWTService)(nil)
