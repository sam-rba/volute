/* Field is a floating point number input field. */
typedef struct {
	char buf[64];
	double value;
} Field;

void init_field(Field *f);
int field(mu_Context *ctx, Field *f);
