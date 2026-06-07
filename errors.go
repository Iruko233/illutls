package illutls

import "errors"

var (
	// ErrProfileNotFound is returned when a requested browser profile does not exist.
	ErrProfileNotFound = errors.New("illutls: browser profile not found")
	// ErrTLSHandshake is returned when the TLS connection fails.
	ErrTLSHandshake = errors.New("illutls: TLS handshake failed")
	// ErrH2Settings is returned when HTTP/2 settings cannot be applied.
	ErrH2Settings = errors.New("illutls: HTTP/2 settings application failed")
	// ErrInvalidProxy is returned when the provided proxy URL is invalid.
	ErrInvalidProxy = errors.New("illutls: invalid proxy URL")
	// ErrNilProfile is returned if a nil profile is provided to the transport.
	ErrNilProfile = errors.New("illutls: profile must not be nil")
	// ErrConnectionFailed is returned for general dial or connection errors.
	ErrConnectionFailed = errors.New("illutls: connection failed")
)
