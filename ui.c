#include <string.h>

#include "microui.h"
#include "widget.h"
#include "unit.h"
#include "engine.h"
#include "ui.h"


#define nelem(arr) (sizeof(arr)/sizeof(arr[0]))


static const char *const displacement_units[] = {"cc", "l", "ci"};
static const char *const map_units[] = {"mbar", "kPa", "bar", "psi"};


void
init_ui(UI *ui) {
	w_init_field(&ui->displacement);
	w_init_select(&ui->displacement_unit, nelem(displacement_units), displacement_units);
	
ui->npoints = 1;

	w_init_field(&ui->rpm[0]);

	w_init_field(&ui->map[0]);
	w_init_select(&ui->map_unit, nelem(map_units), map_units);

	w_init_field(&ui->ve[0]);

	init_engine(&ui->points[0]);
}

void
insert_point(UI *ui, int idx) {
	int i;

	if (idx < 0 || idx >= ui->npoints || ui->npoints >= MAX_POINTS) {
		return;
	}

	for (i = ui->npoints; i > idx; i--) {
		memmove(&ui->rpm[i], &ui->rpm[i-1], sizeof(ui->rpm[i-1]));
		memmove(&ui->map[i], &ui->map[i-1], sizeof(ui->map[i-1]));
		memmove(&ui->ve[i], &ui->ve[i-1], sizeof(ui->ve[i-1]));
		memmove(&ui->points[i], &ui->points[i-1], sizeof(ui->points[i-1]));
	}
	ui->npoints++;
}

void
remove_point(UI *ui, int idx) {
	if (idx < 0 || idx >= ui->npoints || ui->npoints <= 1) {
		return;
	}

	for (; idx < ui->npoints-1; idx++) {
		memmove(&ui->rpm[idx], &ui->rpm[idx+1], sizeof(ui->rpm[idx]));
		memmove(&ui->map[idx], &ui->map[idx+1], sizeof(ui->map[idx]));
		memmove(&ui->ve[idx], &ui->ve[idx+1], sizeof(ui->ve[idx]));
		memmove(&ui->points[idx], &ui->points[idx+1], sizeof(ui->points[idx]));
	}
	ui->npoints--;
}
