#include <stdlib.h>

#include "util.h"

void
free_arr(void **arr, int n) {
	while (n-- > 0) {
		free(arr[n]);
	}
	free(arr);
}
