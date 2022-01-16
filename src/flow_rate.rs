use crate::{mass::Mass, volume::Volume};
use std::time::Duration;

pub struct MassFlowRate {
    pub mass: Mass,
    pub duration: Duration,
}

impl MassFlowRate {
    pub fn as_kilograms_per_minute(&self) -> f64 {
        self.mass.as_kilograms() / (self.duration.as_secs() as f64 / 60.)
    }

    pub fn as_pounds_per_minute(&self) -> f64 {
        self.mass.as_pounds() / (self.duration.as_secs() as f64 / 60.)
    }
}

pub struct VolumetricFlowRate {
    pub volume: Volume,
    pub duration: Duration,
}

impl VolumetricFlowRate {
    pub fn as_cubic_metres_per_second(&self) -> f64 {
        self.volume.as_cubic_metres() / self.duration.as_secs() as f64
    }

    pub fn as_cubic_feet_per_minute(&self) -> f64 {
        self.volume.as_cubic_feet() / (self.duration.as_secs() as f64 / 60.)
    }
}
