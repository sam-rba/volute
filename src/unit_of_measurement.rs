pub trait UnitOfMeasurement {
    type Unit;

    fn from_unit(unit: Self::Unit, n: i32) -> Self;
    fn as_unit(&self, unit: Self::Unit) -> i32;
}

pub mod pressure {
    use super::UnitOfMeasurement;

    #[derive(Default, Clone)]
    pub struct Pressure(i32);

    #[derive(Copy, Clone)]
    pub enum Unit {
        Pascal = 1,
        KiloPascal = 1000,
    }

    impl Unit {
        // Pseudo iter::Cycle behavior.
        pub fn next(&mut self) {
            match self {
                Unit::Pascal => {
                    *self = Unit::KiloPascal;
                }
                Unit::KiloPascal => {
                    *self = Unit::Pascal;
                }
            }
        }
    }

    impl ToString for Unit {
        fn to_string(&self) -> String {
            match self {
                Self::Pascal => String::from("Pa"),
                Self::KiloPascal => String::from("kPa"),
            }
        }
    }

    impl UnitOfMeasurement for Pressure {
        type Unit = Unit;

        fn from_unit(unit: Self::Unit, n: i32) -> Self {
            Self(n * unit as i32)
        }

        fn as_unit(&self, unit: Self::Unit) -> i32 {
            self.0 / unit as i32
        }
    }
}
