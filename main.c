/* Headers. */

#include <limits.h>
#include <stdio.h>
#include <stdlib.h>

#include "microui.h"
#include "renderer.h"
#include "unit.h"
#include "compressor.h"
#include "widget.h"
#include "engine.h"
#include "ui.h"


/* Macros. */

#define nelem(arr) (sizeof(arr)/sizeof(arr[0]))


/* Constants. */

static const char TITLE[] = "volute";
enum window {
	WIN_OPTS = MU_OPT_NOINTERACT | MU_OPT_NOTITLE | MU_OPT_AUTOSIZE | MU_OPT_NOFRAME,
};


static const mu_Color BLACK = {0, 0, 0, 255};
static const mu_Color WHITE = {255, 255, 255, 255};
static const mu_Color LIGHT_GRAY = {222, 222, 222, 255};
static const mu_Color DARK_GRAY = {128, 128, 128, 255};

static const mu_Color COLOR_TEXT = BLACK;
static const mu_Color COLOR_BORDER = BLACK;
static const mu_Color COLOR_WINDOWBG = WHITE;
static const mu_Color COLOR_TITLEBG = LIGHT_GRAY;
static const mu_Color COLOR_TITLETEXT = COLOR_TEXT;
static const mu_Color COLOR_PANELBG = COLOR_WINDOWBG;
static const mu_Color COLOR_BUTTON = WHITE;
static const mu_Color COLOR_BUTTONHOVER = LIGHT_GRAY;
static const mu_Color COLOR_BUTTONFOCUS = DARK_GRAY;
static const mu_Color COLOR_BASE = WHITE;
static const mu_Color COLOR_BASEHOVER = COLOR_BASE;
static const mu_Color COLOR_BASEFOCUS = COLOR_BASE;
static const mu_Color COLOR_SCROLLBASE = WHITE;
static const mu_Color COLOR_SCROLLTHUMB = WHITE;


/* Function declarations. */

static void set_style(mu_Context *ctx);
static void main_loop(mu_Context *ctx, UI *ui);
static void process_frame(mu_Context *ctx, UI *ui);
static void main_window(mu_Context *ctx, UI *ui);
static void displacement_row(mu_Context *ctx, UI *ui);
static void ambient_temperature_row(mu_Context *ctx, UI *ui);
static void ambient_pressure_row(mu_Context *ctx, UI *ui);
static void global_input_row(mu_Context *ctx, UI *ui, const char *label, w_Field *input, void (*callback)(UI *ui), w_Select *unit, void (*unit_callback)(UI *ui));
static void rpm_row(mu_Context *ctx, UI *ui);
static void map_row(mu_Context *ctx, UI *ui);
static void ve_row(mu_Context *ctx, UI *ui);
static void comp_efficiency_row(mu_Context *ctx, UI *ui);
static void intercooler_efficiency_row(mu_Context *ctx, UI *ui);
static void intercooler_deltap_row(mu_Context *ctx, UI *ui);
static void input_row(mu_Context *ctx, UI *ui, const char *label, w_Select *unit, void (*unit_callback)(UI *ui), w_Field inputs[], void (*callback)(UI *ui, int idx));
static void input_row_static_unit(mu_Context *ctx, UI *ui, const char *label, const char *unit, w_Field inputs[], void (*callback)(UI *ui, int idx));
static void dup_del_row(mu_Context *ctx, UI *ui);
static void pressure_ratio_row(mu_Context *ctx, UI *ui);
static void comp_outlet_temperature_row(mu_Context *ctx, UI *ui);
static void manifold_temperature_row(mu_Context *ctx, UI *ui);
static void volume_flow_rate_row(mu_Context *ctx, UI *ui);
static void mass_flow_rate_row(mu_Context *ctx, UI *ui);
static void mass_flow_rate_corrected_row(mu_Context *ctx, UI *ui);
static void comp_select(mu_Context *ctx, UI *ui);
static void comp_img(mu_Context *ctx, UI *ui);
static void scale_preserve_aspect_ratio(mu_Rect *r, int wmax, int hmax);
static void output_row(mu_Context *ctx, UI *ui, const char *label, w_Select *unit, w_Number outputs[]);
static void hpad(mu_Context *ctx, int w);
static void vpad(mu_Context *ctx, int h);


