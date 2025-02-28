enum { MAX_POINTS = 16 };

typedef struct {
	w_Field displacement;
	w_Select displacement_unit;

	w_Field rpm[MAX_POINTS];

	w_Field map[MAX_POINTS];
	w_Select map_unit;

	w_Field ve[MAX_POINTS];

	int npoints;
} UI;

void init_ui(UI *ui);
