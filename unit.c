#include <math.h>

#include "unit.h"


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


Pressure pascal(double x);
Pressure millibar(double x);
Pressure kilopascal(double x);
Pressure bar(double x);
Pressure psi(double x);
double as_pascal(Pressure x);
double as_millibar(Pressure x);
double as_kilopascal(Pressure x);
double as_bar(Pressure x);
double as_psi(Pressure x);


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
