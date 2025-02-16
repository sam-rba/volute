CFLAGS = -std=c99 -Wall -Wextra -pedantic -Wno-deprecated-declarations
LDFLAGS = $(shell sdl2-config --libs)

# Link OpenGL.
ifeq ($(OS),Windows_NT)
	GLFLAG := -lopengl32
else
	UNAME := `uname -o 2>/dev/null || uname -s`
	ifeq ($(UNAME),"Darwin")
		GLFLAG := -framework OpenGL
	else
		GLFLAG := -lGL
	endif
endif
LDFLAGS += $(GLFLAG)

SRC = main.c microui.c renderer.c widget.c ui.c
OBJ = ${SRC:.c=.o}

all: volute

clean:
	rm -f volute *.o

volute: ${OBJ}
	${CC} -o $@ $^ ${LDFLAGS}

%.o: %.c
	${CC} -c ${CFLAGS} $<

${SRC}: microui.h renderer.h atlas.inl widget.h ui.h
