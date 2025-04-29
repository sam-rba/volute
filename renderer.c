#include <limits.h>
#include <stdio.h>
#include <stdlib.h>

#include <SDL2/SDL.h>
#include <SDL2/SDL_ttf.h>
#include <SDL2/SDL_image.h>
#include "microui.h"
#include "renderer.h"


#define expect(x) do {                                               \
    if (!(x)) {                                                      \
      fprintf(stderr, "Fatal error: %s:%d: assertion '%s' failed\n", \
        __FILE__, __LINE__, #x);                                     \
      abort();                                                       \
    }                                                                \
  } while (0)


enum window {
	WIDTH = 640,
	HEIGHT = 480,
	WINFLAGS = SDL_WINDOW_RESIZABLE,
	RENDERFLAGS = SDL_RENDERER_PRESENTVSYNC,
};

enum { ICONLIST_SIZE = 32 };

static const char FONT[] = "font/P052-Roman.ttf";
enum font { FONTSIZE = 14, };

static const mu_Color bg = {255, 255, 255, 255};

static const char button_map[256] = {
	[SDL_BUTTON_LEFT & 0xff] = MU_MOUSE_LEFT,
	[SDL_BUTTON_RIGHT & 0xff] = MU_MOUSE_RIGHT,
	[SDL_BUTTON_MIDDLE & 0xff] = MU_MOUSE_MIDDLE,
};

static const char key_map[256] = {
	[SDLK_LSHIFT & 0xff] = MU_KEY_SHIFT,
	[SDLK_RSHIFT & 0xff] = MU_KEY_SHIFT,
	[SDLK_LCTRL & 0xff] = MU_KEY_CTRL,
	[SDLK_RCTRL & 0xff] = MU_KEY_CTRL,
	[SDLK_LALT & 0xff] = MU_KEY_ALT,
	[SDLK_RALT & 0xff] = MU_KEY_ALT,
	[SDLK_RETURN & 0xff] = MU_KEY_RETURN,
	[SDLK_BACKSPACE & 0xff] = MU_KEY_BACKSPACE,
};

mu_stack(SDL_Texture *, ICONLIST_SIZE) icon_list;


static void print_info(void);
static int text_width(mu_Font mufont, const char *str, int len);
static int text_height(mu_Font mufont);
static void handle_event(SDL_Event e, mu_Context *ctx);
static void clear(void);
static void render_command(mu_Command *cmd);
static void clip(mu_Rect rect);
static void draw_rect(mu_Rect rect, mu_Color color);
static void draw_text(mu_Font font, mu_Vec2 pos, mu_Color color, const char *str);
static void draw_icon(int id, mu_Rect r);

static SDL_Window *window = NULL;
static SDL_Renderer *renderer = NULL;

/* Initialize the window and renderer. Returns non-zero on error. */
int
r_init(mu_Context *ctx, const char *title) {
	icon_list.idx = 0;

	if (SDL_Init(SDL_INIT_VIDEO) != 0) {
		fprintf(stderr, "%s\n", SDL_GetError());
		return 1;
	}
	window = SDL_CreateWindow(title,
		SDL_WINDOWPOS_UNDEFINED, SDL_WINDOWPOS_UNDEFINED,
		WIDTH, HEIGHT,
		WINFLAGS);
	if (!window) {
		fprintf(stderr, "%s\n", SDL_GetError());
		return 1;
	}
	renderer = SDL_CreateRenderer(window, -1, RENDERFLAGS);
	if (!renderer) {
		fprintf(stderr, "%s\n", SDL_GetError());
		return 1;
	}

	if (TTF_Init() != 0) {
		fprintf(stderr, "%s\n", SDL_GetError());
		return 1;
	}
	TTF_Font *font = TTF_OpenFont(FONT, FONTSIZE);
	if (!font) {
		fprintf(stderr, "Failed to open font %s\n", FONT);
		return 1;
	}
	ctx->style->font = font;

	print_info();

	ctx->text_width = text_width;
	ctx->text_height = text_height;

	return 0;
}

void
r_free(void) {
	while (icon_list.idx-- > 0) {
		SDL_DestroyTexture(icon_list.items[icon_list.idx]);
	}
}

static void
print_info(void) {
	SDL_RendererInfo info;
	if (SDL_GetRendererInfo(renderer, &info) != 0) {
		fprintf(stderr, "%s\n", SDL_GetError());
		return;
	}
	printf("Using renderer %s\n", info.name);
	fflush(stdout);
}

static int
text_width(mu_Font font, const char *str, int len) {
	if (!str || !*str) { return 0; }

	int w = 0;
	int c = 0;
	if (TTF_MeasureUTF8(font, str, INT_MAX, &w, &c) != 0) {
		w = 0;
	}
	return w;
}

static int
text_height(mu_Font font) {
	return TTF_FontHeight(font);
}

void
r_input(mu_Context *ctx) {
	SDL_Event e;
	while (SDL_PollEvent(&e)) {
		handle_event(e, ctx);
	}
}

