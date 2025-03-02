#include "test.h"

int
main(void) {
	test_rad_per_sec();
	test_deg_per_sec();
	test_rpm();
	test_as_rad_per_sec();
	test_as_deg_per_sec();
	test_as_rpm();

	test_percent();
	test_as_percent();

	test_pascal();
	test_millibar();
	test_kilopascal();
	test_bar();
	test_psi();
	test_inch_mercury();
	test_as_pascal();
	test_as_millibar();
	test_as_kilopascal();
	test_as_bar();
	test_as_psi();
	test_as_inch_mercury();

	test_kelvin();
	test_celsius();
	test_fahrenheit();
	test_rankine();
	test_as_kelvin();
	test_as_celsius();
	test_as_fahrenheit();
	test_as_rankine();

	test_cubic_centimetre();
	test_litre();
	test_cubic_metre();
	test_cubic_inch();
	test_as_cubic_centimetre();
	test_as_litre();
	test_as_cubic_metre();
	test_as_cubic_inch();

	test_cubic_metre_per_sec();
	test_cubic_metre_per_min();
	test_cubic_foot_per_min();
	test_as_cubic_metre_per_sec();
	test_as_cubic_metre_per_min();
	test_as_cubic_foot_per_min();

	test_kilo_per_sec();
	test_pound_per_min();
	test_as_kilo_per_sec();
	test_as_pound_per_min();

	test_comp_outlet_pressure();
	test_pressure_ratio();
	test_pressure_ratio_intercooled();
	test_comp_outlet_temperature_adiabatic();
	test_comp_outlet_temperature();
	test_manifold_temperature();
	test_volume_flow_rate();
	test_mass_flow_rate();
}
