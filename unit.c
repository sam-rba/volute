#include <math.h>

#include "unit.h"


/* Kilograms per pound. */
#define KG_PER_LB 0.45359237

/* Acceleration of gravity [m/(s*s)]. */
#define G 9.80665

/* Metres per inch. */
#define M_PER_IN 0.0254


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

Fraction percent(double x);
double as_percent(Fraction x);


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
