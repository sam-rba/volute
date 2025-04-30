enum layout {
	LABEL_WIDTH = 128,
	UNIT_WIDTH = 52,
	FIELD_WIDTH = 64,
};

enum {
	FIELD_SIZE = 64,
	NUMBER_SIZE = 128,
	NAME_SIZE = 256,
};


/* Field is a floating point number input field. */
typedef struct {
	char buf[FIELD_SIZE];
	double value;
	int invalid;
} w_Field;

void w_init_field(w_Field *f);
int w_field(mu_Context *ctx, w_Field *f);
void w_set_field(w_Field *f, double v);


typedef struct {
	int nopts;
	const char *const *opts;
	int idx; /* index of selected option. */
	int oldidx; /* index of previously selected option. */
	int active;
} w_Select;

void w_init_select(w_Select *select, int nopts, const char *const opts[]);
int w_select(mu_Context *ctx, w_Select *select);


typedef struct {
	const Compressor *comps;
	int n; /* len(comps) */

	char **names; /* buffer to hold names of compressors. */

	char brand_filter[NAME_SIZE];
	char series_filter[NAME_SIZE];
	char model_filter[NAME_SIZE];

	int *filtered; /* indices of compressors accepted by the filter. */
	int nfiltered; /* len(filtered) */

	int idx; /* index of selected in comps. */
	int oldidx; /* index of previously selected. */

	int active;
} w_Select_Compressor;

int w_init_select_compressor(w_Select_Compressor *select, int n, const Compressor *const comps);
void w_free_select_compressor(w_Select_Compressor *select);
int w_select_compressor(mu_Context *ctx, w_Select_Compressor *select);


typedef char w_Number[NUMBER_SIZE];

void w_init_number(w_Number num);
void w_set_number(w_Number num, double val);
void w_number(mu_Context *ctx, const w_Number num);


typedef struct {
	int id; /* icon id. */
} w_Image;

void w_init_image(w_Image *img);
void w_free_image(w_Image *img);
int w_set_image(w_Image *img, const char *path);
void w_image(mu_Context *ctx, w_Image *img);


typedef struct {
	int id; /* renderer canvas id. */
	int dirty; /* the canvas must be re-rendered. */
	int icon_id; /* renderer icon id. Only valid if !dirty. */
} w_Canvas;

int w_init_canvas(w_Canvas *c, const char *bg_img_path);
void w_free_canvas(w_Canvas *c);
void w_canvas(mu_Context *ctx, w_Canvas *canvas);
