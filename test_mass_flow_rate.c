#include <assert.h>
#include <stdio.h>

#include "test.h"
#include "unit.h"

void
test_kilo_per_sec(void) {
	test(kilo_per_sec(123.456), 123.456);
}

void
test_pound_per_min(void) {
	test(pound_per_min(123.456), 0.9333117);
}

void
test_as_kilo_per_sec(void) {
	test(as_kilo_per_sec(kilo_per_sec(123.456)), 123.456);
}

void
test_as_pound_per_min(void) {
	test(as_pound_per_min(pound_per_min(123.456)), 123.456);
}
