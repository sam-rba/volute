typedef struct {
	Volume displacement;
	AngularSpeed rpm;
	Pressure map;
	Fraction ve;
} Engine;

void init_engine(Engine *e);
