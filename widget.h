enum { LABEL_SIZE = 128 };


/* Field is a floating point number input field. */
typedef struct {
	char buf[64];
	double value;
} w_Field;

void w_init_field(w_Field *f);
int w_field(mu_Context *ctx, w_Field *f);


typedef struct {
	int nopts;
	const char *const *opts;
	int idx; /* index of selected option. */
	int active;
} w_Select;

void w_init_select(w_Select *select, int nopts, const char *const opts[]);
int w_select(mu_Context *ctx, w_Select *select);


typedef char w_Label[LABEL_SIZE];

void w_init_label(w_Label label);
void w_label(mu_Context *ctx, const w_Label label);
