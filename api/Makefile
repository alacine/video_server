.PHONY: all run gomod clean stop restore docker-%
#.DELETE_ON_ERROR:

rwildcard = $(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2))
target = $(shell basename $(CURDIR))

all: $(target)

$(target): $(call rwildcard,./,*.go)
	go build

run: $(target)
	./$(target)

stop:
	pkill -u $$USER $(target)

clean:
	go clean -x

install: 
	go install
	cp config/conf.json $$GOPATH/bin

uninstall:
	go clean -i -x
	rm -f $$GOPATH/bin/conf.json

restore: clean uninstall
