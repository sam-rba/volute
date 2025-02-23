/* Headers. */

#include <stdio.h>
#include <stdlib.h>

#include "microui.h"
#include "renderer.h"
#include "widget.h"
#include "ui.h"


/* Constants. */

static const char TITLE[] = "volute";
enum window {
	WIN_OPTS = MU_OPT_NOINTERACT | MU_OPT_NOTITLE | MU_OPT_AUTOSIZE | MU_OPT_NOFRAME,
};

static const mu_Color BLACK = {0, 0, 0, 255};
static const mu_Color WHITE = {255, 255, 255, 255};
static const mu_Color LIGHT_GRAY = {222, 222, 222, 255};

static const mu_Color COLOR_TEXT = BLACK;
static const mu_Color COLOR_BORDER = BLACK;
static const mu_Color COLOR_WINDOWBG = WHITE;
static const mu_Color COLOR_TITLEBG = LIGHT_GRAY;
static const mu_Color COLOR_TITLETEXT = COLOR_TEXT;
static const mu_Color COLOR_PANELBG = COLOR_WINDOWBG;
static const mu_Color COLOR_BUTTON = LIGHT_GRAY;
static const mu_Color COLOR_BUTTONHOVER = COLOR_BUTTON;
static const mu_Color COLOR_BUTTONFOCUS = COLOR_BUTTON;
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


/* Function Definitions. */

int
main(void) {
	/* Init microui. */
	static mu_Context ctx;
	mu_init(&ctx);
	r_init(&ctx);
	set_style(&ctx);

	/* Init data structures. */
	static UI ui;
	init_ui(&ui);

	main_loop(&ctx, &ui);

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
	/* TODO */
	mu_layout_row(ctx, 2, (int[]) {0, 0}, 0);
	static double value = 0.0;
	if (field(ctx, &ui->displacement) & MU_RES_CHANGE) {
		/* TODO */
		value = ui->displacement.value;
	}
	static char buf[64];
	snprintf(buf, sizeof(buf), "%lf", value);
	mu_label(ctx, buf);
	mu_end_window(ctx);
}
