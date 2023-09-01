build:
	go build -o cvm cmd/cvm.go
run: build
	./cvm
test:
	go test ./...