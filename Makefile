#.DELETE_ON_ERROR:

ALL_SERVICES = api scheduler streamserver

.PHONY: all clean status startdb run-deamo stop $(ALL_SERVICES)

all: $(ALL_SERVICES)

$(ALL_SERVICES):
	$(MAKE) -C $@

status:
	pgrep -au $$USER "api|streamserver|scheduler|deployserver"

startdb:
	docker-compose start -d db

run-deamon: | startdb $(ALL_SERVICES)
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
