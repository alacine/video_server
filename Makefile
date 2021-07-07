#.DELETE_ON_ERROR:

ALL_SERVICES = api scheduler streamserver deployserver

.PHONY: all clean status startdb run-deamo stop $(ALL_SERVICES)

all: $(ALL_SERVICES)

$(ALL_SERVICES):
	$(MAKE) -C $@

status:
	pgrep -au $$USER "api|streamserver|scheduler|deployserver"

startdb:
	docker start mysql-test

run-deamon: | startdb $(ALL_SERVICES)
	cd streamserver && nohup ./streamserver &
	cd scheduler && nohup ./scheduler &
	cd api && nohup ./api &

build-in-docker:
	docker build . -t video_server_build

install stop:
	@for dir in $(ALL_SERVICES); do \
		$(MAKE) -C $$dir $@; \
	done

docker-%:
	@for dir in $(ALL_SERVICES); do \
		$(MAKE) -C $$dir $@; \
	done

stopall: stop
	docker stop mysql-test

clean restore:
	@for dir in $(ALL_SERVICES); do \
		$(MAKE) -C $$dir $@; \
	done
	@#find . -type f ! -regex '^\./\.git/.*' ! -regex '.+\..+' ! -name Makefile -delete
	find . -name nohup.out -delete
	docker rmi video_server/build
