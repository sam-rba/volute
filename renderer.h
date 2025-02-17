#ifndef RENDERER_H
#define RENDERER_H

#include "microui.h"

void r_init(void);
void r_handle_input(mu_Context *ctx);
void r_render(mu_Context *ctx);
void r_present(void);
int r_get_text_width(const char *text, int len);
int r_get_text_height(void);
void r_get_window_size(int *w, int *h);

#endif

