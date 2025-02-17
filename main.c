/* Headers. */

#include <stdio.h>

#include <SDL2/SDL.h>
#include "renderer.h"
#include "microui.h"
#include "widget.h"
#include "ui.h"


/* Constants. */

static const char TITLE[] = "volute";

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

static const char button_map[256] = {
	[ SDL_BUTTON_LEFT & 0xff ] = MU_MOUSE_LEFT,
	[ SDL_BUTTON_RIGHT & 0xff ] = MU_MOUSE_RIGHT,
	[ SDL_BUTTON_MIDDLE & 0xff ] = MU_MOUSE_MIDDLE,
};

static const char key_map[256] = {
	[ SDLK_LSHIFT & 0xff ] = MU_KEY_SHIFT,
	[ SDLK_RSHIFT & 0xff ] = MU_KEY_SHIFT,
	[ SDLK_LCTRL & 0xff ] = MU_KEY_CTRL,
	[ SDLK_RCTRL & 0xff ] = MU_KEY_CTRL,
	[ SDLK_LALT & 0xff ] = MU_KEY_ALT,
	[ SDLK_RALT & 0xff ] = MU_KEY_ALT,
	[ SDLK_RETURN & 0xff ] = MU_KEY_RETURN,
	[ SDLK_BACKSPACE & 0xff ] = MU_KEY_BACKSPACE,
};


/* Function declarations. */

static void set_style(mu_Context *ctx);
static int text_width(mu_Font font, const char *text, int len);
static int text_height(mu_Font font);
static void main_loop(mu_Context *ctx, UI *ui);
static void handle_events(mu_Context *ctx);
static void handle_event(SDL_Event e, mu_Context *ctx);
static void process_frame(mu_Context *ctx, UI *ui);
static void main_window(mu_Context *ctx, UI *ui);
static void render(mu_Context *ctx);
static void render_command(mu_Command *cmd);


/* Function Definitions. */

int
main(void) {
	/* Init SDL and renderer. */
	SDL_Init(SDL_INIT_EVERYTHING);
	r_init();

	/* Init microui. */
	static mu_Context ctx;
	mu_init(&ctx);
	ctx.text_width = text_width;
	ctx.text_height = text_height;
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

static int
text_width(mu_Font font, const char *text, int len) {
	if (len < 0) {
		len = strlen(text);
	}
	return r_get_text_width(text, len);
}

static int
text_height(mu_Font font) {
	return r_get_text_height();
}

static void
main_loop(mu_Context *ctx, UI *ui) {
	for (;;) {
		handle_events(ctx);
		process_frame(ctx, ui);
		render(ctx);
	}
}

static void
handle_events(mu_Context *ctx) {
	SDL_Event e;
	while (SDL_PollEvent(&e)) {
		handle_event(e, ctx);
	}
}

static void
handle_event(SDL_Event e, mu_Context *ctx) {
	switch (e.type) {
		case SDL_QUIT: {
			exit(EXIT_SUCCESS);
		}
		break; case SDL_MOUSEMOTION: {
			mu_input_mousemove(ctx, e.motion.x, e.motion.y);
		}
		break; case SDL_MOUSEWHEEL: {
			mu_input_scroll(ctx, 0, e.wheel.y * -30);
		}
		break; case SDL_TEXTINPUT: {
			mu_input_text(ctx, e.text.text);
		}
		break; case SDL_MOUSEBUTTONDOWN: case SDL_MOUSEBUTTONUP: {
			int b = button_map[e.button.button & 0xff];
			if (b && e.type == SDL_MOUSEBUTTONDOWN) {
				mu_input_mousedown(ctx, e.button.x, e.button.y, b);
			}
			if (b && e.type == SDL_MOUSEBUTTONUP) {
				mu_input_mouseup(ctx, e.button.x, e.button.y, b);
			}
		}
		break; case SDL_KEYDOWN: case SDL_KEYUP: {
			int c = key_map[e.key.keysym.sym & 0xff];
			if (c && e.type == SDL_KEYDOWN) {
				mu_input_keydown(ctx, c);
			}
			if (c && e.type == SDL_KEYUP) {
				mu_input_keyup(ctx, c);
			}
		}
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

	if (!mu_begin_window(ctx, TITLE, mu_rect(0, 0, w, h))) {
		exit(EXIT_SUCCESS);
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

static void
render(mu_Context *ctx) {
	r_clear(COLOR_WINDOWBG);
	mu_Command *cmd = NULL;
	while (mu_next_command(ctx, &cmd)) {
		render_command(cmd);
	}
	r_present();
}

static void
render_command(mu_Command *cmd) {
	switch (cmd->type) {
		case MU_COMMAND_TEXT: {
			r_draw_text(cmd->text.str, cmd->text.pos, cmd->text.color);
		}
		break; case MU_COMMAND_RECT: {
			r_draw_rect(cmd->rect.rect, cmd->rect.color);
		}
		break; case MU_COMMAND_ICON: {
			r_draw_icon(cmd->icon.id, cmd->icon.rect, cmd->icon.color);
		}
		break; case MU_COMMAND_CLIP: {
			r_set_clip_rect(cmd->clip.rect);
		}
	}
}
