typedef struct {
	union {
		MassFlowRate mfr;
		VolumeFlowRate vfr;
	} u;
	enum {
		MASS_FLOW,
		VOLUME_FLOW
	} t;
} Flow;

typedef struct {
	int x, y; /* pixel coordinates. */
	float pr; /* pressure ratio. */
	Flow flow;
} Point;

typedef struct {
	char brand[NAME_MAX+1]; /* e.g. Borgwarner. */
	char series[NAME_MAX+1]; /* e.g. Airwerks. */
	char model[NAME_MAX+1]; /* e.g. S200SX-E. */
	char imgfile[NAME_MAX+1]; /* name of file containing image of the compressor map. */
	Point origin, ref;
} Compressor;


int load_compressors(Compressor **comps, int *n);
