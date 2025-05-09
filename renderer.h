int r_init(mu_Context *ctx, const char *title);
void r_free(void);
void r_input(mu_Context *ctx);
void r_render(mu_Context *ctx);
void r_get_window_size(int *w, int *h);

int r_add_icon(const char *path);
void r_remove_icon(int);
void r_get_icon_size(int id, int *w, int *h);

int r_add_canvas(const char *bg_img_path);
void r_remove_canvas(int id);
void r_canvas_draw_circle(int id, int x, int y, int r, mu_Color color);
void r_clear_canvas(int id);
int r_render_canvas(int id);
