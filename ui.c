#include <assert.h>
#include <stdio.h>
#include <string.h>

#include "microui.h"
#include "widget.h"
#include "unit.h"
#include "engine.h"
#include "ui.h"


#define nelem(arr) (sizeof(arr)/sizeof(arr[0]))


static const char *const volume_units[] = {"cc", "l", "ci"};
static Volume (*volume_converters[nelem(volume_units)])(double) = {
	cubic_centimetre, litre, cubic_inch,
};

static const char *const pressure_units[] = {"mbar", "kPa", "bar", "psi"};
static Pressure (*pressure_converters[nelem(pressure_units)])(double) = {
	millibar, kilopascal, bar, psi,
};

static const char *const volume_flow_rate_units[] = {"m³/s", "CFM"};
static double (*volume_flow_rate_converters[nelem(volume_flow_rate_units)])(VolumeFlowRate) = {
	as_cubic_metre_per_sec, as_cubic_foot_per_min,
};


void
init_ui(UI *ui) {
	w_init_field(&ui->displacement);
	w_init_select(&ui->displacement_unit, nelem(volume_units), volume_units);
	
	ui->npoints = 1;

	w_init_field(&ui->rpm[0]);

	w_init_field(&ui->map[0]);
	w_init_select(&ui->map_unit, nelem(pressure_units), pressure_units);

	w_init_field(&ui->ve[0]);

	init_engine(&ui->points[0]);

	w_init_select(&ui->volume_flow_rate_unit, nelem(volume_flow_rate_units), volume_flow_rate_units);
	w_init_label(ui->volume_flow_rate[0]);
}

void
set_displacement(UI *ui) {
	int idx, i;
	Volume (*convert)(double), disp;

	idx = ui->displacement_unit.idx;
	assert(idx >= 0 && (long unsigned int) idx < nelem(volume_units));

	convert = volume_converters[idx];
	disp = convert(ui->displacement.value);

	for (i = 0; i < ui->npoints; i++) {
		ui->points[i].displacement = disp;
	}
}

void
set_map(UI *ui, int idx) {
	int unit_idx;
	Pressure (*convert)(double), p;

	unit_idx = ui->map_unit.idx;
	assert(unit_idx >= 0 && (long unsigned int) unit_idx < nelem(pressure_units));

	convert = pressure_converters[unit_idx];
	p = convert(ui->map[idx].value);
	ui->points[idx].map = p;
}

void
set_ve(UI *ui, int idx) {
	ui->points[idx].ve = percent(ui->ve[idx].value);
}

void
set_volume_flow_rate(UI *ui, int idx) {
	int unit_idx;
	double (*convert)(VolumeFlowRate), v;

	unit_idx = ui->volume_flow_rate_unit.idx;
	assert(unit_idx >= 0 && (long unsigned int) unit_idx < nelem(volume_flow_rate_units));

	convert = volume_flow_rate_converters[unit_idx];
	v = convert(volume_flow_rate(&ui->points[idx]));
	snprintf(ui->volume_flow_rate[idx], sizeof(ui->volume_flow_rate[idx]), "%lf", v);
}

void
set_all_volume_flow_rate(UI *ui) {
	int i;

	for (i = 0; i < ui->npoints; i++) {
		set_volume_flow_rate(ui, i);
	}
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
		memmove(&ui->volume_flow_rate[i], &ui->volume_flow_rate[i-1], sizeof(ui->volume_flow_rate[i-1]));
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
		memmove(&ui->volume_flow_rate[idx], &ui->volume_flow_rate[idx+1], sizeof(ui->volume_flow_rate[idx]));
	}
	ui->npoints--;
}