/* Function Definitions. */

int
main(void) {
	/* Init microui. */
	static mu_Context ctx;
	mu_init(&ctx);
	r_init(&ctx, TITLE);
	set_style(&ctx);

	/* Init data structures. */
	static UI ui;
	if (init_ui(&ui) != 0) {
		return 1;
	}

	main_loop(&ctx, &ui);

	free_ui(&ui);
	r_free();

	return 0;
}

static void
set_style(mu_Context *ctx) {
	ctx->style->colors[MU_COLOR_TEXT] = COLOR_TEXT;
	ctx->style->colors[MU_COLOR_BORDER] = COLOR_BORDER;
	ctx->style->colors[MU_COLOR_WINDOWBG] = COLOR_WINDOWBG;
	ctx->style->colors[MU_COLOR_TITLEBG] = COLOR_TITLEBG;
	ctx->style->colors[MU_COLOR_TITLETEXT] = COLOR_TITLETEXT;
	ctx->style->colors[MU_COLOR_PANELBG] = COLOR_PANELBG;
	ctx->style->colors[MU_COLOR_BUTTON] = COLOR_BUTTON;
	ctx->style->colors[MU_COLOR_BUTTONHOVER] = COLOR_BUTTONHOVER;
	ctx->style->colors[MU_COLOR_BUTTONFOCUS] = COLOR_BUTTONFOCUS;
	ctx->style->colors[MU_COLOR_BASE] = COLOR_BASE;
	ctx->style->colors[MU_COLOR_BASEHOVER] = COLOR_BASEHOVER;
	ctx->style->colors[MU_COLOR_BASEFOCUS] = COLOR_BASEFOCUS;
	ctx->style->colors[MU_COLOR_SCROLLBASE] = COLOR_SCROLLBASE;
	ctx->style->colors[MU_COLOR_SCROLLTHUMB] = COLOR_SCROLLTHUMB;
}

static void
main_loop(mu_Context *ctx, UI *ui) {
	for (;;) {
		r_input(ctx);
		process_frame(ctx, ui);
		r_render(ctx);
	}
}

static void
process_frame(mu_Context *ctx, UI *ui) {
	mu_begin(ctx);
	main_window(ctx, ui);
	mu_end(ctx);
}

static void
main_window(mu_Context *ctx, UI *ui) {
	int w, h;

	r_get_window_size(&w, &h);
	if (!mu_begin_window_ex(ctx, TITLE, mu_rect(0, 0, w, h), WIN_OPTS)) {
		exit(EXIT_FAILURE);
	}

	displacement_row(ctx, ui);
	ambient_temperature_row(ctx, ui);
	ambient_pressure_row(ctx, ui);

	vpad(ctx, 0);

	rpm_row(ctx, ui);
	map_row(ctx, ui);
	ve_row(ctx, ui);
	comp_efficiency_row(ctx, ui);
	intercooler_efficiency_row(ctx, ui);
	intercooler_deltap_row(ctx, ui);
	dup_del_row(ctx, ui);

	vpad(ctx, 0);

	pressure_ratio_row(ctx, ui);
	comp_outlet_temperature_row(ctx, ui);
	manifold_temperature_row(ctx, ui);
	volume_flow_rate_row(ctx, ui);
	mass_flow_rate_row(ctx, ui);
	mass_flow_rate_corrected_row(ctx, ui);

	vpad(ctx, 0);

	comp_select(ctx, ui);

	comp_img(ctx, ui);

	mu_end_window(ctx);
}

static void
displacement_row(mu_Context *ctx, UI *ui) {
	global_input_row(ctx, ui,
		"Displacement:",
		&ui->displacement, set_displacement,
		&ui->displacement_unit, set_displacement_unit);
}

static void
ambient_temperature_row(mu_Context *ctx, UI *ui) {
	global_input_row(ctx, ui,
		"Ambient T:",
		&ui->ambient_temperature, set_ambient_temperature,
		&ui->ambient_temperature_unit, set_ambient_temperature_unit);
}

