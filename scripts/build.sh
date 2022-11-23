#!/bin/bash

echo "Build qtun for Linux"
#Linux amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/vtun-linux-amd64 ./main.go
#Linux arm64
# CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/vtun-linux-arm64 ./main.go
#Mac amd64
# CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/vtun-darwin-amd64 ./main.go
#Mac arm64
# CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/vtun-darwin-arm64 ./main.go
#Windows amd64
# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/vtun-win-amd64.exe ./main.go


echo "Build HTTP test server"
## Testing client and server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/clienttest ./testing/client/clienttest.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/servertest ./testing/server/servertest.go

echo "DONE!!!"
