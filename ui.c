#include <assert.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "microui.h"
#include "unit.h"
#include "compressor.h"
#include "widget.h"
#include "engine.h"
#include "ui.h"
#include "eprintf.h"


#define DEFAULT_DISPLACEMENT (litre(1.5))
#define DEFAULT_AMBIENT_TEMPERATURE (celsius(20))
#define DEFAULT_AMBIENT_PRESSURE (millibar(1013))
#define DEFAULT_RPM (rpm(2000))
#define DEFAULT_MAP (millibar(1013))
#define DEFAULT_VE (percent(90))
#define DEFAULT_COMPRESSOR_EFFICIENCY (percent(70))
#define DEFAULT_INTERCOOLER_EFFICIENCY (percent(90))
#define DEFAULT_INTERCOOLER_DELTAP (psi(0.2))


static void init_displacement(UI *ui);
static void init_ambient_temperature(UI *ui);
static void init_ambient_pressure(UI *ui);
static void init_rpm(UI *ui);
static void init_map(UI *ui);
static void init_ve(UI *ui);
static void init_comp_efficiency(UI *ui);
static void init_intercooler_efficiency(UI *ui);
static void init_intercooler_deltap(UI *ui);
static void init_pressure_ratio(UI *ui);
static void init_comp_outlet_temperature(UI *ui);
static void init_manifold_temperature(UI *ui);
static void init_volume_flow_rate(UI *ui);
static void init_mass_flow_rate(UI *ui);
static void init_mass_flow_rate_corrected(UI *ui);
static int init_comps(UI *ui);
static void compute_pressure_ratio(UI *ui, int idx);
static void compute_comp_outlet_temperature(UI *ui, int idx);
static void compute_manifold_temperature(UI *ui, int idx);
static void compute_volume_flow_rate(UI *ui, int idx);
static void compute_mass_flow_rate(UI *ui, int idx);
static void compute_mass_flow_rate_corrected(UI *ui, int idx);


/* Returns non-zero on error. */
int
init_ui(UI *ui) {
	ui->npoints = 1;
	init_engine(&ui->points[0]);

	init_displacement(ui);
	init_ambient_temperature(ui);
	init_ambient_pressure(ui);

	init_rpm(ui);
	init_map(ui);
	init_ve(ui);
	init_comp_efficiency(ui);
	init_intercooler_efficiency(ui);
	init_intercooler_deltap(ui);

	init_pressure_ratio(ui);
	init_comp_outlet_temperature(ui);
	init_manifold_temperature(ui);
	init_volume_flow_rate(ui);
	init_mass_flow_rate(ui);
	init_mass_flow_rate_corrected(ui);

	if (init_comps(ui) != 0) {
		return 1;
	}

	compute(ui, 0);

	return 0;
}

void
free_ui(UI *ui) {
	w_free_select_compressor(&ui->comp_select);
	free(ui->comps);
}

static void
init_displacement(UI *ui) {
	int i;
	double v;

	w_init_field(&ui->displacement);
	w_init_select(&ui->displacement_unit, n_volume_units, volume_units);
	i = ui->displacement_unit.idx;
	v = volume_readers[i](DEFAULT_DISPLACEMENT);
	w_set_field(&ui->displacement, v);

	set_displacement(ui);
}

static void
init_ambient_temperature(UI *ui) {
	int i;
	double v;

	w_init_field(&ui->ambient_temperature);
	w_init_select(&ui->ambient_temperature_unit, n_temperature_units, temperature_units);
	i = ui->ambient_temperature_unit.idx;
	v = temperature_readers[i](DEFAULT_AMBIENT_TEMPERATURE);
	w_set_field(&ui->ambient_temperature, v);

	set_ambient_temperature(ui);
}

static void
init_ambient_pressure(UI *ui) {
	int i;
	double v;

	w_init_field(&ui->ambient_pressure);
	w_init_select(&ui->ambient_pressure_unit, n_pressure_units, pressure_units);
	i = ui->ambient_pressure_unit.idx;
	v = pressure_readers[i](DEFAULT_AMBIENT_PRESSURE);
	w_set_field(&ui->ambient_pressure, v);

	set_ambient_pressure(ui);
}

