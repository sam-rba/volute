#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "microui.h"
#include "unit.h"
#include "compressor.h"
#include "widget.h"
#include "util.h"


#define FORMAT "%.5g"

static const mu_Color RED = {255, 0, 0, 255};

static const char *sc_selected_name(w_Select_Compressor *select);
static int select_compressor_active(mu_Context *ctx, w_Select_Compressor *select);
static void sc_filter(w_Select_Compressor *select);
static void update_active(mu_Context *ctx, mu_Id id, mu_Rect r, int *active);


void
w_init_field(w_Field *f) {
	f->buf[0] = '\0';
	f->value = 0.0;
	f->invalid = 0;
}

/* field draws a Field widget and updates its value.
 * It returns MU_RES_CHANGE if the value has changed. */
int
w_field(mu_Context *ctx, w_Field *f) {
	mu_Rect rect;
	int changed;
	char s[2];
	double value;

	rect = mu_layout_next(ctx);
	mu_layout_set_next(ctx, rect, 0);

	changed = 0;
	if (mu_textbox(ctx, f->buf, sizeof(f->buf)) & MU_RES_CHANGE) {
		/* s used to catch erroneous chars at end of field. */
		if (sscanf(f->buf, "%lf %1s", &value, s) == 1) {
			f->value = value;
			f->invalid = 0;
			changed = 1;
		} else if (f->buf[0] == '\0') {
			f->value = 0.0;
			f->invalid = 0;
			changed = 1;
		} else {
			f->invalid = 1;
		}
	}

	if (f->invalid) {
		mu_draw_box(ctx, rect, RED);
	}

	return changed ? MU_RES_CHANGE : 0;
}

void
w_set_field(w_Field *f, double val) {
	f->value = val;
	snprintf(f->buf, sizeof(f->buf), FORMAT, val);
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
	update_active(ctx, id, r, &select->active);

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

/* Returns non-zero on error. */
int
w_init_select_compressor(w_Select_Compressor *select, int n, const Compressor *comps) {
	int i;
	size_t namesize;

	select->comps = comps;
	select->n = n;

	namesize = sizeof((*comps).brand) + sizeof((*comps).series) + sizeof((*comps).model) + 3;
	select->names = malloc(n * sizeof(*select->names));
	if (select->names == NULL) {
		free(select->filtered);
		return 1;
	}
	/* TODO: parallelize. */
	for (i = 0; i < n; i++) {
		select->names[i] = malloc(namesize * sizeof(char));
		if (select->names[i] == NULL) {
			free_arr((void **) select->names, i);
			free(select->filtered);
			return 1;
		}
		snprintf(select->names[i], namesize, "%s %s %s",
			comps[i].brand, comps[i].series, comps[i].model);
	}

	memset(select->brand_filter, 0, sizeof(select->brand_filter));
	memset(select->series_filter, 0, sizeof(select->series_filter));
	memset(select->model_filter, 0, sizeof(select->model_filter));

	select->filtered = malloc(n * sizeof(*select->filtered));
	if (select->filtered == NULL) {
		return 1;
	}
	/* TODO: parallelize. */
	for (i = 0; i < n; i++) {
		select->filtered[i] = i;
	}
	select->nfiltered = n;

	select->idx = 0;
	select -> oldidx = 0;

	select->active = 0;

	return 0;
}

void
w_free_select_compressor(w_Select_Compressor *select) {
	free(select->filtered);
	free_arr((void **) select->names, select->n);
}

int
w_select_compressor(mu_Context *ctx, w_Select_Compressor *select) {
	int width;
	mu_Id id;
	mu_Rect r;

	width = 3*LABEL_WIDTH + 2*ctx->style->spacing;
	mu_layout_row(ctx, 1, &width, 0);
	id = mu_get_id(ctx, &select, sizeof(select));
	r = mu_layout_next(ctx);
	update_active(ctx, id, r, &select->active);

	mu_draw_control_frame(ctx, id, r, MU_COLOR_BUTTON, 0);
	mu_draw_control_text(ctx, sc_selected_name(select), r, MU_COLOR_TEXT, 0);

	if (select->active) {
		return select_compressor_active(ctx, select);
	}

	return 0;
}

static const char *
sc_selected_name(w_Select_Compressor *select) {
	static const char *none = "";

	if (select->idx < 0 || select->idx >= select->n) {
		return none;
	}
	return select->names[select->idx];
}

static int
select_compressor_active(mu_Context *ctx, w_Select_Compressor *select) {
	int filter_changed, res, i, j, width;

	const int widths[] = {LABEL_WIDTH, LABEL_WIDTH, LABEL_WIDTH};

	mu_layout_row(ctx, nelem(widths), widths, 0);
	mu_label(ctx, "brand");
	mu_label(ctx, "series");
	mu_label(ctx, "model");

	mu_layout_row(ctx, nelem(widths), widths, 0);
	filter_changed = 0;
	filter_changed |= mu_textbox(ctx, select->brand_filter, sizeof(select->brand_filter));
	filter_changed |= mu_textbox(ctx, select->series_filter, sizeof(select->series_filter));
	filter_changed |= mu_textbox(ctx, select->model_filter, sizeof(select->model_filter));
	filter_changed &= MU_RES_SUBMIT;
	if (filter_changed) {
		sc_filter(select);
	}

	res = 0;
	for (i = 0; i < select->nfiltered; i++) {
		width = 3*LABEL_WIDTH + 2*ctx->style->spacing;
		mu_layout_row(ctx, 1, &width, 0);
		j = select->filtered[i];
		if (mu_button(ctx, select->names[j])) {
			select->oldidx = select->idx;
			select->idx = j;
			res = MU_RES_CHANGE;
			select->active = 0;
		}
	}

	return res;
}

static void
sc_filter(w_Select_Compressor *select) {
	const char *brand, *series, *model;
	int i;
	const Compressor *comp;

	brand = select->brand_filter;
	series = select->series_filter;
	model = select->model_filter;

	/* TODO: parallelize */
	select->nfiltered = 0;
	for (i = 0; i < select->n; i++) {
		comp = &select->comps[i];

		if (strspn(comp->brand, brand) != strlen(brand)) {
			continue;
		} else if (strspn(comp->series, series) != strlen(series)) {
			continue;
		} else if (strspn(comp->model, model) != strlen(model)) {
			continue;
		}

		select->filtered[select->nfiltered++] = i;
	}
}

void
w_init_number(w_Number num) {
	num[0] = '\0';
}

void
w_set_number(w_Number num, double val) {
	snprintf(num, sizeof(w_Number), FORMAT, val);
}

void
w_number(mu_Context *ctx, const w_Number num) {
	mu_label(ctx, num);
}

/* Update the active/selected status of a widget. id is the microui ID of the widget.
 * *active is the active flag of the widget that will be toggled if anywhere in r is clicked. */
static void
update_active(mu_Context *ctx, mu_Id id, mu_Rect r, int *active) {
	mu_update_control(ctx, id, r, 0);
	*active ^= (ctx->mouse_pressed == MU_MOUSE_LEFT && ctx->focus == id);
}
