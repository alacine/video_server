.PHONY: all run clean
.DELETE_ON_ERROR:

rwildcard = $(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2))
target = $(shell basename $(CURDIR))

all: $(target)

$(target): $(call rwildcard,./,*.go)
	go build

run: $(target)
	./$(target)

clean:
	rm -f $(target)
