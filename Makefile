go-gen:
	go generate ./...

gen: go-gen

lint:
	golangci-lint run ./...