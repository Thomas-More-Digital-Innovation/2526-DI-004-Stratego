lint:
    cd code/backend && golangci-lint run --config=../../.golangci.yaml

fmt:
    cd code/backend && gofmt -l .