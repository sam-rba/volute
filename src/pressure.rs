pub struct Pressure(f64); // Base unit is pascals

impl Pressure {
    pub fn from_pascals(pascals: f64) -> Self {
        Self(pascals)
    }

    pub fn as_pascals(&self) -> f64 {
        self.0
    }
}
