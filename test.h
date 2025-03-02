#define EPSILON (1e-7)

#define test(got, want) { \
	if (got < want-EPSILON || got > want+EPSILON) { \
		fprintf(stderr, "got %.7f; want %.7f\n", got, want); \
		assert(got == want); \
	} \
}


void test_rad_per_sec(void);
void test_deg_per_sec(void);
void test_rpm(void);
void test_as_rad_per_sec(void);
void test_as_deg_per_sec(void);
void test_as_rpm(void);

void test_percent(void);
void test_as_percent(void);

void test_pascal(void);
void test_millibar(void);
void test_kilopascal(void);
void test_bar(void);
void test_psi(void);
void test_inch_mercury(void);
void test_as_pascal(void);
void test_as_millibar(void);
void test_as_kilopascal(void);
void test_as_bar(void);
void test_as_psi(void);
void test_as_inch_mercury(void);

void test_kelvin(void);
void test_celsius(void);
void test_fahrenheit(void);
void test_rankine(void);
void test_as_kelvin(void);
void test_as_celsius(void);
void test_as_fahrenheit(void);
void test_as_rankine(void);

void test_cubic_centimetre(void);
void test_litre(void);
void test_cubic_metre(void);
void test_cubic_inch(void);
void test_as_cubic_centimetre(void);
void test_as_litre(void);
void test_as_cubic_metre(void);
void test_as_cubic_inch(void);

void test_cubic_metre_per_sec(void);
void test_cubic_metre_per_min(void);
void test_cubic_foot_per_min(void);
void test_as_cubic_metre_per_sec(void);
void test_as_cubic_metre_per_min(void);
void test_as_cubic_foot_per_min(void);

void test_kilo_per_sec(void);
void test_pound_per_min(void);
void test_as_kilo_per_sec(void);
void test_as_pound_per_min(void);

void test_comp_outlet_pressure(void);
void test_pressure_ratio(void);
void test_pressure_ratio_intercooled(void);
void test_comp_outlet_temperature_adiabatic(void);
void test_comp_outlet_temperature(void);
void test_manifold_temperature(void);
void test_volume_flow_rate(void);
