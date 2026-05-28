//go:build rustffi

package validator

/*
#cgo LDFLAGS: -L${SRCDIR}/../../rust/health_validator/target/release -lhealth_validator
#include <stdlib.h>
int validate_reading(const char* metric, double value, const char* patient_id, const char* sensor_id);
*/
import "C"

import (
	"fmt"
	"unsafe"

	"laba14-health-pipeline/internal/health"
)

func Validate(r health.Reading) error {
	metric := C.CString(r.Metric)
	patient := C.CString(r.PatientID)
	sensor := C.CString(r.SensorID)
	defer C.free(unsafe.Pointer(metric))
	defer C.free(unsafe.Pointer(patient))
	defer C.free(unsafe.Pointer(sensor))

	code := C.validate_reading(metric, C.double(r.Value), patient, sensor)
	if code != 0 {
		return fmt.Errorf("rust validator rejected reading with code %d", int(code))
	}
	return nil
}
