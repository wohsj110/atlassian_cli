.PHONY: test build clean

test:
	go test ./shared/... ./tools/atk-jira/... ./tools/atk-cfl/...

build:
	mkdir -p bin
	go build -o bin/atk-jira ./tools/atk-jira/cmd/atk-jira
	go build -o bin/atk-cfl ./tools/atk-cfl/cmd/atk-cfl

clean:
	rm -rf bin
