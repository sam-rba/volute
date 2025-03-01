typedef struct {
	Volume displacement;
	Temperature ambient_temperature;
	Pressure ambient_pressure;
	AngularSpeed rpm;
	Pressure map;
	Fraction ve;
	Fraction comp_efficiency;
} Engine;

void init_engine(Engine *e);
VolumeFlowRate volume_flow_rate(const Engine *e);
