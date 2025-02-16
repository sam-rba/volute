#include <stdio.h>

#include "microui.h"
#include "widget.h"

void
init_field(Field *f) {
	f->buf[0] = '\0';
	f->value = 0.0;
}

/* field draws a Field widget and updates its value.
 * It returns MU_RES_CHANGE if the value has changed. */
int
field(mu_Context *ctx, Field *f) {
	double value;
	int changed = 0;
	if (mu_textbox(ctx, f->buf, sizeof(f->buf)) & MU_RES_CHANGE) {
		if (sscanf(f->buf, "%lf", &value) == 1) {
			f->value = value;
			changed = 1;
		} else if (f->buf[0] == '\0') {
			f->value = 0.0;
			changed = 1;
		}
	}
	return changed ? MU_RES_CHANGE : 0;
}
