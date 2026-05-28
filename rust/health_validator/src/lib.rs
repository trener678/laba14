use std::ffi::CStr;
use std::os::raw::{c_char, c_double};

#[derive(Debug, PartialEq)]
pub enum ValidationError {
    EmptyField,
    UnknownMetric,
    OutOfRange,
    NonFinite,
}

pub fn validate(
    metric: &str,
    value: f64,
    patient_id: &str,
    sensor_id: &str,
) -> Result<(), ValidationError> {
    if metric.is_empty() || patient_id.is_empty() || sensor_id.is_empty() {
        return Err(ValidationError::EmptyField);
    }
    if !value.is_finite() {
        return Err(ValidationError::NonFinite);
    }
    match metric {
        "heart_rate" if (25.0..=240.0).contains(&value) => Ok(()),
        "spo2" if (50.0..=100.0).contains(&value) => Ok(()),
        "temperature" if (30.0..=45.0).contains(&value) => Ok(()),
        "heart_rate" | "spo2" | "temperature" => Err(ValidationError::OutOfRange),
        _ => Err(ValidationError::UnknownMetric),
    }
}

#[unsafe(no_mangle)]
pub extern "C" fn validate_reading(
    metric: *const c_char,
    value: c_double,
    patient_id: *const c_char,
    sensor_id: *const c_char,
) -> i32 {
    let Some(metric) = read_cstr(metric) else {
        return 1;
    };
    let Some(patient_id) = read_cstr(patient_id) else {
        return 1;
    };
    let Some(sensor_id) = read_cstr(sensor_id) else {
        return 1;
    };

    match validate(metric, value, patient_id, sensor_id) {
        Ok(()) => 0,
        Err(ValidationError::EmptyField) => 1,
        Err(ValidationError::UnknownMetric) => 2,
        Err(ValidationError::OutOfRange) => 3,
        Err(ValidationError::NonFinite) => 4,
    }
}

fn read_cstr(ptr: *const c_char) -> Option<&'static str> {
    if ptr.is_null() {
        return None;
    }
    unsafe { CStr::from_ptr(ptr) }.to_str().ok()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn accepts_valid_spo2() {
        assert_eq!(validate("spo2", 97.0, "p1", "s1"), Ok(()));
    }

    #[test]
    fn rejects_out_of_range_temperature() {
        assert_eq!(
            validate("temperature", 55.0, "p1", "s1"),
            Err(ValidationError::OutOfRange)
        );
    }
}