static void
ambient_pressure_row(mu_Context *ctx, UI *ui) {
	global_input_row(ctx, ui,
		"Ambient P:",
		&ui->ambient_pressure, set_ambient_pressure,
		&ui->ambient_pressure_unit, set_ambient_pressure_unit);
}

static void
global_input_row(mu_Context *ctx, UI *ui,
	const char *label,
	w_Field *input, void (*callback)(UI *ui),
	w_Select *unit, void (*unit_callback)(UI *ui)
) {
	mu_layout_row(ctx, 3, (int[]) {LABEL_WIDTH, FIELD_WIDTH, UNIT_WIDTH}, 0);
	mu_label(ctx, label);
	if (w_field(ctx, input) & MU_RES_CHANGE) {
		callback(ui);
		compute_all(ui);
	}
	if (w_select(ctx, unit) & MU_RES_CHANGE) {
		unit_callback(ui);
	}
}

static void
rpm_row(mu_Context *ctx, UI *ui) {
	input_row_static_unit(ctx, ui,
		"Speed:",
		"(rpm)",
		ui->rpm, set_rpm);
}

static void
map_row(mu_Context *ctx, UI *ui) {
	input_row(ctx, ui,
		"Manifold P:",
		&ui->map_unit, set_map_unit,
		ui->map, set_map);
}

static void
ve_row(mu_Context *ctx, UI *ui) {
	input_row_static_unit(ctx, ui,
		"Volumetric η:",
		"(%)",
		ui->ve, set_ve);
}

static void
comp_efficiency_row(mu_Context *ctx, UI *ui) {
	input_row_static_unit(ctx, ui,
		"Compressor η:",
		"(%)",
		ui->comp_efficiency, set_comp_efficiency);
}

static void
intercooler_efficiency_row(mu_Context *ctx, UI *ui) {
	input_row_static_unit(ctx, ui,
		"Intercooler η:",
		"(%)",
		ui->intercooler_efficiency, set_intercooler_efficiency);
}

static void
intercooler_deltap_row(mu_Context *ctx, UI *ui) {
	input_row(ctx, ui,
		"Intercooler ΔP:",
		&ui->intercooler_deltap_unit, set_intercooler_deltap_unit,
		ui->intercooler_deltap, set_intercooler_deltap);
}

static void
input_row(mu_Context *ctx, UI *ui,
	const char *label,
	w_Select *unit, void (*unit_callback)(UI *ui),
	w_Field inputs[], void (*callback)(UI *ui, int idx)
) {
	int i;

	mu_layout_row(ctx, 0, NULL, 0);
	mu_layout_width(ctx, LABEL_WIDTH);
	mu_label(ctx, label);
	mu_layout_width(ctx, UNIT_WIDTH);
	if (w_select(ctx, unit) & MU_RES_CHANGE) {
		unit_callback(ui);
	}
	mu_layout_width(ctx, FIELD_WIDTH);
	for (i = 0; i < ui->npoints; i++) {
		if (w_field(ctx, &inputs[i])) {
			callback(ui, i);
			compute(ui, i);
		}
	}
}

static void
input_row_static_unit(mu_Context *ctx, UI *ui,
	const char *label,
	const char *unit,
	w_Field inputs[], void (*callback)(UI *ui, int idx)
) {
	int i;

	mu_layout_row(ctx, 0, NULL, 0);
	mu_layout_width(ctx, LABEL_WIDTH);
	mu_label(ctx, label);
	mu_layout_width(ctx, UNIT_WIDTH);
	mu_label(ctx, unit);
	mu_layout_width(ctx, FIELD_WIDTH);
	for (i = 0; i < ui->npoints; i++) {
		if (w_field(ctx, &inputs[i])) {
			callback(ui, i);
			compute(ui, i);
		}
	}
}

static void
dup_del_row(mu_Context *ctx, UI *ui) {
	int i;

	mu_layout_row(ctx, 0, NULL, 0);
	hpad(ctx, LABEL_WIDTH);
	hpad(ctx, UNIT_WIDTH);
	mu_layout_width(ctx, (FIELD_WIDTH - ctx->style->spacing)/2);
	for (i = 0; i < ui->npoints; i++) {
		mu_push_id(ctx, &i, sizeof(i));
		if (mu_button(ctx, "Dup")) {
			insert_point(ui, i);
		}
		if (mu_button(ctx, "Del")) {
			remove_point(ui, i);
		}
		mu_pop_id(ctx);
	}
}

