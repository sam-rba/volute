use std::time::Duration;
use volute::{flow_rate::MassFlowRate, mass::Mass, pressure::Pressure};

fn main() {
    let mass = Mass::from_grams(1600.);
    println!("{} Kg", mass.as_kilograms());

    let mass_flow_rate = MassFlowRate {
        mass: mass,
        duration: Duration::from_secs(5),
    };
    println!("{:.2} Kg/min.", mass_flow_rate.as_kilograms_per_minute());

    let pressure = Pressure::from_pascals(1500.52);
    println!("{} Pa", pressure.as_pascals());
}
