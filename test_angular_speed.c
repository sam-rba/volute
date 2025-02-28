#include <assert.h>
#include <stdio.h>

#include "test.h"
#include "unit.h"

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