static void
init_rpm(UI *ui) {
	w_init_field(&ui->rpm[0]);
	w_set_field(&ui->rpm[0], as_rpm(DEFAULT_RPM));

	set_rpm(ui, 0);
}

static void
init_map(UI *ui) {
	int i;
	double v;

	w_init_field(&ui->map[0]);
	w_init_select(&ui->map_unit, n_pressure_units, pressure_units);
	i = ui->map_unit.idx;
	v = pressure_readers[i](DEFAULT_MAP);
	w_set_field(&ui->map[0], v);

	set_map(ui, 0);
}

static void
init_ve(UI *ui) {
	w_init_field(&ui->ve[0]);
	w_set_field(&ui->ve[0], as_percent(DEFAULT_VE));

	set_ve(ui, 0);
}

static void
init_comp_efficiency(UI *ui) {
	w_init_field(&ui->comp_efficiency[0]);
	w_set_field(&ui->comp_efficiency[0], as_percent(DEFAULT_COMPRESSOR_EFFICIENCY));

	set_comp_efficiency(ui, 0);
}

static void
init_intercooler_efficiency(UI *ui) {
	w_init_field(&ui->intercooler_efficiency[0]);
	w_set_field(&ui->intercooler_efficiency[0], as_percent(DEFAULT_INTERCOOLER_EFFICIENCY));

	set_intercooler_efficiency(ui, 0);
}

static void
init_intercooler_deltap(UI *ui) {
	int i;
	double v;

	w_init_field(&ui->intercooler_deltap[0]);
	w_init_select(&ui->intercooler_deltap_unit, n_pressure_units, pressure_units);
	i = ui->intercooler_deltap_unit.idx;
	v = pressure_readers[i](DEFAULT_INTERCOOLER_DELTAP);
	w_set_field(&ui->intercooler_deltap[0], v);

	set_intercooler_deltap(ui, 0);
}

static void
init_pressure_ratio(UI *ui) {
	w_init_number(ui->pressure_ratio[0]);
}

static void
init_comp_outlet_temperature(UI *ui) {
	w_init_number(ui->comp_outlet_temperature[0]);
	w_init_select(&ui->comp_outlet_temperature_unit, n_temperature_units, temperature_units);
}

static void
init_manifold_temperature(UI *ui) {
	w_init_number(ui->manifold_temperature[0]);
	w_init_select(&ui->manifold_temperature_unit, n_temperature_units, temperature_units);
}

static void
init_volume_flow_rate(UI *ui) {
	w_init_select(&ui->volume_flow_rate_unit, n_volume_flow_rate_units, volume_flow_rate_units);
	w_init_number(ui->volume_flow_rate[0]);
}

static void
init_mass_flow_rate(UI *ui) {
	w_init_select(&ui->mass_flow_rate_unit, n_mass_flow_rate_units, mass_flow_rate_units);
	w_init_number(ui->mass_flow_rate[0]);
}

static void
init_mass_flow_rate_corrected(UI *ui) {
	w_init_select(&ui->mass_flow_rate_corrected_unit, n_mass_flow_rate_units, mass_flow_rate_units);
	w_init_number(ui->mass_flow_rate_corrected[0]);
}

static int
init_comps(UI *ui) {
	int i;

	if (load_compressors(&ui->comps, &ui->ncomps) != 0) {
		weprintf("failed to load compressors");
		return 1;
	}
	for (i = 0; i < ui->ncomps; i++) {
		printf("%s %s %s\n", ui->comps[i].brand, ui->comps[i].series, ui->comps[i].model);
	}

	if (w_init_select_compressor(&ui->comp_select, ui->ncomps, ui->comps) != 0) {
		free(ui->comps);
		return 1;
	}
	printf("init'd comp select widget\n");

	return 0;
}

