#include <assert.h>
#include <string.h>

#include "microui.h"
#include "widget.h"
#include "unit.h"
#include "engine.h"
#include "ui.h"


#define nelem(arr) (sizeof(arr)/sizeof(arr[0]))


static const char *const volume_units[] = {"cc", "l", "ci"};
static const VolumeMaker volume_makers[nelem(volume_units)] = {
	cubic_centimetre, litre, cubic_inch,
};
static const VolumeReader volume_readers[nelem(volume_units)] = {
	as_cubic_centimetre, as_litre, as_cubic_inch,
};

static const char *const temperature_units[] = {"°C", "K", "°F", "°R"};
static const TemperatureMaker temperature_makers[nelem(temperature_units)] = {
	celsius, kelvin, fahrenheit, rankine,
};
static const TemperatureReader temperature_readers[nelem(temperature_units)] = {
	as_celsius, as_kelvin, as_fahrenheit, as_rankine,
};

static const char *const pressure_units[] = {"mbar", "kPa", "bar", "psi"};
static const PressureMaker pressure_makers[nelem(pressure_units)] = {
	millibar, kilopascal, bar, psi,
};
static const PressureReader pressure_readers[nelem(pressure_units)] = {
	as_millibar, as_kilopascal, as_bar, as_psi,
};

static const char *const volume_flow_rate_units[] = {"m³/s", "CFM"};
static const VolumeFlowRateReader volume_flow_rate_readers[nelem(volume_flow_rate_units)] = {
	as_cubic_metre_per_sec, as_cubic_foot_per_min,
};


void
init_ui(UI *ui) {
	w_init_field(&ui->displacement);
	w_init_select(&ui->displacement_unit, nelem(volume_units), volume_units);

	w_init_field(&ui->ambient_temperature);
	w_init_select(&ui->ambient_temperature_unit, nelem(temperature_units), temperature_units);
	
	w_init_field(&ui->ambient_pressure);
	w_init_select(&ui->ambient_pressure_unit, nelem(pressure_units), pressure_units);

	ui->npoints = 1;

	w_init_field(&ui->rpm[0]);

	w_init_field(&ui->map[0]);
	w_init_select(&ui->map_unit, nelem(pressure_units), pressure_units);

	w_init_field(&ui->ve[0]);

	init_engine(&ui->points[0]);

	w_init_select(&ui->volume_flow_rate_unit, nelem(volume_flow_rate_units), volume_flow_rate_units);
	w_init_number(ui->volume_flow_rate[0]);
}

void
set_displacement(UI *ui) {
	int idx, i;
	VolumeMaker convert;
	Volume disp;

	idx = ui->displacement_unit.idx;
	assert(idx >= 0 && (long unsigned int) idx < nelem(volume_units));

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
	assert(idx >= 0 && (long unsigned int) idx < nelem(temperature_units));

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
	assert(idx >= 0 && (long unsigned int) idx < nelem(pressure_units));

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
	assert(unit_idx >= 0 && (long unsigned int) unit_idx < nelem(pressure_units));

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
set_volume_flow_rate(UI *ui, int idx) {
	int unit_idx;
	VolumeFlowRateReader convert;
	VolumeFlowRate v;

	unit_idx = ui->volume_flow_rate_unit.idx;
	assert(unit_idx >= 0 && (long unsigned int) unit_idx < nelem(volume_flow_rate_units));

	convert = volume_flow_rate_readers[unit_idx];
	v = convert(volume_flow_rate(&ui->points[idx]));
	w_set_number(ui->volume_flow_rate[idx], v);
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
		memmove(&ui->comp_efficiency[i], &ui->comp_efficiency[i-1], sizeof(ui->comp_efficiency[i-1]));
		memmove(&ui->intercooler_efficiency[i], &ui->intercooler_efficiency[i-1], sizeof(ui->intercooler_efficiency[i-1]));
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
		memmove(&ui->comp_efficiency[idx], &ui->comp_efficiency[idx+1], sizeof(ui->comp_efficiency[idx]));
		memmove(&ui->intercooler_efficiency[idx], &ui->intercooler_efficiency[idx+1], sizeof(ui->intercooler_efficiency[idx]));
		memmove(&ui->points[idx], &ui->points[idx+1], sizeof(ui->points[idx]));
		memmove(&ui->volume_flow_rate[idx], &ui->volume_flow_rate[idx+1], sizeof(ui->volume_flow_rate[idx]));
	}
	ui->npoints--;
}
