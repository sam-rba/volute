CFLAGS = -std=c99 -Wall -Wextra -pedantic -Wno-deprecated-declarations -D_XOPEN_SOURCE=700L
LDFLAGS = -lSDL2 -lSDL2_ttf

SRC = main.c microui.c renderer.c widget.c ui.c unit.c
OBJ = ${SRC:.c=.o}

all: volute

clean:
	rm -f volute *.o

test: test_unit
	for t in $^; do \
		./$$t; \
	done

test_unit: test_unit.o unit.o
	${CC} -o $@ $^ ${LDFLAGS}

volute: ${OBJ}
	${CC} -o $@ $^ ${LDFLAGS}

%.o: %.c
	${CC} -c ${CFLAGS} $<

${OBJ}: microui.h renderer.h widget.h ui.h unit.h
