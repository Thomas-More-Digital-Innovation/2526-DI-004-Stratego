lint:
    cd code/backend && golangci-lint run --config=../../.golangci.yaml

fmt:
    cd code/backend && gofmt -l .

dev: 
    docker compose up --build

simulation:
    cd code/backend && go run cmd/simulation/main.go 

update-swagger:
    ./scripts/update-swagger.sh

update-asyncapi:
    ./scripts/update-asyncapi.sh

update-docs:
    ./scripts/update-swagger.sh && ./scripts/update-asyncapi.sh