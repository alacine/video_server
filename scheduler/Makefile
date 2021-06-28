.PHONY: all run clean stop restore
#.DELETE_ON_ERROR:

rwildcard = $(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2))
target = $(shell basename $(CURDIR))

all: $(target)

$(target): $(call rwildcard,./,*.go)
	go build

run: $(target)
	./$(target)

docker-build: Dockerfile
	docker build . -t video_server/$(target)

docker-run:
	docker run -d --name local-$(target) video_server/$(target)

docker-run-once:
	docker run --rm -it video_server/$(target) bash

docker-stop:
	docker stop local-$(target)

docker-rm: docker-stop
	docker rm local-$(target)

stop:
	pkill -u $$USER $(target)

clean:
	go clean -x
	docker rm local-$(target)

install: 
	go install

uninstall:
	go clean -i -x

restore: clean
