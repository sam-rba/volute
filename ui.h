enum { MAX_POINTS = 16 };

typedef struct {
	w_Field displacement;
	w_Select displacement_unit;

	w_Field ambient_temperature;
	w_Select ambient_temperature_unit;

	w_Field ambient_pressure;
	w_Select ambient_pressure_unit;

	int npoints;

	w_Field rpm[MAX_POINTS];

	w_Field map[MAX_POINTS];
	w_Select map_unit;

	w_Field ve[MAX_POINTS];

	w_Field comp_efficiency[MAX_POINTS];

	Engine points[MAX_POINTS];

	w_Select volume_flow_rate_unit;
	w_Number volume_flow_rate[MAX_POINTS];
} UI;

void init_ui(UI *ui);
void set_displacement(UI *ui);
void set_displacement_unit(UI* ui);
void set_ambient_temperature(UI *ui);
void set_ambient_temperature_unit(UI *ui);
void set_ambient_pressure(UI *ui);
void set_ambient_pressure_unit(UI *ui);
void set_map(UI *ui, int idx);
void set_map_unit(UI *ui);
void set_ve(UI *ui, int idx);
void set_comp_efficiency(UI *ui, int idx);
void set_volume_flow_rate(UI *ui, int idx);
void set_all_volume_flow_rate(UI *ui);
void insert_point(UI *ui, int idx);
void remove_point(UI *ui, int idx);