void
set_displacement(UI *ui) {
	int idx, i;
	VolumeMaker convert;
	Volume disp;

	idx = ui->displacement_unit.idx;
	assert(idx >= 0 && (size_t) idx < n_volume_units);

	convert = volume_makers[idx];
	disp = convert(ui->displacement.value);

	for (i = 0; i < ui->npoints; i++) {
		ui->points[i].displacement = disp;
	}
}

void
set_displacement_unit(UI *ui) {
	VolumeMaker maker;
	Volume disp;
	VolumeReader reader;

	maker = volume_makers[ui->displacement_unit.oldidx];
	disp = maker(ui->displacement.value);
	reader = volume_readers[ui->displacement_unit.idx];
	w_set_field(&ui->displacement, reader(disp));
}

void
set_ambient_temperature(UI *ui) {
	int idx, i;
	TemperatureMaker convert;
	Temperature t;

	idx = ui->ambient_temperature_unit.idx;
	assert(idx >= 0 && (size_t) idx < n_temperature_units);

	convert = temperature_makers[idx];
	t = convert(ui->ambient_temperature.value);

	for (i = 0; i < ui->npoints; i++) {
		ui->points[i].ambient_temperature = t;
	}
}

void
set_ambient_temperature_unit(UI *ui) {
	TemperatureMaker maker;
	Temperature t;
	TemperatureReader reader;

	maker = temperature_makers[ui->ambient_temperature_unit.oldidx];
	t = maker(ui->ambient_temperature.value);
	reader = temperature_readers[ui->ambient_temperature_unit.idx];
	w_set_field(&ui->ambient_temperature, reader(t));
}

void
set_ambient_pressure(UI *ui) {
	int idx, i;
	PressureMaker convert;
	Pressure p;

	idx = ui->ambient_pressure_unit.idx;
	assert(idx >= 0 && (size_t) idx < n_pressure_units);

	convert = pressure_makers[idx];
	p = convert(ui->ambient_pressure.value);

	for (i = 0; i < ui->npoints; i++) {
		ui->points[i].ambient_pressure = p;
	}
}

void
set_ambient_pressure_unit(UI *ui) {
	PressureMaker maker;
	Pressure p;
	PressureReader reader;

	maker = pressure_makers[ui->ambient_pressure_unit.oldidx];
	p = maker(ui->ambient_pressure.value);
	reader = pressure_readers[ui->ambient_pressure_unit.idx];
	w_set_field(&ui->ambient_pressure, reader(p));
}

void
set_rpm(UI *ui, int idx) {
	ui->points[idx].rpm = rpm(ui->rpm[idx].value);
}

void
set_map(UI *ui, int idx) {
	int unit_idx;
	PressureMaker convert;
	Pressure p;

	unit_idx = ui->map_unit.idx;
	assert(unit_idx >= 0 && (size_t) unit_idx < n_pressure_units);

	convert = pressure_makers[unit_idx];
	p = convert(ui->map[idx].value);
	ui->points[idx].map = p;
}

void
set_map_unit(UI *ui) {
	PressureMaker maker;
	PressureReader reader;
	int i;
	Pressure map;

	maker = pressure_makers[ui->map_unit.oldidx];
	reader = pressure_readers[ui->map_unit.idx];
	for (i = 0; i < ui->npoints; i++) {
		map = maker(ui->map[i].value);
		w_set_field(&ui->map[i], reader(map));
	}
}

void
set_ve(UI *ui, int idx) {
	ui->points[idx].ve = percent(ui->ve[idx].value);
}

void
set_comp_efficiency(UI *ui, int idx) {
	ui->points[idx].comp_efficiency = percent(ui->comp_efficiency[idx].value);
}

void
set_intercooler_efficiency(UI *ui, int idx) {
	ui->points[idx].intercooler_efficiency = percent(ui->intercooler_efficiency[idx].value);
}

void
set_intercooler_deltap(UI *ui, int idx) {
	int unit_idx;
	PressureMaker convert;
	Pressure p;

	unit_idx = ui->intercooler_deltap_unit.idx;
	assert(unit_idx >= 0 && (size_t) unit_idx < n_pressure_units);

	convert = pressure_makers[unit_idx];
	p = convert(ui->intercooler_deltap[idx].value);
	ui->points[idx].intercooler_deltap = p;
}

