#include <errno.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <dirent.h>

#include "cwalk.h"
#include "toml.h"
#include "util.h"
#include "unit.h"
#include "compressor.h"
#include "eprintf.h"
#include "util.h"


static const char ROOT[NAME_MAX+1] = "compressor_maps";


static int load_compressor(const char *path, Compressor *comp);
static int load_point(const toml_table_t *tbl, const char *key, const char *flowunit, Point *pt);
static int parse_flow(double val, const char *unit, Flow *flow);
static int parse_mass_flow(double val, const char *unit, Flow *flow);
static int parse_volume_flow(double val, const char *unit, Flow *flow);
static int toml_filter(const struct dirent *de);
static int cmp_flow_unit(const void *key, const void *datum);


/* Load descriptions of all of the compressor maps.
 * Places a malloc-allocated array of compressors at *comps,
 * and the number of compressors at *n.
 * Returns 0 on success. */
int
load_compressors(Compressor **comps, int *n) {
	struct dirent **files;
	int nfiles;
	char path[3 * NAME_MAX];
	Compressor comp;

	*comps = NULL;
	*n = 0;

	nfiles = scandir(ROOT, &files, toml_filter, alphasort);
	if (nfiles < 0) {
		weprintf("failed to scan %s", ROOT);
		return 1;
	}

	*comps = malloc(nfiles * sizeof(Compressor));
	if (*comps == NULL) {
		weprintf("malloc failed");
		return 1;
	}

	/* TODO: parallelize. */
	while (nfiles > 0) {
		(void) cwk_path_join(ROOT, files[nfiles-1]->d_name, path, nelem(path));
		if (load_compressor(path, &comp) != 0) {
			weprintf("failed to load compressor from %s", path);
			free_arr((void **) files, nfiles);
			free(*comps);
			return 1;
		}
		(*comps)[(*n)++] = comp;
		free(files[--nfiles]);
	}
	free(files);

	return 0;
}

/* Load a compressor toml file into *comp. Returns 0 on success. */
static int
load_compressor(const char *path, Compressor *comp) {
	FILE *f;
	char errbuf[256];
	toml_table_t *tbl;
	toml_value_t brand, series, model, flowunit;
	int err;

	f = fopen(path, "r");
	if (f == NULL) {
		weprintf("failed to open %s", path);
		return 1;
	}

	tbl = toml_parse_file(f, errbuf, sizeof(errbuf));
	if (!tbl) {
		weprintf("failed to parse %s: %s", path, errbuf);
		return 1;
	}

	if (fclose(f) != 0) {
		weprintf("failed to close %s", path);
		toml_free(tbl);
		return 1;
	}

	err = 0;
	brand = toml_table_string(tbl, "brand");
	if (!brand.ok) {
		weprintf("%s: missing 'brand'", path);
		err = 1;
	}
	series = toml_table_string(tbl, "series");
	if (!series.ok) {
		weprintf("%s: missing 'series'", path);
		err = 1;
	}
	model = toml_table_string(tbl, "model");
	if (!model.ok) {
		weprintf("%s: missing 'model'", path);
		err = 1;
	}
	flowunit = toml_table_string(tbl, "flowunit");
	if (!flowunit.ok) {
		weprintf("%s: missing 'flowunit'", path);
		err = 1;
	}
	if (err) {
		toml_free(tbl);
		return 1;
	}

	strncpy(comp->brand, brand.u.s, nelem(comp->brand)-1);
	strncpy(comp->series, series.u.s, nelem(comp->series)-1);
	strncpy(comp->model, model.u.s, nelem(comp->model)-1);

	(void) cwk_path_change_extension(path, "jpg", comp->imgfile, sizeof(comp->imgfile));

	if (load_point(tbl, "origin", flowunit.u.s, &comp->origin) != 0) {
		weprintf("%s: failed to load 'origin'", path);
		toml_free(tbl);
		return 1;
	}
	if (load_point(tbl, "ref", flowunit.u.s, &comp->ref) != 0) {
		weprintf("%s: failed to load 'ref'", path);
		toml_free(tbl);
		return 1;
	}

	toml_free(tbl);

	return 0;
}

/* load a Point from a compressor toml file.
 * key - the name of the toml subtable containing the point.
 * pt - the returned point.
 * Returns 0 on success. */
static int
load_point(const toml_table_t *tbl, const char *key, const char *flowunit, Point *pt) {
	toml_table_t *subtbl;
	toml_value_t x, y, pr, flowval;
	Flow flow;

	subtbl = toml_table_table(tbl, key);
	if (!subtbl) {
		weprintf("missing table '%s'", key);
		return 1;
	}

	x = toml_table_int(subtbl, "x");
	if (!x.ok) {
		weprintf("%s: missing 'x'", key);
		return 1;
	}
	pt->x = x.u.i;

	y = toml_table_int(subtbl, "y");
	if (!y.ok) {
		weprintf("%s: missing 'y'", key);
		return 1;
	}
	pt->y = y.u.i;

	pr = toml_table_double(subtbl, "pr");
	if (!pr.ok) {
		weprintf("%s: missing 'pr'", key);
		return 1;
	}
	pt->pr = pr.u.d;

	flowval = toml_table_double(subtbl, "flow");
	if (!flowval.ok) {
		weprintf("%s: missing 'flow'", key);
		return 1;
	}
	if (parse_flow(flowval.u.d, flowunit, &flow) != 0) {
		weprintf("invalid flow: %f %s", flowval.u.d, flowunit);
		return 1;
	}
	pt->flow = flow;

	return 0;
}

static int
parse_flow(double val, const char *unit, Flow *flow) {
	if (parse_mass_flow(val, unit, flow) == 0) {
		return 0;
	}
	if (parse_volume_flow(val, unit, flow) == 0) {
		return 0;
	}
	return 1;
}

static int
parse_mass_flow(double val, const char *unit, Flow *flow) {
	int i;

	i = lsearch(unit, mass_flow_rate_units, n_mass_flow_rate_units, sizeof(*mass_flow_rate_units), cmp_flow_unit);
	if (i >= 0) {
		flow->u.mfr = mass_flow_rate_makers[i](val);
		flow->t = MASS_FLOW;
		return 0;
	}

	return 1;
}

static int
parse_volume_flow(double val, const char *unit, Flow *flow) {
	int i;

	i = lsearch(unit, volume_flow_rate_units, n_volume_flow_rate_units, sizeof(volume_flow_rate_units[0]), cmp_flow_unit);
	if (i >= 0) {
		flow->u.vfr = volume_flow_rate_makers[i](val);
		flow->t = VOLUME_FLOW;
		return 0;
	}

	return 1;
}

static int
cmp_flow_unit(const void *key, const void *datum) {
	return strcmp((char *) key, *(char **) datum);
}

static int
toml_filter(const struct dirent *de) {
	const char *extension;
	size_t length;

	if (!cwk_path_get_extension(de->d_name, &extension, &length)) {
		return 0; /* no extension. */
	}
	return strcmp(".toml", extension) == 0; /* extension is ".toml". */
}
