#include <stdio.h>
#include <string.h>

#include "microui.h"
#include "widget.h"


#define nelem(arr) (sizeof(arr)/sizeof(arr[0]))


void
w_init_field(w_Field *f) {
	f->buf[0] = '\0';
	f->value = 0.0;
}

/* field draws a Field widget and updates its value.
 * It returns MU_RES_CHANGE if the value has changed. */
int
w_field(mu_Context *ctx, w_Field *f) {
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

void
w_set_field(w_Field *f, double val) {
	f->value = val;
	snprintf(f->buf, sizeof(f->buf), "%f", val);
}

void
w_init_select(w_Select *select, int nopts, const char *const opts[]) {
	select->nopts = nopts;
	select->opts = opts;
	select->idx = 0;
	select->oldidx = 0;
	select ->active = 0;
}

int
w_select(mu_Context *ctx, w_Select *select) {
	mu_Id id;
	mu_Rect r;
	int width, res, i;

	mu_layout_begin_column(ctx);
	width = -1;
	mu_layout_row(ctx, 1, &width, 0);

	id = mu_get_id(ctx, &select, sizeof(select));
	r = mu_layout_next(ctx);
	mu_update_control(ctx, id, r, 0);

	select->active ^= (ctx->mouse_pressed == MU_MOUSE_LEFT && ctx->focus == id);

	mu_draw_control_frame(ctx, id, r, MU_COLOR_BUTTON, 0);
	const char *label = select->opts[select->idx];
	mu_draw_control_text(ctx, label, r, MU_COLOR_TEXT, 0);

	res = 0;
	if (select->active) {
		res = MU_RES_ACTIVE;
		for (i = 0; i < select->nopts; i++) {
			if (mu_button(ctx, select->opts[i])) {
				select->oldidx = select->idx;
				select->idx = i;
				res |= MU_RES_CHANGE;
				select->active = 0;
			}
		}
	}

	mu_layout_end_column(ctx);

	return res;
}

void
w_init_label(w_Label label) {
	label[0] = '\0';
}

void
w_label(mu_Context *ctx, const w_Label label) {
	mu_label(ctx, label);
}
