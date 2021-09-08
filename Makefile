#.DELETE_ON_ERROR:

ALL_SERVICES = api scheduler streamserver

.PHONY: all clean status startdb run-deamo stop $(ALL_SERVICES) help

all: $(ALL_SERVICES)

$(ALL_SERVICES):
	$(MAKE) -C $@

status:
	pgrep -au $$USER "api|streamserver|scheduler|deployserver"

startdb:
	docker-compose start db

run-daemon: | startdb $(ALL_SERVICES)
	cd streamserver && nohup ./streamserver &
	cd scheduler && nohup ./scheduler &
	cd api && nohup ./api &

build-in-docker:
	mkdir -pv local-cache/db local-cache/api local-cache/scheduler local-cache/streamserver/videos
	docker build . -t video_server_build -f build.Dockerfile
	docker build . -t video_server_base -f base.Dockerfile

install stop:
	@for dir in $(ALL_SERVICES); do \
		$(MAKE) -C $$dir $@; \
	done

stopall: stop
	docker-compose stop db

clean:
	@for dir in $(ALL_SERVICES); do \
		$(MAKE) -C $$dir $@; \
	done
	@#find . -type f ! -regex '^\./\.git/.*' ! -regex '.+\..+' ! -name Makefile -delete
	find . -path ./local-cache -prune -false -or -name nohup.out -exec rm -f {} \;
	docker rmi video_server_build
	docker rmi video_server_base

restore: clean
	@for dir in $(ALL_SERVICES); do \
		$(MAKE) -C $$dir $@; \
	done
	docker-compose down -v
	sudo rm -rf ./local-cache

help:
	@echo "(none):          build all submodules"
	@echo "startdb:         start database in local docker"
	@echo "run-deamon:      start submodules in local environment"
	@echo "build-in-docker: build two docker images:"
	@echo "                     1. video_server_build: contains all built submodules binary"
	@echo "                     2. video_server_base: base image for submodules' image"
	@echo "install:         install submodules in GOPATH/bin"
	@echo "stop:            kill all local submodules' process"
	@echo "stopall:         do 'stop' and stop database docker"
	@echo "clean:           delete all binarys, nohup logs and images: video_server_base, video_server_build"
	@echo "restore:         do 'clean' and database volume"
	@echo ""
	@echo "If you have docker-compose installed, you can try docker-compose as you like."
