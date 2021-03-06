TEAMID ?= 7357
repo=tompscanlan/q3errord
bin=q3errord

all: $(bin)-local

$(bin): deps
	CGO_ENABLED=0 go build -a -v --installsuffix cgo  ./cmd/$(bin)

$(bin)-local: deps
	go build -v -o $(bin)-local  ./cmd/$(bin)

deps:
	go get -v ./...

docker: $(bin)
	docker build -t $(repo) --rm=true .

dockerclean: stop
	echo "Cleaning up Docker Engine before building."
	docker rm $$(docker ps -a | awk '/$(bin)/ { print $$1}') || echo -
	docker rmi $(repo)

clean: stop dockerclean
	go clean
	rm -f $(bin)

run:
	docker run -d -p9999:80   $(repo)

stop:
	docker kill $$(docker ps -a | awk '/$(bin)/ { print $$1}') || echo -

valid:
	go tool vet .
	go test -v -race ./...


.PHONY: imports docker clean run stop deps

