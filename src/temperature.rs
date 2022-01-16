pub struct Temperature(f64); // Base unit is kelvin

impl Temperature {
    pub fn from_kelvin(kelvin: f64) -> Temperature {
        Temperature(kelvin)
    }

    pub fn as_kelvin(&self) -> f64 {
        self.0
    }
}
