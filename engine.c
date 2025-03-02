#include <math.h>
#include <string.h>

#include "unit.h"
#include "engine.h"


/* A four-stroke piston engine takes two revolutions per cycle. */
static const double REV_PER_CYCLE = 2.0;

/* Specific heat of dry air at constant pressure at T=300K [J/(kg*K)]. */
static const double C_P_AIR = 1005.0;

/* Specific heat of dry air at constant volume at T=300K [J/(kg*K)]. */
static const double C_V_AIR = 718.0;

/* Heat capacity ratio of dry air at T=300K [J/(kg*K)]. */
static const double GAMMA_AIR = C_P_AIR / C_V_AIR;


static VolumeFlowRate port_volume_flow_rate(const Engine *e);
static double density_ratio(const Engine *e);


void
init_engine(Engine *e) {
	memset(e, 0, sizeof(*e));
}

Pressure
comp_outlet_pressure(const Engine *e) {
	return e->map + e->intercooler_deltap;
}

/* Pressure ratio across the compressor. */
double
pressure_ratio(const Engine *e) {
	Pressure p1, p2;

	p1 = e->ambient_pressure;
	p2 = comp_outlet_pressure(e);
	return p2 / p1;
}

Temperature
comp_outlet_temperature(const Engine *e) {
	Temperature t1, dt;
	Pressure p1, p2;
	double exp;

	t1 = e->ambient_temperature;
	p1 = e->ambient_pressure;
	p2 = comp_outlet_pressure(e);
	exp = (GAMMA_AIR - 1.0) / GAMMA_AIR;
	dt = t1 * (pow(p2/p1, exp) - 1.0) / e->comp_efficiency;

	return  t1 + dt;
}

Temperature
manifold_temperature(const Engine *e) {
	Temperature t1, t2;

	t1 = e->ambient_temperature;
	t2 = comp_outlet_temperature(e);
	return t2 - (t2 - t1)*e->intercooler_efficiency;
}

/* Volume flow rate throught the compressor inlet. */
VolumeFlowRate
volume_flow_rate(const Engine *e) {
	VolumeFlowRate v3;
	double r;

	v3 = port_volume_flow_rate(e);
	r = density_ratio(e);
	return v3 * r;
}

/* Volume flow rate through the intake ports. */
static VolumeFlowRate
port_volume_flow_rate(const Engine *e) {
	double n, d, ve;

	n = as_rpm(e->rpm);
	d = as_cubic_metre(e->displacement);
	ve = e->ve;
	return cubic_metre_per_min(n * d * ve / REV_PER_CYCLE);
}

/* Density ratio between the ports and the compressor inlet. */
static double
density_ratio(const Engine *e) {
	Pressure p1, p3;
	Temperature t1, t3;

	p1 = e->ambient_pressure;
	p3 = e->map;
	t1 = e->ambient_temperature;
	t3 = manifold_temperature(e);
	return (p1 * t3) / (p3 * t1);
}
