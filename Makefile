build:
	go mod download
	go build -o in-proxy ./examples/inproxy
	go build -o out-proxy ./examples/outproxy
	# use `-x -v` option to print more build log information detail
	chmod +x .githooks/*
	git config core.hooksPath .githooks

test:
	go test -count 1 -v ./... -gcflags "all=-N -l" && go test -race -v ./... -gcflags "-l"

install:

