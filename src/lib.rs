pub mod app;
pub mod input;
pub mod unit_of_measurement;

use crate::unit_of_measurement::{Pressure, Temperature, Volume};

const GAS_CONSTANT: f64 = 8.314472;
const MOLAR_MASS_OF_AIR: f64 = 0.0289647; // Kg/mol

fn moles_from_gas_law(pres: Pressure, vol: Volume, temp: Temperature) -> f64 {
    (pres.as_pascals() * vol.as_cubic_metres()) / (GAS_CONSTANT * temp.as_kelvin())
}
