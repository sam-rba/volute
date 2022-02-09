pub mod pressure {
    pub enum PressureUnit {
        Pascal = 1, // base unit. Every other variant will be a multiple of this.
        KiloPascal = 1000,
    }

    #[derive(Default)]
    pub struct Pressure {
        val: i32, // Base unit is pascals.
    }

    impl Pressure {
        pub fn from_unit(unit: PressureUnit, n: i32) -> Self {
            Self {
                val: n * unit as i32,
            }
        }

        pub fn as_unit(&self, unit: PressureUnit) -> i32 {
            self.val / unit as i32
        }
    }
}
