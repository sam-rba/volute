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

#[derive(Default)]
pub struct Mass(f64); // Base unit is grams
impl Mass {
    /* constructors */
    pub fn from_grams(grams: f64) -> Mass {
        Mass(grams)
    }

    pub fn from_kilograms(kilos: f64) -> Mass {
        Mass(kilos / 1000.)
    }

    pub fn from_moles(moles: f64, molar_mass: f64) -> Mass {
        let kilos = moles * molar_mass;
        Mass::from_kilograms(kilos)
    }

    /* metric */
    pub fn as_grams(&self) -> f64 {
        self.0
    }

    pub fn as_kilograms(&self) -> f64 {
        self.0 / 1000.
    }

    /* imperial */
    pub fn as_pounds(&self) -> f64 {
        self.0 * 0.002204623
    }
}

pub struct Pressure(f64); // Base unit is pascals
impl Pressure {
    pub fn from_pascals(pascals: f64) -> Self {
        Self(pascals)
    }

    pub fn as_pascals(&self) -> f64 {
        self.0
    }
}

#[derive(Default)]
pub struct Temperature(f64); // Base unit is kelvin
impl Temperature {
    pub fn from_kelvin(kelvin: f64) -> Temperature {
        Temperature(kelvin)
    }

    pub fn as_kelvin(&self) -> f64 {
        self.0
    }
}

pub struct Volume(f64); // Base unit is cubic metres
impl Volume {
    pub fn from_cubic_metres(cubic_metres: f64) -> Volume {
        Volume(cubic_metres)
    }

    pub fn from_cubic_centimetres(cubic_centimetres: f64) -> Volume {
        Volume(cubic_centimetres / 1_000_000.)
    }

    pub fn as_cubic_metres(&self) -> f64 {
        self.0
    }

    pub fn as_cubic_centimetres(&self) -> f64 {
        self.0 * 1_000_000.
    }

    pub fn as_cubic_feet(&self) -> f64 {
        self.0 * 35.3147
    }
}
