typedef struct {
	Volume displacement;
	Temperature ambient_temperature;
	Pressure ambient_pressure;
	AngularSpeed rpm;
	Pressure map;
	Fraction ve;
	Fraction comp_efficiency;
	Fraction intercooler_efficiency;
	Pressure intercooler_deltap;
} Engine;

void init_engine(Engine *e);
Pressure comp_outlet_pressure(const Engine *e);
double pressure_ratio(const Engine *e);
Temperature comp_outlet_temperature(const Engine *e);
Temperature manifold_temperature(const Engine *e);
VolumeFlowRate volume_flow_rate(const Engine *e);
