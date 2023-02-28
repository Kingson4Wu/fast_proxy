build:
	go mod download
	go build -o in-proxy ./examples/inproxy
	go build -o out-proxy ./examples/outproxy
