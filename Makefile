default: build

test:
	go test ./...

build:
	go build

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-avm ~/.tflint.d/plugins

e2e:
	cd integration && go test -v && cd ../

.PHONY: test build install
