SHELL:=/bin/bash
TARGET=conductor
SRC=$(wildcard src/*.go)

GDMP_VERSION=1.0.0

all: $(TARGET)

$(TARGET):
	if [ ! -f build/embed/gdmp/$(GDMP_VERSION) ]; then 																							\
		rm -f -r build/embed/gdmp; 																												\
		wget --output-document=build/embed/gdmp/$(GDMP_VERSION) https://github.com/MrNavaStar/GDMP/releases/download/$(GDMP_VERSION)/gdmp; 		\
	fi; 																																		\
	mkdir src/embed
	cp build/embed/gdmp/$(GDMP_VERSION) src/embed/gdmp;
	go build -o build/$(TARGET) $(SRC);
	rm -f -r src/embed;

clean:
	rm -f -r build

run: $(TARGET)
	./build/$(TARGET)

install: $(TARGET)
	go install

.PHONEY: clean run install