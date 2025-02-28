#include <assert.h>
#include <stdio.h>

#include "test.h"
#include "unit.h"

void
test_cubic_centimetre(void) {
	test(cubic_centimetre(123.456), 0.000123456);
}

void
test_litre(void) {
	test(litre(123.456), 0.123456);
}

void
test_cubic_metre(void) {
	test(cubic_metre(123.456), 123.456);
}

void
test_cubic_inch(void) {
	test(cubic_inch(123.456), 0.0020230814);
}

void
test_as_cubic_centimetre(void) {
	test(as_cubic_centimetre(cubic_centimetre(123.456)), 123.456);
}

void
test_as_litre(void) {
	test(as_litre(litre(123.456)), 123.456);
}

void
test_as_cubic_metre(void) {
	test(as_cubic_metre(cubic_metre(123.456)), 123.456);
}

void
test_as_cubic_inch(void) {
	test(as_cubic_inch(cubic_inch(123.456)), 123.456);
}
