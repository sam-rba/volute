pub mod flow_rate;
pub mod mass;
pub mod pressure;
pub mod temperature;
pub mod volume;

use crate::{pressure::Pressure, temperature::Temperature, volume::Volume};

const GAS_CONSTANT: f64 = 8.314472;
const MOLAR_MASS_OF_AIR: f64 = 0.0289647; // Kg/mol

fn moles_from_gas_law(pressure: Pressure, volume: Volume, temperature: Temperature) -> f64 {
    (pressure.as_pascals() * volume.as_cubic_metres()) / (GAS_CONSTANT * temperature.as_kelvin())
}