static void
handle_event(SDL_Event e, mu_Context *ctx) {
	switch (e.type) {
		case SDL_QUIT: {
			TTF_CloseFont(ctx->style->font);
			TTF_Quit();
			SDL_Quit();
			exit(EXIT_SUCCESS);
		}
		break; case SDL_MOUSEMOTION: {
			mu_input_mousemove(ctx, e.motion.x, e.motion.y);
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
		break; case SDL_MOUSEWHEEL: {
			mu_input_scroll(ctx, 0, e.wheel.y * -30);
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
		break; case SDL_TEXTINPUT: {
			mu_input_text(ctx, e.text.text);
		}
	}
}

void
r_render(mu_Context *ctx) {
	clear();

	mu_Command *cmd = NULL;
	while (mu_next_command(ctx, &cmd)) {
		render_command(cmd);
	}

	SDL_RenderPresent(renderer);
}

static void
clear(void) {
	if (SDL_SetRenderDrawColor(renderer, bg.r, bg.g, bg.b, bg.a) != 0) {
		fprintf(stderr, "%s", SDL_GetError());
	}
	if (SDL_RenderClear(renderer) != 0) {
		fprintf(stderr, "%s", SDL_GetError());
	}
}

static void
render_command(mu_Command *cmd) {
	switch (cmd->type) {
		case MU_COMMAND_CLIP: {
			clip(cmd->clip.rect);
		}
		break; case MU_COMMAND_RECT: {
			draw_rect(cmd->rect.rect, cmd->rect.color);
		}
		break; case MU_COMMAND_TEXT: {
			draw_text(cmd->text.font, cmd->text.pos, cmd->text.color, cmd->text.str);
		}
		break; case MU_COMMAND_ICON: {
			draw_icon(cmd->icon.id, cmd->icon.rect);
		}
	}
}

static void
clip(mu_Rect rect) {
	SDL_Rect r = {rect.x, rect.y, rect.w, rect.h};
	if (SDL_RenderSetClipRect(renderer, &r) != 0) {
		fprintf(stderr, "%s\n", SDL_GetError());
	}
}

static void
draw_rect(mu_Rect rect, mu_Color color) {
	if (SDL_SetRenderDrawColor(renderer, color.r, color.g, color.b, color.a) != 0) {
		fprintf(stderr, "%s\n", SDL_GetError());
	}
	SDL_Rect r = {rect.x, rect.y, rect.w, rect.h};
	if (SDL_RenderFillRect(renderer, &r) != 0) {
		fprintf(stderr, "%s\n", SDL_GetError());
	}
}

static void
draw_text(mu_Font font, mu_Vec2 pos, mu_Color color, const char *str) {
	if (!str || !*str) { return; }

	SDL_Color sdl_color = {color.r, color.g, color.b, color.a};
	SDL_Surface *surface = TTF_RenderUTF8_Blended(font, str, sdl_color);
	if (!surface) {
		fprintf(stderr, "%s\n", TTF_GetError());
		return;
	}
	SDL_Texture *texture = SDL_CreateTextureFromSurface(renderer, surface);
	if (!texture) {
		fprintf(stderr, "%s\n", SDL_GetError());
		SDL_FreeSurface(surface);
		return;
	}

	int descent = TTF_FontDescent(font);
	if (descent > 0) { /* v-align fonts with positive descent. */
		pos.y -= descent;
	}
	SDL_Rect dst = {pos.x, pos.y, surface->w, surface->h};
	if (SDL_RenderCopy(renderer, texture, NULL, &dst) != 0) {
		fprintf(stderr, "%s\n", SDL_GetError());
	}

	SDL_DestroyTexture(texture);
	SDL_FreeSurface(surface);
}

static void
draw_icon(int id, mu_Rect r) {
	SDL_Rect rect;

	expect(id >= 0 && id < icon_list.idx);

	rect = (SDL_Rect) {r.x, r.y, r.w, r.h};
	if (SDL_RenderCopy(renderer, icon_list.items[id], NULL, &rect) != 0) {
		fprintf(stderr, "%s\n", SDL_GetError());
	}
}

void
r_get_window_size(int *w, int*h) {
	if (SDL_GetRendererOutputSize(renderer, w, h) != 0) {
		fprintf(stderr, "%s\n", SDL_GetError());
	}
}

/* Load an image and add it to the list of icons. Returns the id of the icon, or -1 on error. */
int
r_add_icon(const char *path) {
	SDL_Surface *surface;

	expect(icon_list.idx < ICONLIST_SIZE);

	surface = IMG_Load(path);
	if (!surface) {
		fprintf(stderr, "failed to load %s: %s\n", path, SDL_GetError());
		return -1;
	}
	icon_list.items[icon_list.idx] = SDL_CreateTextureFromSurface(renderer, surface);
	if (!icon_list.items[icon_list.idx]) {
		fprintf(stderr, "%s\n", SDL_GetError());
		SDL_FreeSurface(surface);
		return -1;
	}
	SDL_FreeSurface(surface);
	return icon_list.idx++;
}

/* Remove the icon with the specified id from the icons list. */
void
r_remove_icon(int id) {
	SDL_Texture **dst, **src;
	size_t size;

	expect(id >= 0 && id < icon_list.idx);

	SDL_DestroyTexture(icon_list.items[id]);

	dst = icon_list.items + id;
	src = icon_list.items + id + 1;
	size = (icon_list.idx - id - 1) * sizeof(*icon_list.items);
	memmove(dst, src, size);

	icon_list.idx--;
}

void
r_get_icon_size(int id, int *w, int *h) {
	expect(id >= 0 && id < icon_list.idx);

	*w = *h = 0;
	if (SDL_QueryTexture(icon_list.items[id], NULL, NULL, w, h) != 0) {
		fprintf(stderr, "%s\n", SDL_GetError());
	}
}
