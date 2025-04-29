int r_init(mu_Context *ctx, const char *title);
void r_free(void);
void r_input(mu_Context *ctx);
void r_render(mu_Context *ctx);
void r_get_window_size(int *w, int *h);
int r_add_icon(const char *path);
void r_remove_icon(int);
