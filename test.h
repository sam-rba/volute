#define EPSILON (1e-7)

#define test(got, want) { \
	if (got < want-EPSILON || got > want+EPSILON) { \
		fprintf(stderr, "got %lf; want %lf\n", got, want); \
		assert(got == want); \
	} \
}


void test_rad_per_sec(void);
void test_deg_per_sec(void);
void test_rpm(void);
void test_as_rad_per_sec(void);
void test_as_deg_per_sec(void);
void test_as_rpm(void);

void test_cubic_centimetre(void);
void test_litre(void);
void test_cubic_metre(void);
void test_cubic_inch(void);
void test_as_cubic_centimetre(void);
void test_as_litre(void);
void test_as_cubic_metre(void);
void test_as_cubic_inch(void);
