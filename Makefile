API = api/api
SCHEDULER = scheduler/scheduler
STREAMSERVER = streamserver/streamserver
DEPLOYSERVER = deployserver/deployserver

all: $(API) $(SCHEDULER) $(STREAMSERVER) $(DEPLOYSERVER)

$(API): api/*.go
	cd api && go build

$(SCHEDULER): scheduler/*.go
	cd scheduler && go build

$(STREAMSERVER): streamserver/*.go
	cd streamserver && go build

$(DEPLOYSERVER): deployserver/*.go
	cd deployserver && go build

run:
	cd streamserver && nohup ./streamserver &
	cd scheduler && nohup ./scheduler &
	cd api && nohup ./api &

stop:
	./admin.sh

clean:
	find . -type f ! -regex '^\./\.git/.*' ! -regex '.+\..+' ! -name Makefile -delete
	find . -name nohup.out -delete
