#include <assert.h>
#include <stdio.h>

#include "unit.h"


#define EPSILON (1e-7)

#define test(got, want) { \
	if (got < want-EPSILON || got > want+EPSILON) { \
		fprintf(stderr, "got %lf; want %lf\n", got, want); \
		assert(got == want); \
	} \
}


void
test_rad_per_sec(void) {
	test(rad_per_sec(123.456), 123.456);
}

void
test_deg_per_sec(void) {
	test(deg_per_sec(123.456), 2.15471367888);
}

void
test_rpm(void) {
	test(rpm(123.456), 12.92828207328);
}

void
test_as_rad_per_sec(void) {
	test(as_rad_per_sec(rad_per_sec(123.456)), 123.456);
}

void
test_as_deg_per_sec(void) {
	test(as_deg_per_sec(deg_per_sec(123.456)), 123.456);
}

void
test_as_rpm(void) {
	test(as_rpm(rpm(123.456)), 123.456);
}

int
main(void) {
	test_rad_per_sec();
	test_deg_per_sec();
	test_rpm();
	test_as_rad_per_sec();
	test_as_deg_per_sec();
	test_as_rpm();
}
