.PHONY: build clean doc test vet

lib_name := appoptics-go

build:
	go build -o $(lib_name)

clean:
	rm $(lib_name)

doc:
	godoc -http=:8080 -index

test:
	go test ./...

live_test:
	cd _live-tests && go test ./...

super_test: test
super_test: live_test

vet:
	go vet

release:
	git tag -a $(shell go run cmd/main.go)
