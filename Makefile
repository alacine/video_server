.DELETE_ON_ERROR:

ALL_SERVICES = api scheduler streamserver deployserver

.PHONY: all clean status startdb run-deamo stop $(ALL_SERVICES)

all: $(ALL_SERVICES)

$(ALL_SERVICES):
	$(MAKE) -C $@

status:
	ps aux | grep -E 'api|streamserver|scheduler|deployserver' | grep -v grep

startdb:
	docker start mysql-test

run-deamon: | startdb $(ALL_SERVICES)
	cd streamserver && nohup ./streamserver &
	cd scheduler && nohup ./scheduler &
	cd api && nohup ./api &

stop:
	./admin.sh

stopall: stop
	docker stop mysql-test

clean: stop
	@for dir in $(ALL_SERVICES); do \
		$(MAKE) -C $$dir $@; \
	done
	@#find . -type f ! -regex '^\./\.git/.*' ! -regex '.+\..+' ! -name Makefile -delete
	find . -name nohup.out -delete