void
set_intercooler_deltap_unit(UI *ui) {
	PressureMaker maker;
	PressureReader reader;
	int i;
	Pressure p;

	maker = pressure_makers[ui->intercooler_deltap_unit.oldidx];
	reader = pressure_readers[ui->intercooler_deltap_unit.idx];
	for (i = 0; i < ui->npoints; i++) {
		p = maker(ui->intercooler_deltap[i].value);
		w_set_field(&ui->intercooler_deltap[i], reader(p));
	}
}

void
compute(UI *ui, int idx) {
	compute_pressure_ratio(ui, idx);
	compute_comp_outlet_temperature(ui, idx);
	compute_manifold_temperature(ui, idx);
	compute_volume_flow_rate(ui, idx);
	compute_mass_flow_rate(ui, idx);
	compute_mass_flow_rate_corrected(ui, idx);
}

void
compute_all(UI *ui) {
	int i;

	for (i = 0; i < ui->npoints; i++) {
		compute(ui, i);
	}
}

static void
compute_pressure_ratio(UI *ui, int idx) {
	double pr;

	pr = pressure_ratio(&ui->points[idx]);
	w_set_number(ui->pressure_ratio[idx], pr);
}

static void
compute_comp_outlet_temperature(UI *ui, int idx) {
	int unit_idx;
	TemperatureReader convert;
	double v;

	unit_idx = ui->comp_outlet_temperature_unit.idx;
	assert(unit_idx >= 0 && (size_t) unit_idx < n_temperature_units);

	convert = temperature_readers[unit_idx];
	v = convert(comp_outlet_temperature(&ui->points[idx]));
	w_set_number(ui->comp_outlet_temperature[idx], v);
}

static void
compute_manifold_temperature(UI *ui, int idx) {
	int unit_idx;
	TemperatureReader convert;
	double v;

	unit_idx = ui->manifold_temperature_unit.idx;
	assert(unit_idx >= 0 && (size_t) unit_idx < n_temperature_units);

	convert = temperature_readers[unit_idx];
	v = convert(manifold_temperature(&ui->points[idx]));
	w_set_number(ui->manifold_temperature[idx], v);
}

static void
compute_volume_flow_rate(UI *ui, int idx) {
	int unit_idx;
	VolumeFlowRateReader convert;
	double v;

	unit_idx = ui->volume_flow_rate_unit.idx;
	assert(unit_idx >= 0 && (size_t) unit_idx < n_volume_flow_rate_units);

	convert = volume_flow_rate_readers[unit_idx];
	v = convert(volume_flow_rate(&ui->points[idx]));
	w_set_number(ui->volume_flow_rate[idx], v);
}

static void
compute_mass_flow_rate(UI *ui, int idx) {
	int unit_idx;
	MassFlowRateReader convert;
	double v;

	unit_idx = ui->mass_flow_rate_unit.idx;
	assert(unit_idx >= 0 && (size_t) unit_idx < n_mass_flow_rate_units);

	convert = mass_flow_rate_readers[unit_idx];
	v = convert(mass_flow_rate(&ui->points[idx]));
	w_set_number(ui->mass_flow_rate[idx], v);
}

static void
compute_mass_flow_rate_corrected(UI *ui, int idx) {
	int unit_idx;
	MassFlowRateReader convert;
	double v;

	unit_idx = ui->mass_flow_rate_corrected_unit.idx;
	assert(unit_idx >= 0 && (size_t) unit_idx < n_mass_flow_rate_units);

	convert = mass_flow_rate_readers[unit_idx];
	v = convert(mass_flow_rate_corrected(&ui->points[idx]));
	w_set_number(ui->mass_flow_rate_corrected[idx], v);
}

