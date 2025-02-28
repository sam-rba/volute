#include "microui.h"
#include "widget.h"
#include "ui.h"


#define nelem(arr) (sizeof(arr)/sizeof(arr[0]))


static const char *const displacement_units[] = {"cc", "l", "ci"};
static const char *const map_units[] = {"mbar", "kPa", "bar", "psi"};


void
init_ui(UI *ui) {
	w_init_field(&ui->displacement);
	w_init_select(&ui->displacement_unit, nelem(displacement_units), displacement_units);

	w_init_field(&ui->rpm[0]);

	w_init_field(&ui->map[0]);
	w_init_select(&ui->map_unit, nelem(map_units), map_units);

	w_init_field(&ui->ve[0]);

	ui->npoints = 1;
}
