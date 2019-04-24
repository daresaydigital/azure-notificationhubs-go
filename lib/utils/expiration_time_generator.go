package utils

import "time"

type (
	// ExpirationTimeGenerator generates an expiration time
	ExpirationTimeGenerator interface {
		GenerateTimestamp() int64
	}

	// ExpirationTimeGeneratorFunc is a function generating expiration times
	ExpirationTimeGeneratorFunc func() int64
)

// NewExpirationTimeGenerator creates the default generator
func NewExpirationTimeGenerator() ExpirationTimeGenerator {
	return ExpirationTimeGeneratorFunc(generateExpirationTimestamp)
}

// GenerateTimestamp calls f()
func (f ExpirationTimeGeneratorFunc) GenerateTimestamp() int64 {
	return f()
}

// generateExpirationTimestamp generates token expiration timestamp value
func generateExpirationTimestamp() int64 {
	return time.Now().Unix() + int64(3600)
}
