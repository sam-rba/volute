CFLAGS = -std=c99 -Wall -Wextra -pedantic -Wno-deprecated-declarations
LDFLAGS = -lSDL2 -lSDL2_ttf

SRC = main.c microui.c renderer.c widget.c ui.c
OBJ = ${SRC:.c=.o}

all: volute

clean:
	rm -f volute *.o

volute: ${OBJ}
	${CC} -o $@ $^ ${LDFLAGS}

%.o: %.c
	${CC} -c ${CFLAGS} $<

${OBJ}: microui.h renderer.h widget.h ui.h
