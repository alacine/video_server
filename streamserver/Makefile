.PHONY: all run lint stop clean install uninstall restore
#.DELETE_ON_ERROR:

rwildcard = $(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2))
target = $(shell basename $(CURDIR))

all: $(target)

$(target): $(call rwildcard,./,*.go)
	go build

run: $(target)
	./$(target)

lint:
	errcheck ./...
	revive ./...
	golangci-lint run

stop:
	pkill -u $$USER $(target)

clean:
	go clean -x

install: 
	go install

uninstall:
	go clean -i -x

restore: clean uninstall
	rm -rf videos
