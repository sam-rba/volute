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
