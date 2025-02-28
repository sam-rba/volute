#include <assert.h>
#include <stdio.h>

#include "test.h"
#include "unit.h"

void
test_cubic_metre_per_sec(void) {
	test(cubic_metre_per_sec(123.456), 123.456);
}

void
test_cubic_metre_per_min(void) {
	test(cubic_metre_per_min(123.456), 2.0576);
}

void
test_cubic_foot_per_min(void) {
	test(cubic_foot_per_min(123.456), 0.0582647436);
}

void
test_as_cubic_metre_per_sec(void) {
	test(as_cubic_metre_per_sec(cubic_metre_per_sec(123.456)), 123.456);
}

void
test_as_cubic_metre_per_min(void) {
	test(as_cubic_metre_per_min(cubic_metre_per_min(123.456)), 123.456);
}

void
test_as_cubic_foot_per_min(void) {
	test(as_cubic_foot_per_min(cubic_foot_per_min(123.456)), 123.456);
}
