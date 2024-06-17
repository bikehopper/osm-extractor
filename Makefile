install:
	go mod download && go mod verify
build:
	go build -o ./bin/osm-extractor-workflow ./workflow/cmd
build_docker: 
	docker build . -t osm-extractor:local
run:
	go run ./workflow/cmd
all: install build