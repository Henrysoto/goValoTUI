run:
	go run main.go api.go apistructs.go store.go tui.go util.go art.go

build:
	go build -o bin/goValoTUI main.go api.go apistructs.go store.go tui.go util.go art.go

test:
	go clean -testcache && go test -v -cover ./...

.PHONY: run build test

