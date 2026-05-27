cd code/backend &&
go get github.com/swaggo/swag/cmd/swag@v1.16.6 &&
go run github.com/swaggo/swag/cmd/swag init -g cmd/server/main.go &&
npx -y @asyncapi/cli@latest generate fromTemplate docs/asyncapi.yaml @asyncapi/html-template -o docs/asyncapi-docs --force-write