void
insert_point(UI *ui, int idx) {
	int i;

	if (idx < 0 || idx >= ui->npoints || ui->npoints >= MAX_POINTS) {
		return;
	}

	for (i = ui->npoints; i > idx; i--) {
		memmove(&ui->points[i], &ui->points[i-1], sizeof(ui->points[i-1]));

		memmove(&ui->rpm[i], &ui->rpm[i-1], sizeof(ui->rpm[i-1]));
		memmove(&ui->map[i], &ui->map[i-1], sizeof(ui->map[i-1]));
		memmove(&ui->ve[i], &ui->ve[i-1], sizeof(ui->ve[i-1]));
		memmove(&ui->comp_efficiency[i], &ui->comp_efficiency[i-1], sizeof(ui->comp_efficiency[i-1]));
		memmove(&ui->intercooler_efficiency[i], &ui->intercooler_efficiency[i-1], sizeof(ui->intercooler_efficiency[i-1]));
		memmove(&ui->intercooler_deltap[i], &ui->intercooler_deltap[i-1], sizeof(ui->intercooler_deltap[i-1]));

		memmove(&ui->pressure_ratio[i], &ui->pressure_ratio[i-1], sizeof(ui->pressure_ratio[i-1]));
		memmove(&ui->comp_outlet_temperature[i], &ui->comp_outlet_temperature[i-1], sizeof(ui->comp_outlet_temperature[i-1]));
		memmove(&ui->manifold_temperature[i], &ui->manifold_temperature[i-1], sizeof(ui->manifold_temperature[i-1]));
		memmove(&ui->volume_flow_rate[i], &ui->volume_flow_rate[i-1], sizeof(ui->volume_flow_rate[i-1]));
		memmove(&ui->mass_flow_rate[i], &ui->mass_flow_rate[i-1], sizeof(ui->mass_flow_rate[i-1]));
		memmove(&ui->mass_flow_rate_corrected[i], &ui->mass_flow_rate_corrected[i-1], sizeof(ui->mass_flow_rate_corrected[i-1]));
	}
	ui->npoints++;
}

void
remove_point(UI *ui, int idx) {
	if (idx < 0 || idx >= ui->npoints || ui->npoints <= 1) {
		return;
	}

	for (; idx < ui->npoints-1; idx++) {
		memmove(&ui->points[idx], &ui->points[idx+1], sizeof(ui->points[idx]));

		memmove(&ui->rpm[idx], &ui->rpm[idx+1], sizeof(ui->rpm[idx]));
		memmove(&ui->map[idx], &ui->map[idx+1], sizeof(ui->map[idx]));
		memmove(&ui->ve[idx], &ui->ve[idx+1], sizeof(ui->ve[idx]));
		memmove(&ui->comp_efficiency[idx], &ui->comp_efficiency[idx+1], sizeof(ui->comp_efficiency[idx]));
		memmove(&ui->intercooler_efficiency[idx], &ui->intercooler_efficiency[idx+1], sizeof(ui->intercooler_efficiency[idx]));
		memmove(&ui->intercooler_deltap[idx], &ui->intercooler_deltap[idx+1], sizeof(ui->intercooler_deltap[idx]));

		memmove(&ui->pressure_ratio[idx], &ui->pressure_ratio[idx+1], sizeof(ui->pressure_ratio[idx]));
		memmove(&ui->comp_outlet_temperature[idx], &ui->comp_outlet_temperature[idx+1], sizeof(ui->comp_outlet_temperature[idx]));
		memmove(&ui->manifold_temperature[idx], &ui->manifold_temperature[idx+1], sizeof(ui->manifold_temperature[idx]));
		memmove(&ui->volume_flow_rate[idx], &ui->volume_flow_rate[idx+1], sizeof(ui->volume_flow_rate[idx]));
		memmove(&ui->mass_flow_rate[idx], &ui->mass_flow_rate[idx+1], sizeof(ui->mass_flow_rate[idx]));
		memmove(&ui->mass_flow_rate_corrected[idx], &ui->mass_flow_rate_corrected[idx+1], sizeof(ui->mass_flow_rate_corrected[idx]));
	}
	ui->npoints--;
}
