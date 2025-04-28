#include <stdlib.h>

#include "util.h"

void
free_arr(void **arr, int n) {
	while (n-- > 0) {
		free(arr[n]);
	}
	free(arr);
}

/* lsearch linearly searches base[0]...base[n-1] for an item that matches *key.
 * The function cmp must return zero if its first argument (the search key)
 * equals its second (a table entry), non-zero if not equal.
 * Returns the index of the first occurrence of key in base, or -1 if not present. */
int
lsearch(const void *key, const void *base, size_t n, size_t size, int (*cmp)(const void *keyval, const void *datum)) {
	size_t i;

	for (i = 0; i < n; i++) {
		if (cmp(key, base) == 0) {
			return i;
		}
		base = (char *) base + size;
	}
	return -1;
}
