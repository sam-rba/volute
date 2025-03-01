#include <assert.h>
#include <stdio.h>

#include "test.h"
#include "unit.h"

void
test_kelvin(void) {
	test(kelvin(123.456), 123.456);
}

void
test_celsius(void) {
	test(celsius(123.456), 396.606);
}

void
test_fahrenheit(void) {
	test(fahrenheit(123.456), 323.9588889);
}

void
test_rankine(void) {
	test(rankine(123.456), 68.5866667);
}

void
test_as_kelvin(void) {
	test(as_kelvin(kelvin(123.456)), 123.456);
}

void
test_as_celsius(void) {
	test(as_celsius(celsius(123.456)), 123.456);
}

void
test_as_fahrenheit(void) {
	test(as_fahrenheit(fahrenheit(123.456)), 123.456);
}

void
test_as_rankine(void) {
	test(as_rankine(rankine(123.456)), 123.456);
}
