#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>

#include "eprintf.h"

void
eprintf(const char *fmt, ...) {
	va_list args;

	va_start(args, fmt);
	fprintf(stderr, "error: ");
	vfprintf(stderr, fmt, args);
	fprintf(stderr, "\n");
	va_end(args);
	exit(1);
}

void weprintf(const char *fmt, ...) {
	va_list args;

	va_start(args, fmt);
	fprintf(stderr, "warning: ");
	vfprintf(stderr, fmt, args);
	fprintf(stderr, "\n");
	va_end(args);
}
