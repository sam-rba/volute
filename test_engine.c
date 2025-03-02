#include <assert.h>
#include <stdio.h>

#include "test.h"
#include "unit.h"
#include "engine.h"

void
test_comp_outlet_pressure(void) {
	Engine e;
	init_engine(&e);
	e.map = millibar(2000);
	e.intercooler_deltap = psi(0.4);
	test(comp_outlet_pressure(&e), millibar(2027.579029173));
}

void
test_pressure_ratio(void) {
	Engine e = {
		.ambient_pressure = psi(14.3),
		.map = psi(14.3+18),
	};
	test(pressure_ratio(&e), 2.2587413);

}
void
test_pressure_ratio_intercooled(void) {
	Engine e = {
		.ambient_pressure = psi(14.3),
		.map = psi(14.3+18),
		.intercooler_deltap = psi(0.4),
	};
	test(pressure_ratio(&e), 2.2867133);
}

void
test_comp_outlet_temperature_adiabatic(void) {
	Engine e = {
		.ambient_temperature = fahrenheit(70),
		.ambient_pressure = psi(14.7),
		.map = psi(31.7),
		.comp_efficiency = percent(100),
	};
	test(comp_outlet_temperature(&e), kelvin(366.4715514));
}

void
test_comp_outlet_temperature(void) {
	Engine e = {
		.ambient_temperature = fahrenheit(70),
		.ambient_pressure = psi(14.7),
		.map = psi(31.7),
		.comp_efficiency = percent(70),
	};
	test(comp_outlet_temperature(&e), kelvin(397.418883));
}

void
test_manifold_temperature(void) {
	Engine e = {
		.ambient_temperature = fahrenheit(80),
		.ambient_pressure = millibar(1015.9166),
		.map = millibar(2031.8332),
		.comp_efficiency = percent(70),
		.intercooler_efficiency = percent(70),
	};
	test(manifold_temperature(&e), kelvin(327.9429247));
}


void
test_volume_flow_rate(void) {
	Pressure p_ambient, p_boost, map;

	p_ambient = inch_mercury(30);
	p_boost = psi(10);
	map = p_ambient + p_boost;
	Engine e = {
		.displacement = cubic_inch(250),
		.rpm = rpm(5000),
		.map = map,
		.ambient_temperature = fahrenheit(70),
		.ambient_pressure = p_ambient,
		.ve = percent(80),
		.comp_efficiency = percent(65),
	};
	test(volume_flow_rate(&e), cubic_metre_per_sec(0.184086));
}

void
test_mass_flow_rate(void) {
	Pressure p_ambient, p_boost, map;

	p_ambient = inch_mercury(30);
	p_boost = psi(10);
	map = p_ambient + p_boost;
	Engine e = {
		.displacement = cubic_inch(250),
		.rpm = rpm(5000),
		.map = map,
		.ambient_temperature = fahrenheit(70),
		.ambient_pressure = p_ambient,
		.ve = percent(80),
		.comp_efficiency = percent(65),
	};
	test(mass_flow_rate(&e), kilo_per_sec(0.2214056));
}
