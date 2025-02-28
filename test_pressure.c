#include <assert.h>
#include <stdio.h>

#include "test.h"
#include "unit.h"

void
test_pascal(void) {
	test(pascal(123.456), 123.456);
}

void
test_millibar(void) {
	test(millibar(123.456), 12345.6);
}

void
test_kilopascal(void) {
	test(kilopascal(123.456), 123456.0);
}

void
test_bar(void) {
	test(bar(123.456), 12345600.0);
}

void
test_psi(void) {
	test(psi(123.456), 851199.15638539323);
}

void
test_as_pascal(void) {
	test(as_pascal(pascal(123.456)), 123.456);
}

void
test_as_millibar(void) {
	test(as_millibar(millibar(123.456)), 123.456);
}

void
test_as_kilopascal(void) {
	test(as_kilopascal(kilopascal(123.456)), 123.456);
}

void
test_as_bar(void) {
	test(as_bar(bar(123.456)), 123.456);
}

void
test_as_psi(void) {
	test(as_psi(psi(123.456)), 123.456);
}
