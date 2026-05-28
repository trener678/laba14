//go:build !rustffi

package validator

import "laba14-health-pipeline/internal/health"

func Validate(r health.Reading) error {
	return health.ValidateReading(r)
}
