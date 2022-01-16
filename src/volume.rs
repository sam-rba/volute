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
