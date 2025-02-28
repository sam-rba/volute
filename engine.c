#include <string.h>

#include "unit.h"
#include "engine.h"


/* A four-stroke piston engine takes two revolutions per cycle. */
#define REV_PER_CYCLE 2.0


void
init_engine(Engine *e) {
	memset(e, 0, sizeof(*e));
}

VolumeFlowRate
volume_flow_rate(const Engine *e) {
	double n = as_rpm(e->rpm);
	double d = as_cubic_metre(e->displacement);
	double ve = e->ve;
	return cubic_metre_per_min(n * d * ve / REV_PER_CYCLE);
}
