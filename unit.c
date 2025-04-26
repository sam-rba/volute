#include <math.h>
#include <stddef.h>

#include "unit.h"
#include "util.h"


/* Kilograms per pound. */
#define KG_PER_LB 0.45359237

/* Acceleration of gravity [m/(s*s)]. */
#define G 9.80665

/* Metres per inch. */
#define M_PER_IN 0.0254

/* Metres per foot. */
#define M_PER_FT 0.3048

/* Seconds per minute. */
#define SEC_PER_MIN 60.0

/* Zero Celsius in Kelvin. */
#define ZERO_C 273.15

const char *const pressure_units[] = {"mbar", "kPa", "bar", "psi", "inHG"};
const PressureMaker pressure_makers[nelem(pressure_units)] = {
	millibar, kilopascal, bar, psi, inch_mercury,
};
const PressureReader pressure_readers[nelem(pressure_units)] = {
	as_millibar, as_kilopascal, as_bar, as_psi, as_inch_mercury,
};
const size_t n_pressure_units = nelem(pressure_units);

const char *const temperature_units[] = {"°C", "K", "°F", "°R"};
const TemperatureMaker temperature_makers[nelem(temperature_units)] = {
	celsius, kelvin, fahrenheit, rankine,
};
const TemperatureReader temperature_readers[nelem(temperature_units)] = {
	as_celsius, as_kelvin, as_fahrenheit, as_rankine,
};
const size_t n_temperature_units = nelem(temperature_units);

const char *const volume_units[] = {"cc", "l", "ci"};
const VolumeMaker volume_makers[nelem(volume_units)] = {
	cubic_centimetre, litre, cubic_inch,
};
const VolumeReader volume_readers[nelem(volume_units)] = {
	as_cubic_centimetre, as_litre, as_cubic_inch,
};
const size_t n_volume_units = nelem(volume_units);

const char *const volume_flow_rate_units[] = {"m³/s", "CFM"};
const VolumeFlowRateMaker volume_flow_rate_makers[nelem(volume_flow_rate_units)] = {
	cubic_metre_per_sec, cubic_foot_per_min,
};
const VolumeFlowRateReader volume_flow_rate_readers[nelem(volume_flow_rate_units)] = {
	as_cubic_metre_per_sec, as_cubic_foot_per_min,
};
const size_t n_volume_flow_rate_units = nelem(volume_flow_rate_units);

const char *const mass_flow_rate_units[] = {"kg/s", "lb/min"};
const MassFlowRateMaker mass_flow_rate_makers[nelem(mass_flow_rate_units)] = {
	kilo_per_sec, pound_per_min,
};
const MassFlowRateReader mass_flow_rate_readers[nelem(mass_flow_rate_units)] = {
	as_kilo_per_sec, as_pound_per_min,
};
const size_t n_mass_flow_rate_units = nelem(mass_flow_rate_units);

AngularSpeed
rad_per_sec(double x) {
	return x;
}

AngularSpeed
deg_per_sec(double x) {
	return x * M_PI / 180.0;
}

AngularSpeed rpm(double x) {
	return x * 2.0 * M_PI / 60.0;
}

double as_rad_per_sec(AngularSpeed x) {
	return x;
}

double as_deg_per_sec(AngularSpeed x) {
	return x * 180.0 / M_PI;
}

double as_rpm(AngularSpeed x) {
	return x * 60.0 / (2.0 * M_PI);
}


Fraction
percent(double x) {
	return x / 100.0;
}

double
as_percent(Fraction x) {
	return x * 100.0;
}


Pressure
pascal(double x) {
	return x;
}

Pressure
millibar(double x) {
	return x * 100.0;
}

Pressure
kilopascal(double x) {
	return x * 1e3;
}

Pressure
bar(double x) {
	return x * 1e5;
}

Pressure
psi(double x) {
	return x * KG_PER_LB * G / pow(M_PER_IN, 2);
}

Pressure
inch_mercury(double x) {
	return x * 3386.389;
}

double
as_pascal(Pressure x) {
	return x;
}

double
as_millibar(Pressure x) {
	return x / 100.0;
}

double
as_kilopascal(Pressure x) {
	return x * 1e-3;
}

double
as_bar(Pressure x) {
	return x * 1e-5;
}

double
as_psi(Pressure x) {
	return x * pow(M_PER_IN, 2) / (KG_PER_LB * G);
}

double
as_inch_mercury(Pressure x) {
	return x / 3386.389;
}


Temperature
kelvin(double x) {
	return x;
}

Temperature
celsius(double x) {
	return x + ZERO_C;
}

Temperature
fahrenheit(double x) {
	double c;

	c = (x - 32.0) * 5.0 / 9.0;
	return c + ZERO_C;
}

Temperature
rankine(double x) {
	return x * 5.0 / 9.0;
}

double
as_kelvin(Temperature t) {
	return t;
}

double
as_celsius(Temperature t) {
	return t - ZERO_C;
}

double
as_fahrenheit(Temperature t) {
	return as_celsius(t) * 9.0 / 5.0 + 32.0;
}

double
as_rankine(Temperature t) {
	return t * 9.0 / 5.0;
}


Volume
cubic_centimetre(double x) {
	return x * 1e-6;
}

Volume
litre(double x) {
	return x * 1e-3;
}

Volume
cubic_metre(double x) {
	return x;
}

Volume
cubic_inch(double x) {
	return x * 1.6387064e-5;
}

double
as_cubic_centimetre(Volume x) {
	return x * 1e6;
}

double
as_litre(Volume x) {
	return x * 1e3;
}

double
as_cubic_metre(double x) {
	return x;
}

double
as_cubic_inch(double x) {
	return x / 1.6387064e-5;
}


VolumeFlowRate
cubic_metre_per_sec(double x) {
	return x;
}

VolumeFlowRate
cubic_metre_per_min(double x) {
	return x / SEC_PER_MIN;
}

VolumeFlowRate
cubic_foot_per_min(double x) {
	return x * pow(M_PER_FT, 3) / SEC_PER_MIN;
}

double
as_cubic_metre_per_sec(VolumeFlowRate x) {
	return x;
}

double
as_cubic_metre_per_min(VolumeFlowRate x) {
	return x * SEC_PER_MIN;
}

double
as_cubic_foot_per_min(VolumeFlowRate x) {
	return x * SEC_PER_MIN / pow(M_PER_FT, 3);
}


MassFlowRate
kilo_per_sec(double x) {
	return x;
}

MassFlowRate
pound_per_min(double x) {
	return x * KG_PER_LB / SEC_PER_MIN;
}

double
as_kilo_per_sec(MassFlowRate m) {
	return m;
}

double
as_pound_per_min(MassFlowRate m) {
	return m / KG_PER_LB * SEC_PER_MIN;
}
