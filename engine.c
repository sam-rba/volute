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


void
init_engine(Engine *e) {
	memset(e, 0, sizeof(*e));
}

/* Pressure ratio across the compressor. */
double
pressure_ratio(const Engine *e) {
	Pressure p1, p2;

	p1 = e->ambient_pressure;
	p2 = comp_outlet_pressure(e);
	return p2 / p1;
}

Pressure
comp_outlet_pressure(const Engine *e) {
	return e->map + e->intercooler_deltap;
}

Temperature
comp_outlet_temperature(const Engine *e) {
	Temperature t1;
	Pressure p1, p2;

	t1 = e->ambient_temperature;
	p1 = e->ambient_pressure;
	p2 = comp_outlet_pressure(e);
	return t1 * pow(p2/p1, (GAMMA_AIR-1.0)/GAMMA_AIR);
}

VolumeFlowRate
volume_flow_rate(const Engine *e) {
	double n = as_rpm(e->rpm);
	double d = as_cubic_metre(e->displacement);
	double ve = e->ve;
	return cubic_metre_per_min(n * d * ve / REV_PER_CYCLE);
}
