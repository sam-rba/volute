#include <errno.h>
#include <limits.h>
#include <stdlib.h>
#include <strings.h>

#include <dirent.h>

#include "unit.h"


static const char ROOT[] = "compressor_maps/";


typedef struct {
	int x, y;
} Point;

typedef enum { MASS_FLOW, VOLUME_FLOW } FlowType;

typedef union {
	MassFlowRate mfr;
	VolumeFlowRate vfr;
} Flow;

typedef struct {
	char brand[NAME_MAX+1]; /* e.g. Borgwarner. */
	char series[NAME_MAX+1]; /* e.g. Airwerks. */
	char model[NAME_MAX+1]; /* e.g. S200SX-E. */
	char imgfile[NAME_MAX+1]; /* name of file containing image of the compressor map. */
	Point origin, refpt; /* pixel coords of origin and reference point. */
	double originpr, refpr; /* pressure ratio at origin and reference point. */
	Flow originflow, refflow; /* flow at origin and reference point. */
	FlowType flowtype; /* mass-flow or volume-flow (x-axis). */
} Compressor;
