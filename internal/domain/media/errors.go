package media

import "errors"

var (
	ErrSourceNameRequired    = errors.New("source name is required")
	ErrInvalidSourceType     = errors.New("invalid source type, must be 'pull' or 'push'")
	ErrPullSourceRequiresURL = errors.New("pull source requires a URL")
	ErrAssetNameRequired     = errors.New("asset name is required")
	ErrAssetPathRequired     = errors.New("asset path is required")
	ErrInvalidAssetType      = errors.New("invalid asset type")
	ErrInvalidAssetStatus    = errors.New("invalid asset status")
)
