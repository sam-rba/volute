#include <assert.h>
#include <stdio.h>

#include "test.h"
#include "unit.h"

void
test_percent(void) {
	test(percent(12.345), 0.12345);
}

void
test_as_percent(void) {
	test(as_percent(percent(12.345)), 12.345);
}
