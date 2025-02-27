#include "microui.h"
#include "widget.h"
#include "ui.h"


#define nelem(arr) (sizeof(arr)/sizeof(arr[0]))


void
init_ui(UI *ui) {
	w_init_field(&ui->displacement);

	static const char *const displacement_units[] = {"cc", "l", "ci"};
	w_init_select(&ui->displacement_unit, nelem(displacement_units), displacement_units);
}
