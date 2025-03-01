#include <math.h>

#include "unit.h"


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
