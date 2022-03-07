pub trait UnitOfMeasurement {
    type Unit;

    fn from_unit(unit: Self::Unit, n: i32) -> Self;
    fn as_unit(&self, unit: Self::Unit) -> i32;
}

pub mod pressure {
    use super::UnitOfMeasurement;

    #[derive(Default)]
    pub struct Pressure {
        val: i32,
    }

    pub enum Unit {
        Pascal = 1,
        KiloPascal = 1000,
    }

    impl UnitOfMeasurement for Pressure {
        type Unit = Unit;

        fn from_unit(unit: Self::Unit, n: i32) -> Self {
            Self {
                val: n * unit as i32,
            }
        }

        fn as_unit(&self, unit: Self::Unit) -> i32 {
            self.val / unit as i32
        }
    }
}
