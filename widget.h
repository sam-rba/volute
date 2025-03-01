enum {
	FIELD_SIZE = 64,
	NUMBER_SIZE = 128,
};


/* Field is a floating point number input field. */
typedef struct {
	char buf[FIELD_SIZE];
	double value;
	int invalid;
} w_Field;

void w_init_field(w_Field *f);
int w_field(mu_Context *ctx, w_Field *f);
void w_set_field(w_Field *f, double v);


typedef struct {
	int nopts;
	const char *const *opts;
	int idx; /* index of selected option. */
	int oldidx; /* index of previously selected option. */
	int active;
} w_Select;

void w_init_select(w_Select *select, int nopts, const char *const opts[]);
int w_select(mu_Context *ctx, w_Select *select);


typedef char w_Number[NUMBER_SIZE];

void w_init_number(w_Number num);
void w_set_number(w_Number num, double val);
void w_number(mu_Context *ctx, const w_Number num);