static void
pressure_ratio_row(mu_Context *ctx, UI *ui) {
	int i;

	mu_layout_row(ctx, 0, NULL, 0);
	mu_layout_width(ctx, LABEL_WIDTH);
	mu_label(ctx, "Pressure ratio:");
	hpad(ctx, UNIT_WIDTH);
	mu_layout_width(ctx, FIELD_WIDTH);
	for (i = 0; i < ui->npoints; i++) {
		w_number(ctx, ui->pressure_ratio[i]);
	}
}

static void
comp_outlet_temperature_row(mu_Context *ctx, UI *ui) {
	output_row(ctx, ui,
		"Compressor T:",
		&ui->comp_outlet_temperature_unit,
		ui->comp_outlet_temperature);
}

static void
manifold_temperature_row(mu_Context *ctx, UI *ui) {
	output_row(ctx, ui,
		"Manifold T:",
		&ui->manifold_temperature_unit,
		ui->manifold_temperature);
}

static void
volume_flow_rate_row(mu_Context *ctx, UI *ui) {
	output_row(ctx, ui,
		"Volume flow:",
		&ui->volume_flow_rate_unit,
		ui->volume_flow_rate);
}

static void
mass_flow_rate_row(mu_Context *ctx, UI *ui) {
	output_row(ctx, ui,
		"Mass flow:",
		&ui->mass_flow_rate_unit,
		ui->mass_flow_rate);
}

static void
mass_flow_rate_corrected_row(mu_Context *ctx, UI *ui) {
	output_row(ctx, ui,
		"Mass flow at STP:",
		&ui->mass_flow_rate_corrected_unit,
		ui->mass_flow_rate_corrected);
}

static void
comp_select(mu_Context *ctx, UI *ui) {
	if (w_select_compressor(ctx, &ui->comp_select) & MU_RES_CHANGE) {
		/* TODO */
	}
}

static void
comp_img(mu_Context *ctx, UI *ui) {
	mu_Rect win, area, img;

	win.x = win.y = 0;
	r_get_window_size(&win.w, &win.h);
	mu_layout_row(ctx, 1, &win.w, win.h);

	area = mu_layout_next(ctx);
	area.w = win.w - area.x - ctx->style->padding;
	area.h = win.h - area.y - ctx->style->padding;

	img.x = area.x;
	img.y = area.y;
	r_get_icon_size(ui->comp_img.id, &img.w, &img.h);

	scale_preserve_aspect_ratio(&img, area.w, area.h);

	mu_layout_set_next(ctx, img, 0);

	w_canvas(ctx, &ui->comp_img);
}

static void
scale_preserve_aspect_ratio(mu_Rect *r, int wmax, int hmax) {
	double wscale, hscale, scale;

	wscale = (double) wmax / (double) r->w;
	hscale = (double) hmax / (double) r->h;
	scale = mu_min(wscale, hscale);
	r->w *= scale;
	r->h *= scale;
}

static void
output_row(mu_Context *ctx, UI *ui, const char *label, w_Select *unit, w_Number outputs[]) {
	int i;

	mu_layout_row(ctx, 0, NULL, 0);
	mu_layout_width(ctx, LABEL_WIDTH);
	mu_label(ctx, label);
	mu_layout_width(ctx, UNIT_WIDTH);
	if (w_select(ctx, unit) & MU_RES_CHANGE) {
		compute_all(ui);
	}
	mu_layout_width(ctx, FIELD_WIDTH);
	for (i = 0; i < ui->npoints; i++) {
		w_number(ctx, outputs[i]);
	}
}

static void
hpad(mu_Context *ctx, int w) {
	mu_layout_width(ctx, w);
	mu_label(ctx, "");
}

static void
vpad(mu_Context *ctx, int h) {
	mu_layout_row(ctx, 0, NULL, h);
	mu_label(ctx, "");
}
