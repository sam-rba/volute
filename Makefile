CFLAGS = -std=c99 -I ./ -fopenmp -Wall -Wextra -pedantic -Wno-deprecated-declarations -D_XOPEN_SOURCE=700L
LDFLAGS = -lSDL2 -lSDL2_ttf -lSDL2_image -lm -fopenmp

SRC = main.c microui.c renderer.c widget.c ui.c unit.c engine.c compressor.c eprintf.c cwalk.c toml.c util.c
OBJ = ${SRC:.c=.o}
HDR = microui.h renderer.h widget.h ui.h unit.h engine.h eprintf.h util.h cwalk.h toml.h

TEST_SRC = test.c test_angular_speed.c test_fraction.c test_pressure.c test_temperature.c test_volume.c test_volume_flow_rate.c test_mass_flow_rate.c test_engine.c unit.c engine.c
TEST_OBJ = ${TEST_SRC:.c=.o}

volute: ${OBJ}
	${CC} -o $@ $^ ${LDFLAGS}

test: ${TEST_OBJ}
	${CC} -o $@ $^ ${LDFLAGS}
	./$@

clean:
	rm -f volute test *.o

%.o: %.c
	${CC} -c ${CFLAGS} $<

${OBJ}: ${HDR}
${TEST_OBJ}: ${HDR} test.h
