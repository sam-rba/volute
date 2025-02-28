#include <string.h>

#include "unit.h"
#include "engine.h"

void
init_engine(Engine *e) {
	memset(e, 0, sizeof(*e));
}
